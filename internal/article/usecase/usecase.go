package usecase

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/slavken/go-clean-architecture/internal/domain/models"
	"github.com/slavken/go-clean-architecture/internal/domain/repositories"
	"github.com/slavken/go-clean-architecture/internal/domain/usecases"
	"github.com/slavken/go-clean-architecture/pkg/logger"
)

type usecase struct {
	pgRepository    repositories.PGArticleRepository
	redisRepository repositories.RedisArticleRepository
	log             logger.Logger
}

const cacheDuration = 3600

func New(
	pg repositories.PGArticleRepository,
	redis repositories.RedisArticleRepository,
	log logger.Logger,
) usecases.ArticleUseCase {
	return &usecase{
		pgRepository:    pg,
		redisRepository: redis,
		log:             log,
	}
}

func (u *usecase) GetAll() ([]models.Article, error) {
	res, err := u.pgRepository.GetAll()
	if err != nil {
		u.log.Errorf("article.pgRepository.GetAll: %v", err)
		return res, err
	}

	return res, nil
}

func (u *usecase) GetByID(id uuid.UUID) (models.Article, error) {
	cachedArticle, err := u.redisRepository.GetByID(id)
	if err != nil {
		u.log.Errorf("article.redisRepository.GetByID: %v", err)
	}

	if cachedArticle.ID != uuid.Nil {
		return cachedArticle, nil
	}

	res, err := u.pgRepository.GetByID(id)
	if err != nil {
		u.log.Errorf("article.pgRepository.GetByID: %v", err)
		return res, err
	}

	if err := u.redisRepository.SetArticle(
		&res,
		time.Second*cacheDuration,
	); err != nil {
		u.log.Errorf("article.redisRepository.SetArticle: %v", err)
		return res, err
	}

	return res, nil
}

func (u *usecase) Store(article *models.Article) (*models.Article, error) {
	if err := article.Validate(); err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	res, err := u.pgRepository.Store(article)
	if err != nil {
		u.log.Errorf("article.pgRepository.Store: %v", err)
		return nil, err
	}

	if err := u.redisRepository.SetArticle(
		res,
		time.Second*cacheDuration,
	); err != nil {
		u.log.Errorf("article.redisRepository.SetArticle: %v", err)
		return nil, err
	}

	return res, nil
}

func (u *usecase) Update(article *models.Article) (*models.Article, error) {
	if err := article.Validate(); err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	res, err := u.pgRepository.Update(article)
	if err != nil {
		u.log.Errorf("article.pgRepository.Update: %v", err)
		return nil, err
	}

	if err := u.redisRepository.SetArticle(
		res,
		time.Second*cacheDuration,
	); err != nil {
		u.log.Errorf("article.redisRepository.SetArticle: %v", err)
		return nil, err
	}

	return res, nil
}

func (u *usecase) Delete(article models.Article) error {
	if err := u.pgRepository.Delete(article); err != nil {
		u.log.Errorf("article.pgRepository.Delete: %v", err)
		return err
	}

	if err := u.redisRepository.Delete(article.ID); err != nil {
		u.log.Errorf("article.redisRepository.Delete: %v", err)
		return err
	}

	return nil
}
