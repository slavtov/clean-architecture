package models

import (
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type Article struct {
	ID        uuid.UUID `json:"id" db:"id" example:"00000000-0000-0000-0000-000000000000"`
	AuthorID  uuid.UUID `json:"author_id" db:"author_id" validate:"required" example:"00000000-0000-0000-0000-000000000000"`
	Title     string    `json:"title" db:"title" validate:"required,min=5,max=250" example:"Title"`
	Desc      string    `json:"desc" db:"desc" validate:"required" example:"Description"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at" example:"0000-01-01T00:00:00.000000Z"`
	CreatedAt time.Time `json:"created_at" db:"created_at" example:"0000-01-01T00:00:00.000000Z"`
}

type ArticlesList struct {
	TotalCount int       `json:"total_count"`
	Articles   []Article `json:"articles"`
}

func (a *Article) Validate() error {
	validate := validator.New()

	a.Title = strings.TrimSpace(a.Title)
	a.Desc = strings.TrimSpace(a.Desc)

	return validate.Struct(a)
}
