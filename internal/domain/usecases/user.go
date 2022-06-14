package usecases

import (
	"github.com/google/uuid"
	"github.com/slavken/clean-architecture/internal/domain/models"
	"github.com/slavken/clean-architecture/pkg/utils"
)

type (
	jwtUseCase interface {
		Auth(user *models.User) (*models.AuthUser, error)
		GetToken(id uuid.UUID, tokenID uuid.UUID) (uuid.UUID, error)
		DeleteToken(id uuid.UUID, tokenID uuid.UUID) error
		Logout(id uuid.UUID, tokenID *utils.TokenDetails) error
		LogoutAll(id uuid.UUID) error
	}

	UserUseCase interface {
		GetAll() ([]models.User, error)
		GetByID(id uuid.UUID) (models.User, error)
		Login(user *models.User) (*models.AuthUser, error)
		Store(user *models.User) (*models.AuthUser, error)
		Update(user *models.User) (*models.User, error)
		Delete(id uuid.UUID) error
		jwtUseCase
	}
)
