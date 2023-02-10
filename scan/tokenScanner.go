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

	return newTokenScanner(tokenRegex, tokenConfig), nil
}

func newTokenScanner(
	tokenRegex *regexp.Regexp,
	tokenConfig *config.TokenScannerConfig,
) *tokenScanner {
	return &tokenScanner{
		token:  tokenRegex,
		config: tokenConfig,
	}

}

func (ts *tokenScanner) Scan(file client.File) *ScanResult {
	result := &ScanResult{
		Passed: true,
	}

	fileReader, err := file.Reader()
	if err != nil {
		result.Err = err
		result.Passed = false
		return result
	}
	defer fileReader.Close()

	reader := bufio.NewReader(fileReader)

	var lineNumber int64
	for {
		lineNumber++

		line, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			result.Passed = false
			result.Err = err
			return result
		}

		if line == nil {
			break
		}

		if ts.token.Match(line) {
			// mark result as failed
			result.Passed = false

			// report finding
			finding := model.Finding{
				ScanData: ts.config.ScanData,
				MetaData: ts.config.MetaData,
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
		}
	}

	return result
}
