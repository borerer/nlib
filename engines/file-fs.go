package engines

import (
	"io"
	"os"
	"path"
)

type FileFS struct {
	BaseDir string
}

func NewFileFS(baseDir string) *FileFS {
	return &FileFS{
		BaseDir: baseDir,
	}
}

func (fileFS *FileFS) GetFile(filename string) (io.Reader, error) {
	absFilePath := path.Join(fileFS.BaseDir, filename)
	f, err := os.Open(absFilePath)
	if err != nil {
		return nil, err
	}
	return f, nil
}
