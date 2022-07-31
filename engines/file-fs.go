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

func (fileFS *FileFS) PutFile(filename string, override bool, fileReader io.Reader) (int64, error) {
	absFilePath := path.Join(fileFS.BaseDir, filename)
	f, err := os.Create(absFilePath)
	if err != nil {
		return 0, err
	}
	n, err := io.Copy(f, fileReader)
	return n, err
}
