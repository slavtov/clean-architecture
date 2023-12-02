package usecase

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/slavtov/clean-architecture/internal/config"
	"github.com/slavtov/clean-architecture/internal/domain/models"
	"github.com/slavtov/clean-architecture/internal/domain/repositories"
	"github.com/slavtov/clean-architecture/internal/domain/usecases"
	"github.com/slavtov/clean-architecture/pkg/logger"
	"github.com/slavtov/clean-architecture/pkg/utils"
)

type usecase struct {
	cfg             *config.Config
	pgRepository    repositories.PGUserRepository
	redisRepository repositories.RedisUserRepository
	log             logger.Logger
}

const cacheDuration = 3600

func New(
	cfg *config.Config,
	pg repositories.PGUserRepository,
	redis repositories.RedisUserRepository,
	log logger.Logger,
) usecases.UserUseCase {
	return &usecase{
		cfg:             cfg,
		pgRepository:    pg,
		redisRepository: redis,
		log:             log,
	}
}

func (u *usecase) GetAll() ([]models.User, error) {
	res, err := u.pgRepository.GetAll()
	if err != nil {
		u.log.Errorf("auth.pgRepository.GetAll: %v", err)
		return res, err
	}

	return res, nil
}

func (u *usecase) GetByID(id uuid.UUID) (models.User, error) {
	cachedUser, err := u.redisRepository.GetByID(id)
	if err != nil {
		u.log.Errorf("auth.redisRepository.GetByID: %v", err)
	}

	if cachedUser.ID != uuid.Nil {
		return cachedUser, nil
	}

	res, err := u.pgRepository.GetByID(id)
	if err != nil {
		u.log.Errorf("auth.pgRepository.GetByID: %v", err)
		return res, err
	}

	res.SanitizePassword()

	if err := u.redisRepository.SetUser(
		&res,
		time.Second*cacheDuration,
	); err != nil {
		u.log.Errorf("auth.redisRepository.SetUser: %v", err)
		return res, err
	}

	return res, nil
}

func (u *usecase) Login(user *models.User) (*models.AuthUser, error) {
	if err := user.Validate(); err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := user.ValidatePassword(); err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	res, err := u.pgRepository.FindByEmail(user.Email)
	if err != nil {
		u.log.Errorf("auth.pgRepository.FindByEmail: %v", err)
		return nil, err
	}

	if err = res.ComparePassword(user.Password); err != nil {
		return nil, echo.ErrUnauthorized
	}

	res.SanitizePassword()

	return u.Auth(&res)
}

func (u *usecase) Store(user *models.User) (*models.AuthUser, error) {
	if err := user.Validate(); err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := user.ValidatePassword(); err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := user.HashPassword(); err != nil {
		return nil, echo.ErrInternalServerError
	}

	res, err := u.pgRepository.Store(user)
	if err != nil {
		u.log.Errorf("auth.pgRepository.Store: %v", err)
		return nil, err
	}

	res.SanitizePassword()

	if err := u.redisRepository.SetUser(
		res,
		time.Second*cacheDuration,
	); err != nil {
		u.log.Errorf("auth.redisRepository.SetUser: %v", err)
		return nil, err
	}

	return u.Auth(res)
}

func (u *usecase) Update(user *models.User) (*models.User, error) {
	if err := user.Validate(); err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if user.Password != "" {
		if err := user.HashPassword(); err != nil {
			return nil, echo.ErrBadRequest
		}
	}

	res, err := u.pgRepository.Update(user)
	if err != nil {
		u.log.Errorf("auth.pgRepository.Update: %v", err)
		return nil, err
	}

	res.SanitizePassword()

	if err := u.redisRepository.SetUser(
		res,
		time.Second*cacheDuration,
	); err != nil {
		u.log.Errorf("auth.redisRepository.SetUser: %v", err)
		return nil, err
	}

	return res, nil
}

func (u *usecase) Delete(id uuid.UUID) error {
	if err := u.pgRepository.Delete(id); err != nil {
		u.log.Errorf("auth.pgRepository.Delete: %v", err)
		return err
	}

	if err := u.redisRepository.Delete(utils.GetRedisKey(
		"users",
		id.String(),
	)); err != nil {
		u.log.Errorf("auth.redisRepository.Delete: %v", err)
		return err
	}

	if err := u.LogoutAll(id); err != nil {
		return err
	}

	return nil
}
