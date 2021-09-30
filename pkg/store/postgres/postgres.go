package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewClient(cfg *config) (*sqlx.DB, error) {
	dsn := fmt.Sprintf(
		`host=%s port=%d user=%s password=%s dbname=%s sslmode=%s`,
		cfg.host,
		cfg.port,
		cfg.user,
		cfg.password,
		cfg.name,
		cfg.ssl,
	)

	db, err := sqlx.Connect(cfg.driver, dsn)
	if err != nil {
		return nil, err
	}

	return db, nil
}
