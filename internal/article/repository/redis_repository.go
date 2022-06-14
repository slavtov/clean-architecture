package repository

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/slavken/clean-architecture/internal/domain/models"
	"github.com/slavken/clean-architecture/internal/domain/repositories"
	"github.com/slavken/clean-architecture/pkg/store/redis"
	"github.com/slavken/clean-architecture/pkg/utils"
)

type redisRepository struct {
	redis redis.Store
}

const prefix = "articles"

func NewRedisRepository(rdb redis.Store) repositories.RedisArticleRepository {
	return &redisRepository{rdb}
}

func (r *redisRepository) GetByID(id uuid.UUID) (models.Article, error) {
	var article models.Article

	res, err := r.redis.Get(utils.GetRedisKey(prefix, id.String()))
	if err != nil {
		return article, echo.ErrNotFound
	}

	if err := json.Unmarshal([]byte(res), &article); err != nil {
		return article, echo.ErrInternalServerError
	}

	return article, nil
}

func (r *redisRepository) SetArticle(
	article *models.Article,
	exp time.Duration,
) error {
	res, err := json.Marshal(article)
	if err != nil {
		return echo.ErrInternalServerError
	}

	if err := r.redis.Set(utils.GetRedisKey(
		prefix,
		article.ID.String(),
	), res, exp); err != nil {
		return echo.ErrInternalServerError
	}

	return nil
}

func (r *redisRepository) Delete(id uuid.UUID) error {
	if err := r.redis.Del(utils.GetRedisKey(
		prefix,
		id.String(),
	)); err != nil {
		return echo.ErrNotFound
	}

	return nil
}
