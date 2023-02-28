package backend

import (
	"fmt"
	"io"
	"net"
	"os"
	"path"

	"github.com/borerer/nlib/configs"
	"github.com/hirochachacha/go-smb2"
)

type SambaClient struct {
	config  *configs.SambaConfig
	conn    net.Conn
	session *smb2.Session
	share   *smb2.Share
}

func NewSambaClient(config *configs.SambaConfig) *SambaClient {
	return &SambaClient{
		config: config,
	}
}

func (c *SambaClient) initClient() error {
	if c.share != nil {
		return nil
	}
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:445", c.config.Endpoint))
	if err != nil {
		return err
	}
	dialer := &smb2.Dialer{
		Initiator: &smb2.NTLMInitiator{
			User:     c.config.User,
			Password: c.config.Password,
		},
	}
	session, err := dialer.Dial(conn)
	if err != nil {
		return err
	}
	share, err := session.Mount(c.config.Share)
	if err != nil {
		return err
	}
	c.conn = conn
	c.session = session
	c.share = share
	return nil
}

func (c *SambaClient) Start() error {
	if err := c.initClient(); err != nil {
		return err
	}
	return nil
}

func (c *SambaClient) Stop() error {
	if err := c.share.Umount(); err != nil {
		return err
	}
	if err := c.session.Logoff(); err != nil {
		return err
	}
	if err := c.conn.Close(); err != nil {
		return err
	}
	return nil
}

func (c *SambaClient) getFullPath(filename string) string {
	return path.Join(c.config.Path, filename)
}

func (c *SambaClient) GetFile(filename string) (io.ReadCloser, error) {
	filename = c.getFullPath(filename)
	reader, err := c.share.Open(filename)
	if os.IsNotExist(err) {
		return nil, ErrFileNotFound
	}
	if err != nil {
		return nil, err
	}
	return reader, nil
}

func (c *SambaClient) PutFile(filename string, contentType string, override bool, fileReader io.Reader) (int64, error) {
	filename = c.getFullPath(filename)
	writer, err := c.share.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if os.IsNotExist(err) {
		writer, err = c.share.Create(filename)
	}
	if err != nil {
		return 0, err
	}
	defer writer.Close()
	return io.Copy(writer, fileReader)
}

func (c *SambaClient) DeleteFile(filename string) error {
	filename = c.getFullPath(filename)
	err := c.share.Remove(filename)
	if err != nil {
		return err
	}
	return nil
}

func (c *SambaClient) HeadFile(filename string) (*FileInfo, error) {
	filename = c.getFullPath(filename)
	info, err := c.share.Stat(filename)
	if err != nil {
		return nil, err
	}
	return &FileInfo{
		Size:         info.Size(),
		LastModified: info.ModTime().UnixMilli(),
	}, nil
}

// func (c *SambaClient) ListFolder(folder string) ([]string, error) {
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
