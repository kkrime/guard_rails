package scan

import (
	"guard_rails/client"
	"guard_rails/model"
)

type ScanResult struct {
	Passed   bool
	Findings model.Findings
	Err      error
}

type RepositoryScanner interface {
	Scan(file client.File) *ScanResult
}
