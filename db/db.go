package db

import (
	"context"
	"fmt"
	"guard_rails/config"
	"guard_rails/logger"
	"guard_rails/model"

	"github.com/jmoiron/sqlx"
	sqldblogger "github.com/simukti/sqldb-logger"
	"github.com/simukti/sqldb-logger/logadapter/logrusadapter"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/lib/pq"
)

// compiler time interface check
var _ RepositoryDb = (*db)(nil)
var _ ScanDb = (*db)(nil)

type db struct {
	db *sqlx.DB
}

func NewDbObject(database *sqlx.DB) *db {

	return &db{
		db: database,
	}
}

type RepositoryDb interface {
	AddRepository(ctx context.Context, repository *model.Repository) error
	GetRepositoryByName(ctx context.Context, repositoryName string) (*model.Repository, error)
	GetRepositoryById(ctx context.Context, repositoryId int64) (*model.Repository, error)
	UpdateRepository(ctx context.Context, repository *model.Repository) (int64, error)
	DeleteRepository(ctx context.Context, repositoryName string) (int64, error)
}
type ScanDb interface {
	GetScanWithStatus(ctx context.Context, repositoryId int64, status []model.ScanStatus) (*model.Scan, error)
	CreateNewScan(ctx context.Context, repositoryId int64) (*model.Scan, error)
	GetScans(ctx context.Context, repositoryName string) ([]model.Scan, error)
	UpdateScanStatus(ctx context.Context, scanId int64, status model.ScanStatus) error
	StartScan(ctx context.Context, scanId int64) error
	StopScan(ctx context.Context, scanId int64, findings model.Findings, status model.ScanStatus) error
}

func Init(config *config.PostgresConfig) (*sqlx.DB, error) {

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.Dbname)

	database, err := sqlx.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	dbLog := logger.CreateNewLogger()

	database.DB = sqldblogger.OpenDriver(dsn, database.DB.Driver(), logrusadapter.New(dbLog),
		sqldblogger.WithTimeFormat(sqldblogger.TimeFormatRFC3339),
		sqldblogger.WithLogDriverErrorSkip(true),
		sqldblogger.WithSQLQueryAsMessage(true))

	err = database.Ping()
	if err != nil {
		return nil, err
	}

	return database, nil
}
