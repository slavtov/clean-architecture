package repository

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/slavtov/clean-architecture/internal/domain/models"
	"github.com/slavtov/clean-architecture/internal/domain/repositories"
	"github.com/slavtov/clean-architecture/pkg/store/redis"
	"github.com/slavtov/clean-architecture/pkg/utils"
)

type redisRepository struct {
	redis redis.Store
}

const (
	authPrefix = "auth"
	userPrefix = "users"
)

func NewRedisRepository(rdb redis.Store) repositories.RedisUserRepository {
	return &redisRepository{rdb}
}

func (r *redisRepository) GetByID(id uuid.UUID) (models.User, error) {
	var user models.User

	res, err := r.redis.Get(utils.GetRedisKey(userPrefix, id.String()))
	if err != nil {
		return user, echo.ErrNotFound
	}

	if err = json.Unmarshal([]byte(res), &user); err != nil {
		return user, echo.ErrInternalServerError
	}

	return user, nil
}

func (r *redisRepository) GetTokenInfo(
	id uuid.UUID,
	tokenID uuid.UUID,
) (uuid.UUID, error) {
	res, err := r.redis.Get(utils.GetRedisKey(
		authPrefix,
		id.String(),
		tokenID.String(),
	))
	if err != nil {
		return uuid.Nil, echo.ErrNotFound
	}

	return uuid.Parse(res)
}

func (r *redisRepository) SetToken(
	id uuid.UUID,
	tokenID uuid.UUID,
	exp int64,
) error {
	t := time.Unix(exp, 0)
	now := time.Now()

	if err := r.redis.Set(utils.GetRedisKey(
		authPrefix,
		id.String(),
		tokenID.String(),
	), id.String(), t.Sub(now)); err != nil {
		return echo.ErrInternalServerError
	}

	return nil
}

func (r *redisRepository) SetUser(
	user *models.User,
	exp time.Duration,
) error {
	res, err := json.Marshal(user)
	if err != nil {
		return echo.ErrInternalServerError
	}

	if err = r.redis.Set(utils.GetRedisKey(
		userPrefix,
		user.ID.String(),
	), res, exp); err != nil {
		return echo.ErrInternalServerError
	}

	return nil
}

func (r *redisRepository) Delete(keys ...string) error {
	if err := r.redis.Del(keys...); err != nil {
		return echo.ErrNotFound
	}

	return nil
}

func (r *redisRepository) DeleteAll(pattern string) error {
	if err := r.redis.DelAll(pattern); err != nil {
		return echo.ErrNotFound
	}

	return nil
}
