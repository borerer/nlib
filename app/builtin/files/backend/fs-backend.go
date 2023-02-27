package backend

import "io"

type FSBackend interface {
	Start() error
	Stop() error
	HeadFile(filename string) (*FileInfo, error)
	GetFile(filename string) (io.ReadCloser, error)
	PutFile(filename string, contentType string, override bool, fileReader io.Reader) (int64, error)
	DeleteFile(filename string) error
}
