package service

import (
	"context"
	"guard_rails/client"
	"guard_rails/db"
	"guard_rails/model"
	"sdk/support"
	"testing"
	"time"

	"guard_rails/client/git"
	clientMock "guard_rails/client/mock"
	"guard_rails/db"
	scanDbMock "guard_rails/db/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ScanServiceTestSuite struct {
	suite.Suite
}

var (
	timezoneString string = "Asia/Singapore"
	loc, _                = time.LoadLocation(timezoneString)
)

func (s *ScanServiceTestSuite) Testscan() {
	scanId := 555

	successFindings := &ScanResult{
		Passed: true,
		Findings: []model.Finding{
			{
				Type: "one",
			},
		},
	}

	repository := model.Repository{
		Name: "repository",
		Url:  "url",
	}

	tests := []struct {
		name                  string
		gitClientProviderMock git.GitClientProvider
		repositoryDbMock      db.RepositoryDb
		scanDbMock            db.ScanDb
		msgOut                string
		expectedLog           []string
	}{
		{
			name: "Happy Path",

			scanDbMock: func() db.ScanDb {
				mock := scanDbMock.NewScanDbb(s.T())

				mock.On(
					"StopScan",
					scanId,
					successFindings.Findings,
					model.Success,
				).
					Return(
						func(
							scanId int64,
							findings model.Findings,
							status model.ScanStatus,
						) error {
							return nil
						},
					).
					Once()
				return mock
			},

			gitClientProviderMock: func() client.GitClientProvider {
				mock := clientMock.NewGitCleintProvider(s.T())
				mock.On(
					"NewGitClient",
				).
					Return(
						func() client.GitClient {
							mock := clientMock.NewGitClient(s.T())

							// clone
							mock.On(
								"CLone",
								repository,
							).
								Return(
									func()).
								Once()

								// first file
							mock.On(
								"GetNextFile",
							).
								Return(
									func() client.File {
										return &git.File{}
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
								).
								Once()

							return mock
						},
					).
					Once()

				return mock

			},
		},
	}

	for _, test := range tests {
		s.T().Run(test.name, func(t2 *testing.T) {
			log, logBuffer, logWriter := support.NewTestLogger()

			out, err := ls.updateTransactionLimit(context.TODO(), nil, nil, &test.transactionLimit)
			logWriter.Flush()

			assert.Equal(s.T(), test.msgOut, out.Message)
			assert.Equal(s.T(), nil, err)

			for _, expectedLog := range test.expectedLog {
				assert.Contains(s.T(), logBuffer.String(), expectedLog)
			}
		})
	}
}
