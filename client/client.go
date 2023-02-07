package client

import (
	"guard_rails/model"
	"io"
)

type HttpClient interface {
	IsUrlReachable(url string) bool
}

type File interface {
	Size() int64
	Name() string
	Reader() (io.ReadCloser, error)
}

type GitClientProvider interface {
	NewGitClient() GitClient
}

type GitClient interface {
	Clone(repository *model.Repository) error
	GetNextFile() (File, error)
}
