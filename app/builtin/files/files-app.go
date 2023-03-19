package files

import (
	"bytes"
	"encoding/base64"
	"errors"
	"io"

	nlibshared "github.com/borerer/nlib-shared/go"
	"github.com/borerer/nlib/app/builtin/files/backend"
	"github.com/borerer/nlib/app/common"
	"github.com/borerer/nlib/configs"
)

var (
	EncodingBase64 = "base64"
)

type FilesApp struct {
	config        *configs.BuiltinConfig
	activeBackend backend.FSBackend
}

func NewFilesApp(config *configs.BuiltinConfig) *FilesApp {
	return &FilesApp{
		config: config,
	}
}

func (app *FilesApp) Start() error {
	if len(app.config.Webdav.Endpoint) > 0 {
		app.activeBackend = backend.NewWebdavClient(&app.config.Webdav)
	} else if len(app.config.Samba.Endpoint) > 0 {
		app.activeBackend = backend.NewSambaClient(&app.config.Samba)
	} else if len(app.config.Minio.Endpoint) > 0 {
		app.activeBackend = backend.NewMinioClient(&app.config.Minio)
	} else {
		return errors.New("no valid files backend")
	}
	if err := app.activeBackend.Start(); err != nil {
		return err
	}
	return nil
}

func (app *FilesApp) Stop() error {
	if err := app.activeBackend.Stop(); err != nil {
		return err
	}
	return nil
}

func (app *FilesApp) AppID() string {
	return "files"
}

func (app *FilesApp) CallFunction(name string, req *nlibshared.Request) *nlibshared.Response {
	switch name {
	case "get":
		return app.get(req)
	case "put":
		return app.put(req)
	default:
		return common.Err404
	}
}

func toBase64(buf []byte) string {
	return base64.StdEncoding.EncodeToString(buf)
}

func fromBase64(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}

func (app *FilesApp) get(req *nlibshared.Request) *nlibshared.Response {
	filename := common.GetQuery(req, "file")
	info, err := app.activeBackend.HeadFile(filename)
	if err != nil {
		return common.Error(err)
	}
	reader, err := app.activeBackend.GetFile(filename)
	if err != nil {
		return common.Error(err)
	}
	defer reader.Close()
	buf, err := io.ReadAll(reader)
	if err != nil {
		return common.Error(err)
	}
	b64Str := toBase64(buf)
	res := common.Text(b64Str)
	res.Headers = append(res.Headers, nlibshared.Header{Name: "Content-Type", Value: info.ContentType})
	res.Content.Encoding = &EncodingBase64
	return res
}

func (app *FilesApp) put(req *nlibshared.Request) *nlibshared.Response {
	filename := common.GetQuery(req, "file")
	contentType := common.GetHeader(req, "Content-Type")
	if req.PostData == nil || req.PostData.Text == nil {
		return common.Err400
	}
	buf, err := fromBase64(*req.PostData.Text)
	if err != nil {
		return common.Error(err)
	}
	reader := bytes.NewReader(buf)
	_, err = app.activeBackend.PutFile(filename, contentType, true, reader)
	if err != nil {
		return common.Error(err)
	}
	return common.Text("ok")
}
