package ssh

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"path/filepath"

	nlibshared "github.com/borerer/nlib-shared/go"
	"github.com/borerer/nlib/app/builtin/kv"
	"github.com/borerer/nlib/app/builtin/logs"
	"github.com/borerer/nlib/app/common"
	"github.com/melbahja/goph"
)

type SSHApp struct {
	kvApp   *kv.KVApp
	logsApp *logs.LogsApp
}

var (
	EncodingBase64 = "base64"
)

func NewSSHApp(kvApp *kv.KVApp, logsApp *logs.LogsApp) *SSHApp {
	return &SSHApp{
		kvApp:   kvApp,
		logsApp: logsApp,
	}
}

func (a *SSHApp) AppID() string {
	return "ssh"
}

func (a *SSHApp) Start() error {
	return nil
}

func (a *SSHApp) Stop() error {
	return nil
}

func (a *SSHApp) CallFunction(name string, req *nlibshared.Request) *nlibshared.Response {
	switch name {
	case "exec":
		return a.exec(req)
	case "download":
		return a.download(req)
	}
	return common.Err404
}

func (a *SSHApp) getSSHClient(sshConfig string) (*goph.Client, error) {
	str, err := a.kvApp.GetKey(sshConfig)
	if err != nil {
		return nil, err
	}
	var config SSHConfig
	err = json.Unmarshal([]byte(str), &config)
	if err != nil {
		return nil, err
	}
	client, err := goph.NewUnknown(config.User, config.Host, goph.Password(config.Password))
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (a *SSHApp) Exec(sshConfig string, command string) (string, error) {
	client, err := a.getSSHClient(sshConfig)
	if err != nil {
		return "", err
	}
	defer client.Close()
	buf, err := client.Run(command)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

func (a *SSHApp) exec(req *nlibshared.Request) *nlibshared.Response {
	sshConfig := common.GetQuery(req, "ssh-config")
	command := common.GetQuery(req, "command")
	output, err := a.Exec(sshConfig, command)
	if err != nil {
		a.logsApp.Error("error executing command", "ssh-config", sshConfig, "command", command, "error", err.Error())
		return common.Error(err)
	}
	a.logsApp.Info("successfully executed", "ssh-config", sshConfig, "command", command, "output", output)
	return common.JSON(map[string]interface{}{
		"output": output,
	})
}

func toBase64(buf []byte) string {
	return base64.StdEncoding.EncodeToString(buf)
}

func fromBase64(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}

func (a *SSHApp) Download(sshConfig string, path string) ([]byte, error) {
	client, err := a.getSSHClient(sshConfig)
	if err != nil {
		return nil, err
	}
	defer client.Close()
	sftp, err := client.NewSftp()
	if err != nil {
		return nil, err
	}
	file, err := sftp.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	buf, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func (a *SSHApp) download(req *nlibshared.Request) *nlibshared.Response {
	sshConfig := common.GetQuery(req, "ssh-config")
	file := common.GetQuery(req, "file")
	buf, err := a.Download(sshConfig, file)
	if err != nil {
		a.logsApp.Error("error downloading", "ssh-config", sshConfig, "file", file, "error", err.Error())
		return common.Error(err)
	}
	contentType := http.DetectContentType(buf)
	if contentType == "application/octet-stream" {
		ext := filepath.Ext(file)
		mimeType := mime.TypeByExtension(ext)
		if len(mimeType) > 0 {
			contentType = mimeType
		}
	}
	filename := filepath.Base(file)
	contentDisposition := fmt.Sprintf("attachment; filename=\"%s\"", filename)
	b64Str := toBase64(buf)
	res := common.Text(b64Str)
	res.Headers = append(res.Headers, nlibshared.Header{Name: "Content-Type", Value: contentType})
	res.Headers = append(res.Headers, nlibshared.Header{Name: "Content-Disposition", Value: contentDisposition})
	res.Content.Encoding = &EncodingBase64
	return res
}
