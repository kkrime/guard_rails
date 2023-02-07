package service

import (
	"context"
	"fmt"
	"sync"

	"guard_rails/client"
	"guard_rails/client/git"
	"guard_rails/db"
	"guard_rails/errors"
	"guard_rails/model"
	"guard_rails/scan"
	scanner "guard_rails/scan"

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

func NewScanServiceProvider(database *sqlx.DB) ScanServiceProvider {
	db := db.NewDb(database)
	queue := make(chan *model.Scan, 4098)

	tokenScanner, err := scan.NewTokenScanner("token")
	if err != nil {
		// return err
	}

	scanners := []scan.RepositoryScanner{
		tokenScanner,
	}

	gitClientProvider := git.NewGitCleintProvider()

	scanService := &scanService{
		repositoryDb:      db,
		scanDb:            db,
		gitClientProvider: gitClientProvider,
		queue:             queue,
		scanners:          scanners,
	}

	go scanService.scan()

	return scanService
}

func (ss *scanService) NewScanServiceInstance(log *logrus.Entry) ScanService {

	return &scanServiceInstance{
		scanService: ss,
		log:         log,
	}
}

type scanServiceInstance struct {
	*scanService
	log *logrus.Entry
}

func (ss *scanServiceInstance) Scan(ctx context.Context, repositoryName string) error {

	// check if repository exists
	repository, err := ss.repositoryDb.GetRepositoryByName(ctx, repositoryName)
	if err != nil {
		return err
	}

	if repository == nil {
		return errors.NewRestError(404, Repository_Not_Found)
	}

	repositoryId := repository.Id

	// check if scan already exists and queued/in progress
	scans, err := ss.scanDb.GetScanWithStatus(ctx, repositoryId, []model.ScanStatus{model.Queued, model.InProgress})
	fmt.Println("AFTER")
	if err != nil {
		fmt.Printf("err = %+v\n", err)
		return err
	}
	fmt.Println("AFTER")

	if scans != nil {
		errors.NewRestError(409, Scan_Already_Exists)
	}

	// create new scan
	scan, err := ss.scanDb.CreateNewScan(ctx, repositoryId)
	if err != nil {
		return err
	}

	go func() {
		// only sending repsitory id and not the repository object because the user is able to update the repository
		ss.queue <- scan
	}()

	return nil
}

func (ss *scanServiceInstance) GetScan(ctx context.Context, repositoryName string) ([]model.Scan, error) {

	// check if repository exists
	scans, err := ss.scanDb.GetScans(ctx, repositoryName)
	if err != nil {
		return nil, err
	}

	if scans == nil {
		return nil, errors.NewRestError(404, Repository_Not_Found)
	}

	return scans, nil

}

func (ss *scanService) scan() {

	for scan := range ss.queue {

		// mark repository as in progress
		err := ss.scanDb.StartScan(scan.Id)
		if err != nil {
			// TODO log
			continue
		}

		go ss.scanRepository(scan)

	}
}

func (ss *scanService) scanRepository(scan *model.Scan) (err error) {
	var (
		repository *model.Repository
		file       client.File
		findings   model.Findings
	)

	status := model.Success

	resultChann := make(chan *scanner.ScanResult)

	defer func() {
		// on error mark the scan as FAILURE
		if err != nil {
			// TODO log
			ss.scanDb.StopScan(scan.Id, nil, model.Failure)
		}
	}()

	repository, err = ss.repositoryDb.GetRepositoryById(scan.RepositoryId)
	if err != nil {
		// TODO log
		return err
	}

	// very unlikey to happen, but just in case...
	if repository == nil {
		// TODO log
		return err
	}

	gitClient := ss.gitClientProvider.NewGitClient()
	err = gitClient.Clone(repository)
	if err != nil {
		return err
	}

	var ScannersWg sync.WaitGroup

	for {

		file, err = gitClient.GetNextFile()
		if err != nil {
			// TODO log
			return err
		}

		if file == nil {
			// TODO log
			break
		}

		ScannersWg.Add(1)

		// each file will be scanned in its own goroutine
		go func(file client.File) {
			defer ScannersWg.Done()
			for _, scanner := range ss.scanners {
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
	fmt.Println(">>>>>>>>>>>>>>>>>>. 1")
	close(resultChann)
	fmt.Println(">>>>>>>>>>>>>>>>>>. 2")
	readingResultsWg.Wait()
	fmt.Println(">>>>>>>>>>>>>>>>>>. 3")

	writeResultsErr := ss.scanDb.StopScan(scan.Id, findings, status)
	if writeResultsErr != nil {
		// TOOD log
	}

	return nil
}
