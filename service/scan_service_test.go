package service

import (
	"bytes"
	"context"
	"fmt"
	"guard_rails/client"
	"guard_rails/config"
	"guard_rails/db"
	"guard_rails/logger"
	"guard_rails/model"
	"guard_rails/scan"
	"testing"

	clientMock "guard_rails/client/mocks"
	dbMock "guard_rails/db/mocks"

	scanMock "guard_rails/scan/mocks"

	mockk "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/tj/assert"
)

type ScanServiceTestSuite struct {
	suite.Suite
}

func TestScanServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ScanServiceTestSuite))
}

func (s *ScanServiceTestSuite) Testscan() {
	repositoryId := int64(222)
	scanId := int64(555)

	successResult := scan.ScanResult{
		Passed: true,
		Findings: model.Findings{
			{
				ScanData: &config.ScanData{Type: "passed"},
			},
		},
	}

	failedResult := scan.ScanResult{
		Passed: false,
		Findings: model.Findings{
			{
				ScanData: &config.ScanData{Type: "failed"},
			},
		},
	}

	repositoryDbObject := model.Repository{
		Id:   repositoryId,
		Name: "repository",
		Url:  "url",
	}

	scanDbObject := model.Scan{
		Id:           scanId,
		RepositoryId: repositoryId,
	}

	tests := []struct {
		name                  string
		gitClientProviderMock client.GitClientProvider
		repositoryDbMock      db.RepositoryDb
		scanDbMock            db.ScanDb
		scanners              []scan.RepositoryScanner
		findings              []model.Finding
		err                   error
		log                   []string
	}{
		{
			name: "Happy Path",

			repositoryDbMock: func() db.RepositoryDb {
				mock := dbMock.NewRepositoryDb(s.T())

				mock.On(
					"GetRepositoryById",
					mockk.Anything,
					repositoryId,
				).
					Return(
						func(
							ctx context.Context,
							repositoryId int64,
						) *model.Repository {
							return &repositoryDbObject
						},
						func(
							ctx context.Context,
							repositoryId int64,
						) error {
							return nil
						},
					).
					Once()

				return mock
			}(),

			scanDbMock: func() db.ScanDb {
				mock := dbMock.NewScanDb(s.T())

				mock.On(
					"StopScan",
					mockk.Anything,
					scanId,
					successResult.Findings,
					model.Success,
				).
					Return(
						func(
							ctx context.Context,
							scanId int64,
							findings model.Findings,
							status model.ScanStatus,
						) error {
							return nil
						},
					).
					Once()
				return mock
			}(),

			gitClientProviderMock: func() client.GitClientProvider {
				mock := clientMock.NewGitClientProvider(s.T())
				mock.On(
					"NewGitClient",
				).
					Return(
						func() client.GitClient {
							mock := clientMock.NewGitClient(s.T())

							// clone
							mock.On(
								"Clone",
								&repositoryDbObject,
							).
								Return(func(
									repository *model.Repository,
								) error {
									return nil
								}).
								Once()

								// first file
							mock.On(
								"GetNextFile",
							).
								Return(
									func() client.File {
										return func() client.File {
											mock := clientMock.NewFile(s.T())
											mock.On(
												"IsBinary",
											).
												Return(
													func() bool {
														return false
													},
													func() error {
														return nil
													},
												).
												Once()
											return mock
										}()
									},
									func() error {
										return nil
									},
								).
								Once()

								// end of files
							mock.On(
								"GetNextFile",
							).
								Return(
									func() client.File {
										return nil
									},
									func() error {
										return nil
									},
								).
								Once()

							return mock
						},
					).
					Once()

				return mock

			}(),

			scanners: []scan.RepositoryScanner{
				func() scan.RepositoryScanner {
					mock := scanMock.NewRepositoryScanner(s.T())

					mock.On(
						"Scan",
						mockk.Anything,
					).
						Return(
							func(
								file client.File,
							) *scan.ScanResult {
								return &successResult
							},
						).
						Once()

					return mock

				}(),
			},

			findings: successResult.Findings,

			log: []string{
				fmt.Sprintf(Scan_Completed_Successful, scanId),
			},
		},
		{
			name: "Scan Failed",

			repositoryDbMock: func() db.RepositoryDb {
				mock := dbMock.NewRepositoryDb(s.T())

				mock.On(
					"GetRepositoryById",
					mockk.Anything,
					repositoryId,
				).
					Return(
						func(
							ctx context.Context,
							repositoryId int64,
						) *model.Repository {
							return &repositoryDbObject
						},
						func(
							ctx context.Context,
							repositoryId int64,
						) error {
							return nil
						},
					).
					Once()

				return mock
			}(),

			scanDbMock: func() db.ScanDb {
				mock := dbMock.NewScanDb(s.T())

				mock.On(
					"StopScan",
					mockk.Anything,
					scanId,
					failedResult.Findings,
					model.Failure,
				).
					Return(
						func(
							ctx context.Context,
							scanId int64,
							findings model.Findings,
							status model.ScanStatus,
						) error {
							return nil
						},
					).
					Once()
				return mock
			}(),

			gitClientProviderMock: func() client.GitClientProvider {
				mock := clientMock.NewGitClientProvider(s.T())
				mock.On(
					"NewGitClient",
				).
					Return(
						func() client.GitClient {
							mock := clientMock.NewGitClient(s.T())

							// clone
							mock.On(
								"Clone",
								&repositoryDbObject,
							).
								Return(func(
									repository *model.Repository,
								) error {
									return nil
								}).
								Once()

								// first file
							mock.On(
								"GetNextFile",
							).
								Return(
									func() client.File {
										return func() client.File {
											mock := clientMock.NewFile(s.T())
											mock.On(
												"IsBinary",
											).
												Return(
													func() bool {
														return false
													},
													func() error {
														return nil
													},
												).
												Once()
											return mock
										}()
									},
									func() error {
										return nil
									},
								).
								Once()

								// end of files
							mock.On(
								"GetNextFile",
							).
								Return(
									func() client.File {
										return nil
									},
									func() error {
										return nil
									},
								).
								Once()

							return mock
						},
					).
					Once()

				return mock

			}(),

			scanners: []scan.RepositoryScanner{
				func() scan.RepositoryScanner {
					mock := scanMock.NewRepositoryScanner(s.T())

					mock.On(
						"Scan",
						mockk.Anything,
					).
						Return(
							func(
								file client.File,
							) *scan.ScanResult {
								return &failedResult
							},
						).
						Once()

					return mock

				}(),
			},

			findings: failedResult.Findings,

			log: []string{
				fmt.Sprintf(Scan_Completed_Successful, scanId),
			},
		},
		{
			// checks if the defer() works and sets the scan to failed in the DB
			name: "Panic",

			// to trigger the panic
			repositoryDbMock: nil,

			// checks scan results gets written to the db
			scanDbMock: func() db.ScanDb {
				mock := dbMock.NewScanDb(s.T())

				mock.On(
					"StopScan",
					mockk.Anything,
					scanId,
					model.Findings(nil),
					model.Failure,
				).
					Return(
						func(
							ctx context.Context,
							scanId int64,
							findings model.Findings,
							status model.ScanStatus,
						) error {
							return nil
						},
					).
					Once()
				return mock

			}(),

			log: []string{
				fmt.Sprintf(Scan_Paniced),
			},
		},
	}

	for _, test := range tests {
		s.T().Run(test.name, func(t2 *testing.T) {
			var logBuffer bytes.Buffer

			log := logger.CreateNewLogger()
			log.Out = &logBuffer

			entry := log.WithField("ScanID", scanId)

			scanService := newScannerService(test.repositoryDbMock, test.scanDbMock, test.gitClientProviderMock, nil, test.scanners)
			ssi := scanService.NewScanServiceInstance(entry)

			findings, err := ssi.scanRepository(&scanDbObject)
			assert.Equal(s.T(), test.findings, findings)
			assert.Equal(s.T(), test.err, err)

			for _, expectedLog := range test.log {
				assert.Contains(s.T(), logBuffer.String(), expectedLog)
			}
		})
	}
}
