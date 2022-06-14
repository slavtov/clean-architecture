package repository

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/slavken/clean-architecture/internal/domain/models"
	"github.com/slavken/clean-architecture/internal/domain/repositories"
)

type pgRepository struct {
	db *sqlx.DB
}

func NewPGRepository(db *sqlx.DB) repositories.PGArticleRepository {
	return &pgRepository{db}
}

func (r *pgRepository) GetAll() ([]models.Article, error) {
	var articles []models.Article

	if err := r.db.Select(
		&articles,
		getArticlesQuery,
	); err != nil {
		return articles, echo.ErrInternalServerError
	}

	return articles, nil
}

func (r *pgRepository) GetByID(id uuid.UUID) (models.Article, error) {
	var article models.Article

	if err := r.db.Get(
		&article,
		getArticleQuery,
		id,
	); err != nil {
		if err == sql.ErrNoRows {
			return article, echo.ErrNotFound
		}

		return article, echo.ErrBadRequest
	}

	return article, nil
}

func (r *pgRepository) Store(a *models.Article) (*models.Article, error) {
	var article models.Article

	if err := r.db.QueryRowx(
		createArticleQuery,
		a.AuthorID,
		a.Title,
		a.Desc,
	).StructScan(&article); err != nil {
		return nil, echo.ErrBadRequest
	}

	return &article, nil
}

func (r *pgRepository) Update(a *models.Article) (*models.Article, error) {
	var article models.Article

	if err := r.db.QueryRowx(
		updateArticleQuery,
		a.Title,
		a.Desc,
		a.ID,
	).StructScan(&article); err != nil {
		if err == sql.ErrNoRows {
			return nil, echo.ErrNotFound
		}

		return nil, echo.ErrBadRequest
	}

	return &article, nil
}

func (r *pgRepository) Delete(a models.Article) error {
	res, err := r.db.Exec(
		deleteArticleQuery,
		a.ID,
		a.AuthorID,
	)
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
