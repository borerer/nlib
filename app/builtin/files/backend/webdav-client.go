package backend

import (
	"io"

	"github.com/borerer/nlib/configs"
	"github.com/studio-b12/gowebdav"
)

type WebdavClient struct {
	config *configs.WebdavConfig
	client *gowebdav.Client
}

func NewWebdavClient(config *configs.WebdavConfig) *WebdavClient {
	return &WebdavClient{
		config: config,
	}
}

func (c *WebdavClient) initClient() error {
	if c.client != nil {
		return nil
	}
	client := gowebdav.NewClient(c.config.Endpoint, c.config.User, c.config.Password)
	c.client = client
	return nil
}

func (c *WebdavClient) Start() error {
	if err := c.initClient(); err != nil {
		return err
	}
	return nil
}

func (c *WebdavClient) Stop() error {
	return nil
}

func (c *WebdavClient) GetFile(filename string) (io.ReadCloser, error) {
	reader, err := c.client.ReadStream(filename)
	if err != nil {
		return nil, err
	}
	return reader, nil
}

func (c *WebdavClient) PutFile(filename string, contentType string, override bool, fileReader io.Reader) (int64, error) {
	err := c.client.WriteStream(filename, fileReader, 0644)
	if err != nil {
		return 0, err
	}
	info, err := c.HeadFile(filename)
	if err != nil {
		return 0, err
	}
	return info.Size, nil
}

func (c *WebdavClient) DeleteFile(filename string) error {
	err := c.client.Remove(filename)
	if err != nil {
		return err
	}
	return nil
}

func (c *WebdavClient) HeadFile(filename string) (*FileInfo, error) {
	info, err := c.client.Stat(filename)
	if err != nil {
		return nil, err
	}
	return &FileInfo{
		Size:         info.Size(),
		LastModified: info.ModTime().UnixMilli(),
	}, nil
}

// func (c *WebdavClient) ListFolder(folder string) ([]string, error) {
// 	objectCh := mc.client.ListObjects(context.Background(), mc.config.Bucket, minio.ListObjectsOptions{
// 		Prefix:    folder,
// 		Recursive: false,
// 	})
// 	var res []string
// 	for obj := range objectCh {
// 		res = append(res, obj.Key)
// 	}
// 	return res, nil
// }
