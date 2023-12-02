package usecases

import (
	"github.com/google/uuid"
	"github.com/slavtov/clean-architecture/internal/domain/models"
)

type ArticleUseCase interface {
	GetAll() ([]models.Article, error)
	GetByID(id uuid.UUID) (models.Article, error)
	Store(a *models.Article) (*models.Article, error)
	Update(a *models.Article) (*models.Article, error)
	Delete(a models.Article) error
}
