package engines

import "io"

type FileEngine interface {
	GetFile(filename string) (io.Reader, error)
	PutFile(filename string, override bool, fileReader io.Reader) (int64, error)
}
