package filesystem

import (
	"io/fs"
	"os"
	"path"
)

type embeddedFileSystem struct {
	fs      fs.FS
	dirname string
}

func NewEmbeddedFileSystem(fs fs.FS, dirname string) *embeddedFileSystem {
	return &embeddedFileSystem{
		fs:      fs,
		dirname: dirname,
	}
}

func (e *embeddedFileSystem) Open(name string) (fs.File, error) {
	file, err := e.fs.Open(path.Join(e.dirname, name))
	if err != nil {
		return nil, os.ErrNotExist
	}

	if fileInfo, err := os.Stat(name); err == nil {
		if fileInfo.IsDir() {
			return nil, os.ErrNotExist
		}
	}

	return file, nil
}
