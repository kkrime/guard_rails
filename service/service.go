package service

import (
	"context"
	"guard_rails/model"

	"github.com/sirupsen/logrus"
)

type RepositoryServiceProvider interface {
	NewRepositoryServiceInstance(log *logrus.Entry) RepositoryService
}

type RepositoryService interface {
	AddRepository(ctx context.Context, repository *model.Repository) error
	GetRepository(ctx context.Context, repositoryName string) (*model.Repository, error)
	UpdateRepository(ctx context.Context, repository *model.Repository) error
	DeleteRepository(ctx context.Context, repositoryName string) error
}

type ScanServiceProvider interface {
	NewScanServiceInstance(log *logrus.Entry) ScanService
}

type ScanService interface {
	Scan(ctx context.Context, repositoryName string) error
	GetScan(ctx context.Context, repositoryName string) ([]model.Scan, error)
}
