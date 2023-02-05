package service

import (
	"context"
	"guard_rails/client"
	"guard_rails/db"
	"guard_rails/model"
	"strings"

	"guard_rails/errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
)

type repositoryService struct {
	httpClient   client.HttpClient
	repositoryDb db.RepositoryDb
}

func NewRepositoryService(database *sqlx.DB) RepositoryService {
	httpClient := client.NewHttpCleint()
	repositoryDb := db.NewRepositoryDb(database)

	return &repositoryService{
		httpClient:   httpClient,
		repositoryDb: repositoryDb,
	}
}

func (rs *repositoryService) AddRepository(ctx context.Context, repository *model.Repository) error {
	// check if repository url is valid
	reachable := rs.httpClient.IsUrlReachable(repository.Url)

	if !reachable {
		return errors.NewRestError(424, Unable_To_Reach_Repository)
	}

	repository.Name = strings.ToLower(repository.Name)

	err := rs.repositoryDb.AddRepository(ctx, repository)

	if err != nil {
		// check if duplicate add
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == db.Duplicate_Record_Code {
				err = errors.NewRestError(409, Repository_Already_Added)
			}
		}
	}

	return err
}

func (rs *repositoryService) GetRepository(ctx context.Context, repositoryName string) (*model.Repository, error) {

	repository, err := rs.repositoryDb.GetRepository(ctx, repositoryName)
	if err != nil {
		return nil, err
	}

	if repository == nil {
		return nil, errors.NewRestError(404, Repository_Not_Found)
	}

	return &repository[0], err
}

func (rs *repositoryService) UpdateRepository(ctx context.Context, repository *model.Repository) error {
	// check if repository url is valid
	reachable := rs.httpClient.IsUrlReachable(repository.Url)

	if !reachable {
		return errors.NewRestError(424, Unable_To_Reach_Repository)
	}

	rowsAffected, err := rs.repositoryDb.UpdateRepository(ctx, repository)
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.NewRestError(404, Repository_Not_Found)
	}

	return nil
}

func (rs *repositoryService) DeleteRepository(ctx context.Context, repositoryName string) error {

	rowsAffected, err := rs.repositoryDb.DeleteRepository(ctx, repositoryName)
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.NewRestError(404, Repository_Not_Found)
	}

	return nil
}
