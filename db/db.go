package db

import (
	"context"
	"fmt"
	"guard_rails/model"

	"github.com/jmoiron/sqlx"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/lib/pq"
)

type RepositoryDb interface {
	AddRepository(ctx context.Context, repository *model.Repository) error
	GetRepository(ctx context.Context, repositoryName string) ([]model.Repository, error)
	UpdateRepository(ctx context.Context, repository *model.Repository) (int64, error)
	DeleteRepository(ctx context.Context, repositoryName string) (int64, error)
}

func Init() (*sqlx.DB, error) {
	const (
		host     = "localhost"
		port     = 5432
		user     = "postgres"
		password = "postgres"
		dbname   = "guard_rails"
	)

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sqlx.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
