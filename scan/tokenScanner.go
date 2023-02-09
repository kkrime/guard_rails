package scan

import (
	"bufio"
	"io"
	"regexp"

	"guard_rails/client"
	"guard_rails/config"
	"guard_rails/model"
)

type tokenScanner struct {
	token  *regexp.Regexp
	config *config.TokenScannerConfig
}

func NewTokenScanner(tokenConfig *config.TokenScannerConfig) (RepositoryScanner, error) {
	tokenRegex, err := regexp.Compile(tokenConfig.ScanData.Token)
	if err != nil {
		return nil, err
	}

	return &tokenScanner{
		token:  tokenRegex,
		config: tokenConfig,
	}, nil
}

func (ts *tokenScanner) Scan(file client.File) *ScanResult {
	result := &ScanResult{
		Passed: true,
	}
	// result := &ScanResult{
	// 	Passed: false,
	// 	Findings: []model.Finding{
	// 		{
	// 			Type: "one",
	// 		},
	// 	},
	// }
	// return result

	fileReader, err := file.Reader()
	if err != nil {
		result.Err = err
		return nil
	}

	reader := bufio.NewReader(fileReader)

	var lineNumber int64
	for {
		lineNumber++

		line, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			result.Err = err
			return nil
		}

		if line == nil {
			break
		}

		if ts.token.Match(line) {
			// report finding
			finding := model.Finding{
				ScanData: &ts.config.ScanData,
				MetaData: &ts.config.MetaData,
				Location: model.Location{
					Path: file.Name(),
					Positions: model.Positions{
						Begin: model.Begin{
							Line: lineNumber,
						},
					},
				},
			}
			result.Findings = append(result.Findings, finding)
			result.Passed = false
		}
	}

	return result
}
