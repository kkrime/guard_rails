package db

import (
	"context"
	"fmt"
	"guard_rails/model"
)

func (rd *db) AddRepository(ctx context.Context, repository *model.Repository) error {

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

func (rd *db) GetRepositoryByName(ctx context.Context, repositoryName string) (*model.Repository, error) {
	var repository []model.Repository

	statement := `
        SELECT * FROM
            repositories
        WHERE
            name = $1 AND
            deleted_at IS null
        ;`

	err := rd.db.SelectContext(ctx, &repository, statement, repositoryName)
	fmt.Printf("err = %+v\n", err)

	if err != nil {
		return nil, err
	}

	if repository == nil {
		return nil, nil
	}

	return &repository[0], err
}

func (rd *db) GetRepositoryById(repositoryId int64) (*model.Repository, error) {
	var repository []model.Repository

	statement := `
        SELECT * FROM
            repositories
        WHERE
            id = $1 AND
            deleted_at IS null
        ;`

	err := rd.db.Select(&repository, statement, repositoryId)

	if err != nil {
		return nil, err
	}

	if repository == nil {
		return nil, nil
	}

	return &repository[0], err
}

func (rd *db) UpdateRepository(ctx context.Context, repository *model.Repository) (int64, error) {
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

func (rd *db) DeleteRepository(ctx context.Context, repositoryName string) (int64, error) {
	var rowsAffected int64

	statement := `
        UPDATE
            repositories
        SET
            deleted_at = now()
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
