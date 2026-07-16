package storage

import (
	"cloud/internal/config"
	"cloud/internal/models"
	"database/sql"

	_ "github.com/lib/pq"
)

func Connect(cfg *config.Config) (*models.Postgres, error) {
	db, err := sql.Open("postgres", cfg.DSN)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &models.Postgres{Db: db}, nil
}
