package file

import (
	"context"
	"io"

	"github.com/borerer/nlib/logs"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.uber.org/zap"
)

type FileMinio struct {
	Config *FileMinioConfig
	client *minio.Client
}

type FileMinioConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	UseSSL    bool
	Bucket    string
}

func NewFileMinio(config *FileMinioConfig) *FileMinio {
	return &FileMinio{
		Config: config,
	}
}

func (f *FileMinio) initClient() error {
	if f.client != nil {
		return nil
	}
	client, err := minio.New(f.Config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(f.Config.AccessKey, f.Config.SecretKey, ""),
		Secure: f.Config.UseSSL,
	})
	if err != nil {
		return err
	}
	f.client = client
	return nil
}

func (f *FileMinio) Start() error {
	logs.Info("file minio start")
	if err := f.initClient(); err != nil {
		return err
	}
	return nil
}

func (f *FileMinio) Stop() error {
	return nil
}

func (f *FileMinio) GetFile(filename string) (io.ReadCloser, error) {
	if err := f.initClient(); err != nil {
		return nil, err
	}
	obj, err := f.client.GetObject(context.Background(), f.Config.Bucket, filename, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func (f *FileMinio) PutFile(filename string, override bool, fileReader io.Reader) (int64, error) {
	if err := f.initClient(); err != nil {
		return 0, err
	}
	info, err := f.client.PutObject(context.Background(), f.Config.Bucket, filename, fileReader, -1, minio.PutObjectOptions{})
	if err != nil {
		return 0, err
	}
	return info.Size, nil
}

func (f *FileMinio) DeleteFile(filename string) error {
	if err := f.initClient(); err != nil {
		return err
	}
	err := f.client.RemoveObject(context.Background(), f.Config.Bucket, filename, minio.RemoveObjectOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (f *FileMinio) HeadFile(filename string) (*FileInfo, error) {
	if err := f.initClient(); err != nil {
		return nil, err
	}
	info, err := f.client.StatObject(context.Background(), f.Config.Bucket, filename, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	return &FileInfo{
		Size:         info.Size,
		LastModified: info.LastModified.UnixMilli(),
		ContentType:  info.ContentType,
	}, nil
}

func (f *FileMinio) ListFolder(folder string) ([]string, error) {
	if err := f.initClient(); err != nil {
		return nil, err
	}
	objectCh := f.client.ListObjects(context.Background(), f.Config.Bucket, minio.ListObjectsOptions{
		Prefix:    folder,
		Recursive: false,
	})
	var res []string
	logs.Info("yyyy")
	for obj := range objectCh {
		logs.Info("xxxx", zap.Any("xxx", obj))
		res = append(res, obj.Key)
	}
	return res, nil
}
