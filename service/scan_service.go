package service

import (
	"context"
	"sync"

	"guard_rails/client"
	"guard_rails/client/git"
	"guard_rails/config"
	"guard_rails/db"
	"guard_rails/errors"
	"guard_rails/logger"
	"guard_rails/model"
	"guard_rails/scan"
	scanner "guard_rails/scan"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type scanService struct {
	repositoryDb      db.RepositoryDb
	scanDb            db.ScanDb
	gitClientProvider client.GitClientProvider
	queue             chan *model.Scan
	scanners          []scan.RepositoryScanner
}

func NewScanServiceProvider(database *sqlx.DB, config *config.Config) (ScanServiceProvider, error) {
	var scanners []scan.RepositoryScanner

	for _, scanConfig := range config.TokenScanner {

		tokenScanner, err := scan.NewTokenScanner(&scanConfig)
		if err != nil {
			return nil, err
		}

		scanners = append(scanners, tokenScanner)
	}

	db := db.NewDbObject(database)
	gitClientProvider := git.NewGitCleintProvider(&config.Git)

	queue := make(chan *model.Scan, config.Queue.QueueSize)

	scanService := &scanService{
		repositoryDb:      db,
		scanDb:            db,
		gitClientProvider: gitClientProvider,
		queue:             queue,
		scanners:          scanners,
	}

	go scanService.scan()

	return scanService, nil
}

func newScannerService(

	repositoryDb db.RepositoryDb,
	scanDb db.ScanDb,
	gitClientProvider client.GitClientProvider,
	queue chan *model.Scan,
	scanners []scan.RepositoryScanner) *scanService {

	return &scanService{
		repositoryDb:      repositoryDb,
		scanDb:            scanDb,
		gitClientProvider: gitClientProvider,
		queue:             queue,
		scanners:          scanners,
	}

}

func (ss *scanService) NewScanServiceInstance(log *logrus.Entry) ScanService {
	return &scanServiceInstance{scanService: ss, log: log}
}

type scanServiceInstance struct {
	*scanService
	log *logrus.Entry
}

func (ssi *scanServiceInstance) QueueScan(ctx context.Context, repositoryName string) error {

	// check if repository exists
	repository, err := ssi.repositoryDb.GetRepositoryByName(ctx, repositoryName)
	if err != nil {
		ssi.log.Errorf("unable to get repository from db, err %v", err)
		return err
	}

	if repository == nil {
		return errors.NewRestError(404, Repository_Not_Found)
	}

	repositoryId := repository.Id

	// check if scan already exists and queued/in progress
	scans, err := ssi.scanDb.GetScanWithStatus(ctx, repositoryId, []model.ScanStatus{model.Queued, model.InProgress})
	if err != nil {
		ssi.log.Errorf("unable to get scans from database, err %v", err)
		return err
	}

	if scans != nil {
		ssi.log.Errorf("scans for repository %v already exists", repositoryName)
		errors.NewRestError(409, Scan_Already_Exists)
	}

	// create new scan
	scan, err := ssi.scanDb.CreateNewScan(ctx, repositoryId)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			// race conditon
			if pgErr.Code != db.Duplicate_Record_Code {
				ssi.log.Errorf("unable to create new scan, err %v", err)
				return err
			}
		} else {
			ssi.log.Errorf("unable to create new scan, err %v", err)
			return err
		}
	}

	go func() {
		ssi.queue <- scan
		ssi.log.Infof("scan qeueud for repository %v", repository.Name)
	}()

	return nil
}

func (ssi *scanServiceInstance) GetScan(ctx context.Context, repositoryName string) ([]model.Scan, error) {

	// check if repository exists
	scans, err := ssi.scanDb.GetScans(ctx, repositoryName)
	if err != nil {
		ssi.log.Errorf("unable to get scans from database, err %v", err)
		return nil, err
	}

	if scans == nil {
		ssi.log.Errorf("repository %v not found in the database, err %v", repositoryName, err)
		return nil, errors.NewRestError(404, Repository_Not_Found)
	}

	return scans, nil

}

func (ss *scanService) scan() {

	for scan := range ss.queue {

		//TODO create scan context

		log := logger.CreateNewLogger()
		log.ReportCaller = true
		entry := log.WithField("ScanID", scan.Id)
		ssi := ss.NewScanServiceInstance(entry)

		// mark repository as in progress
		err := ss.scanDb.StartScan(context.Background(), scan.Id)
		if err != nil {
			log.Errorf("unable to set scan as in progress, err %v", err)
			continue
		}

		log.Infof("scan %v started", scan.Id)

		go ssi.scanRepository(scan)
	}
}

func (ssi *scanServiceInstance) scanRepository(scan *model.Scan) ( /* return parameters added for testing purposes*/ findings []model.Finding, err error) {
	var (
		repository *model.Repository
		file       client.File
	)

	status := model.Success

	resultChann := make(chan *scanner.ScanResult)

	// on error or panic mark the scan as FAILURE
	defer func() {
		paniced := recover()
		if err != nil || paniced != nil {
			if err != nil {
				ssi.log.Errorf("scan errored, err %v", err)
			} else {
				ssi.log.Error(Scan_Paniced)
			}
			err = ssi.scanDb.StopScan(context.Background(), scan.Id, nil, model.Failure)
			if err != nil {
				ssi.log.Errorf("unable to set scan %v as failed, err %v", scan.Id, err)
			}
		}
	}()

	repository, err = ssi.repositoryDb.GetRepositoryById(context.Background(), scan.RepositoryId)
	if err != nil {
		ssi.log.Errorf("unable to get repository id %v from the database, err %v", scan.RepositoryId, err)
		return nil, err
	}

	// very unlikey to happen, but just in case...
	if repository == nil {
		ssi.log.Errorf("repository id %v not found in the database, err %v", repository.Id, err)
		return nil, err
	}

	gitClient := ssi.gitClientProvider.NewGitClient()
	ssi.log.Infof("cloneing repo %v", repository.Name)
	err = gitClient.Clone(repository)
	if err != nil {
		ssi.log.Errorf("unable to clone repository %v, err %v", repository.Name, err)
		return nil, err
	}
	ssi.log.Infof("finished cloneing repo %v", repository.Name)

	var ScannersWg sync.WaitGroup

	for {
		var isBinary bool

		file, err = gitClient.GetNextFile()
		if err != nil {
			ssi.log.Errorf("unable to get next file  err %v", err)
			return nil, err
		}

		if file == nil {
			// end of the repository
			break
		}

		isBinary, err = file.IsBinary()
		if err != nil {
			ssi.log.Errorf("error on isBinary() err %v", err)
			break
		}

		// skip if file is binary
		if isBinary {
			continue
		}

		ScannersWg.Add(1)

		// each file will be scanned in its own goroutine
		go func(file client.File) {
			defer ScannersWg.Done()
			for _, scanner := range ssi.scanners {
				result := scanner.Scan(file)
				resultChann <- result
			}
		}(file)
	}

	var readingResultsWg sync.WaitGroup
	readingResultsWg.Add(1)
	go func() {
		defer readingResultsWg.Done()

		for result := range resultChann {
			if result.Err != nil {
				err = result.Err
				// drain resultChann
				for len(resultChann) != 0 {
					<-resultChann
				}
				return
			}

			if !result.Passed {
				status = model.Failure
			}

			findings = append(findings, result.Findings...)
		}
	}()

	ScannersWg.Wait()
	close(resultChann)
	readingResultsWg.Wait()

	writeResultsErr := ssi.scanDb.StopScan(context.Background(), scan.Id, findings, status)
	if writeResultsErr != nil {
		ssi.log.Errorf("unable to get mark scan %v as successful err %v", scan.Id, err)
	}

	ssi.log.Infof(Scan_Completed_Successful, scan.Id)

	return findings, nil
}
