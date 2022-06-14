package repositories

import (
	"time"

	"github.com/google/uuid"
	"github.com/slavken/clean-architecture/internal/domain/models"
)

type (
	PGUserRepository interface {
		GetAll() ([]models.User, error)
		GetByID(id uuid.UUID) (models.User, error)
		FindByEmail(email string) (models.User, error)
		Store(u *models.User) (*models.User, error)
		Update(u *models.User) (*models.User, error)
		Delete(id uuid.UUID) error
	}

	RedisUserRepository interface {
		GetByID(id uuid.UUID) (models.User, error)
		GetTokenInfo(id uuid.UUID, tokenID uuid.UUID) (uuid.UUID, error)
		SetToken(id uuid.UUID, tokenID uuid.UUID, exp int64) error
		SetUser(user *models.User, exp time.Duration) error
		Delete(keys ...string) error
		DeleteAll(pattern string) error
	}
)
