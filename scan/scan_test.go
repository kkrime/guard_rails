package scan

import (
	"fmt"
	"guard_rails/client"
	"guard_rails/config"
	"guard_rails/model"
	"io"
	"regexp"
	"strings"

	"testing"

	clientMock "guard_rails/client/mocks"

	"github.com/stretchr/testify/suite"
	"github.com/tj/assert"
)

type TokenScannerTestSuite struct {
	suite.Suite
}

func TestTokenScannerTestSuite(t *testing.T) {
	suite.Run(t, new(TokenScannerTestSuite))
}

func (s *TokenScannerTestSuite) Testscan() {

	token := `private_key\w`
	fileName := "test_file.go"
	scanType := "test_type"
	ruleId := "test_rule"

	scanData := config.ScanData{
		Type:   scanType,
		RuleId: ruleId,
	}

	metaData := config.MetaData{}

	err := fmt.Errorf("tokenScanner error")

	tests := []struct {
		name        string
		file        client.File
		token       string
		scanResult  ScanResult
		err         error
		msgOut      string
		expectedLog []string
	}{
		{
			name: "Happy Path - token found",

			file: func() client.File {
				mock := clientMock.NewFile(s.T())

				mock.On(
					"Reader",
				).
					Return(
						func() io.ReadCloser {
							file := "this is line 1\nthis is line 2, this line contains the token: private_keyPRIVATEKEY"
							fileReader := strings.NewReader(file)
							return io.NopCloser(fileReader)
						},
						func() error {
							return nil
						},
					).
					Once()

				mock.On(
					"Name",
				).
					Return(
						func() string {
							return fileName
						},
					).
					Once()

				return mock
			}(),

			token: token,

			scanResult: ScanResult{
				Passed: false,
				Findings: []model.Finding{
					{
						ScanData: &scanData,
						Location: model.Location{
							Path: fileName,
							Positions: model.Positions{
								Begin: model.Begin{
									Line: 2,
								},
							},
						},
						MetaData: &metaData,
					},
				},
			},
		},
		{
			name: "Happy Path - token not found",

			file: func() client.File {
				mock := clientMock.NewFile(s.T())

				mock.On(
					"Reader",
				).
					Return(
						func() io.ReadCloser {
							// INFO: Just the prefix, NOTHING AFTER, so regex should not match
							file := "this is line 1\nthis is line 2, this line does not contains the token: private_key"
							fileReader := strings.NewReader(file)
							return io.NopCloser(fileReader)
						},
						func() error {
							return nil
						},
					).
					Once()

				return mock
			}(),

			token: token,

			scanResult: ScanResult{
				Passed: true,
			},
		},
		{
			name: "Error on Reader()",

			file: func() client.File {
				mock := clientMock.NewFile(s.T())

				mock.On(
					"Reader",
				).
					Return(
						func() io.ReadCloser {
							return nil
						},
						func() error {
							return err
						},
					).
					Once()

				return mock
			}(),

			token: token,

			scanResult: ScanResult{
				Passed: false,
				Err:    err,
			},
		},
	}

	for _, test := range tests {
		s.T().Run(test.name, func(t2 *testing.T) {

			tokenRegex, err := regexp.Compile(test.token)
			if err != nil {
				s.Fail("unable to compile regex")
			}

			config := config.TokenScannerConfig{
				ScanData: &scanData,
				MetaData: &metaData,
			}

			tokenScanner := newTokenScanner(tokenRegex, &config)
			scanResult := tokenScanner.Scan(test.file)

			assert.Equal(s.T(), &test.scanResult, scanResult)
		})
	}
}
