package git

import (
	"io"

	"guard_rails/client"

	"github.com/go-git/go-git/v5/plumbing/object"
)

type File struct {
	file *object.File
}

func newFile(fileObject *object.File) client.File {
	return &File{
		file: fileObject,
	}
}

func (f *File) Size() int64 {
	return f.file.Size
}

func (f *File) Name() string {
	return f.file.Name
}

func (f *File) Reader() (io.ReadCloser, error) {
	return f.file.Reader()
}

func (f *File) IsBinary() (bool, error) {
	return f.file.IsBinary()
}
