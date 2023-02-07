package scan

import (
	"bufio"
	"io"
	"regexp"

	"guard_rails/client"
	"guard_rails/model"
)

type tokenScanner struct {
	token *regexp.Regexp
}

func NewTokenScanner(token string) (RepositoryScanner, error) {
	tokenRegex, err := regexp.Compile(`zzzzzzzzzzzzzzzzzzzzzzz`)
	if err != nil {
		return nil, err
	}

	return &tokenScanner{
		token: tokenRegex,
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

	lineNumber := 0
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

		// fmt.Printf("string(line) = %+v\n", string(line))

		if ts.token.Match(line) {
			finding := model.Finding{
				Type: "ddd",
			}
			result.Findings = append(result.Findings, finding)
			result.Passed = false
		}
	}

	return result
}
