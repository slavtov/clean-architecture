package repository

import (
	"database/sql"
	"net/http"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	"github.com/slavken/clean-architecture/internal/domain/models"
	"github.com/slavken/clean-architecture/internal/domain/repositories"
)

type pgRepository struct {
	db *sqlx.DB
}

func NewPGRepository(db *sqlx.DB) repositories.PGUserRepository {
	return &pgRepository{db}
}

func (r *pgRepository) GetAll() ([]models.User, error) {
	var users []models.User

	if err := r.db.Select(
		&users,
		getUsersQuery,
	); err != nil {
		return users, echo.ErrInternalServerError
	}

	return users, nil
}

func (r *pgRepository) GetByID(id uuid.UUID) (models.User, error) {
	var user models.User

	if err := r.db.Get(
		&user,
		getUserQuery,
		id,
	); err != nil {
		if err == sql.ErrNoRows {
			return user, echo.ErrNotFound
		}

		return user, echo.ErrBadRequest
	}

	return user, nil
}

func (r *pgRepository) FindByEmail(email string) (models.User, error) {
	var user models.User

	if err := r.db.QueryRowx(
		findUserByEmailQuery,
		email,
	).StructScan(&user); err != nil {
		if err == sql.ErrNoRows {
			return user, echo.NewHTTPError(
				http.StatusBadRequest,
				"user is not found",
			)
		}

		return user, echo.ErrBadRequest
	}

	return user, nil
}

func (r *pgRepository) Store(u *models.User) (*models.User, error) {
	var user models.User

	if err := r.db.QueryRowx(
		createUserQuery,
		u.Email,
		u.Password,
	).StructScan(&user); err != nil {
		if err.(*pq.Error).Code == "23505" {
			return nil, echo.NewHTTPError(
				http.StatusBadRequest,
				"email already exists",
			)
		}

		return nil, echo.ErrBadRequest
	}

	return &user, nil
}

func (r *pgRepository) Update(a *models.User) (*models.User, error) {
	var user models.User

	if err := r.db.QueryRowx(
		updateUserQuery,
		a.Email,
		a.Password,
		a.ID,
	).StructScan(&user); err != nil {
		if err == sql.ErrNoRows {
			return nil, echo.ErrNotFound
		}

		return nil, echo.ErrBadRequest
	}

	return &user, nil
}

func (r *pgRepository) Delete(id uuid.UUID) error {
	res, err := r.db.Exec(deleteUserQuery, id)
	if err != nil {
		return echo.ErrBadRequest
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return echo.ErrInternalServerError
	}

	if rowsAffected == 0 {
		return echo.ErrNotFound
	}

	return nil
}
