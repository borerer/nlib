package engines

import "io"

type FileEngine interface {
	GetFile(file string) (io.Reader, error)
}
