package service

import (
	"context"
	"guard_rails/model"
)

type RepositoryService interface {
	AddRepository(ctx context.Context, repository *model.Repository) error
	GetRepository(ctx context.Context, repositoryName string) (*model.Repository, error)
	UpdateRepository(ctx context.Context, repository *model.Repository) error
	DeleteRepository(ctx context.Context, repositoryName string) error
}
