package repositories

import (
	"time"

	"github.com/google/uuid"
	"github.com/slavken/go-clean-architecture/internal/domain/models"
)

type (
	PGArticleRepository interface {
		GetAll() ([]models.Article, error)
		GetByID(id uuid.UUID) (models.Article, error)
		Store(a *models.Article) (*models.Article, error)
		Update(a *models.Article) (*models.Article, error)
		Delete(a models.Article) error
	}

	RedisArticleRepository interface {
		GetByID(id uuid.UUID) (models.Article, error)
		SetArticle(article *models.Article, exp time.Duration) error
		Delete(id uuid.UUID) error
	}
)
