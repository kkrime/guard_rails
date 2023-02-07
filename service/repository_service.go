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
	"github.com/sirupsen/logrus"
)

type repositoryService struct {
	httpClient   client.HttpClient
	repositoryDb db.RepositoryDb
}

func NewRepositoryServiceProvider(database *sqlx.DB) RepositoryServiceProvider {
	httpClient := client.NewHttpCleint()
	repositoryDb := db.NewDb(database)

	return &repositoryService{
		httpClient:   httpClient,
		repositoryDb: repositoryDb,
	}
}

func (rs *repositoryService) NewRepositoryServiceInstance(log *logrus.Entry) RepositoryService {

	return &repositoryServiceInstance{
		repositoryService: rs,
		log:               log,
	}
}

type repositoryServiceInstance struct {
	*repositoryService
	log *logrus.Entry
}

func (rsi *repositoryServiceInstance) AddRepository(ctx context.Context, repository *model.Repository) error {
	// check if repository url is valid
	reachable := rsi.httpClient.IsUrlReachable(repository.Url)

	if !reachable {
		return errors.NewRestError(424, Unable_To_Reach_Repository)
	}

	repository.Name = strings.ToLower(repository.Name)

	err := rsi.repositoryDb.AddRepository(ctx, repository)

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

func (rsi *repositoryServiceInstance) GetRepository(ctx context.Context, repositoryName string) (*model.Repository, error) {

	repository, err := rsi.repositoryDb.GetRepositoryByName(ctx, repositoryName)
	if err != nil {
		return nil, err
	}

	if repository == nil {
		return nil, errors.NewRestError(404, Repository_Not_Found)
	}

	return repository, err
}

func (rsi *repositoryServiceInstance) UpdateRepository(ctx context.Context, repository *model.Repository) error {
	// check if repository url is valid
	reachable := rsi.httpClient.IsUrlReachable(repository.Url)

	if !reachable {
		return errors.NewRestError(424, Unable_To_Reach_Repository)
	}

	rowsAffected, err := rsi.repositoryDb.UpdateRepository(ctx, repository)
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.NewRestError(404, Repository_Not_Found)
	}

	return nil
}

func (rsi *repositoryServiceInstance) DeleteRepository(ctx context.Context, repositoryName string) error {

	rowsAffected, err := rsi.repositoryDb.DeleteRepository(ctx, repositoryName)
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.NewRestError(404, Repository_Not_Found)
	}

	return nil
}
