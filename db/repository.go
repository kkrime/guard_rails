package db

import (
	"context"
	"guard_rails/model"

	"github.com/jmoiron/sqlx"
)

type repositoryDb struct {
	db *sqlx.DB
}

func NewRepositoryDb(db *sqlx.DB) RepositoryDb {

	return &repositoryDb{
		db: db,
	}
}

func (rd *repositoryDb) AddRepository(ctx context.Context, repository *model.Repository) error {

	statement := `
        INSERT INTO
            repositories
            (
                name,
                url
            )
            VALUES
            (
                $1,
                $2
            )
        ;`

	_, err := rd.db.ExecContext(ctx, statement, repository.Name, repository.Url)

	return err
}

func (rd *repositoryDb) GetRepository(ctx context.Context, repositoryName string) ([]model.Repository, error) {
	var repository []model.Repository

	statement := `
        SELECT * FROM
            repositories
        WHERE
            name = $1
        ;`

	err := rd.db.SelectContext(ctx, &repository, statement, repositoryName)

	return repository, err
}

func (rd *repositoryDb) UpdateRepository(ctx context.Context, repository *model.Repository) (int64, error) {
	var rowsAffected int64

	statement := `
        UPDATE 
            repositories
        SET
            url = $1
        WHERE
            name = $2
        ;`

	result, err := rd.db.ExecContext(ctx, statement, repository.Url, repository.Name)
	if err != nil {
		return 0, err
	}

	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

func (rd *repositoryDb) DeleteRepository(ctx context.Context, repositoryName string) (int64, error) {
	var rowsAffected int64

	statement := `
        DELETE FROM
            repositories
        WHERE
            name = $1
        ;`

	result, err := rd.db.ExecContext(ctx, statement, repositoryName)
	if err != nil {
		return 0, err
	}

	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}
