package file

import (
	"io"

	"github.com/borerer/nlib/configs"
)

type FileHelper interface {
	GetFile(filename string) (io.ReadCloser, error)
	PutFile(filename string, override bool, fileReader io.Reader) (int64, error)
}

func NewFileHelper(config *configs.FileConfig) FileHelper {
	switch config.Engine {
	case "fs":
		return NewFileFS(config.FS.Dir)
	case "minio":
		return NewFileMinio(&FileMinioConfig{
			Endpoint:  config.Minio.Endpoint,
			AccessKey: config.Minio.AccessKey,
			SecretKey: config.Minio.SecretKey,
			UseSSL:    config.Minio.UseSSL,
			Bucket:    config.Minio.Bucket,
		})
	default:
		return nil
	}
}
