package ssh

import (
	"encoding/json"

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
	}
	return common.Err404
}

func (a *SSHApp) Exec(sshConfig string, command string) (string, error) {
	str, err := a.kvApp.GetKey(sshConfig)
	if err != nil {
		return "", err
	}
	var config SSHConfig
	err = json.Unmarshal([]byte(str), &config)
	if err != nil {
		return "", err
	}
	client, err := goph.NewUnknown(config.User, config.Host, goph.Password(config.Password))
	if err != nil {
		return "", err
	}
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
		a.logsApp.Error("error executing command", "ssh-config", sshConfig, "command", command)
		return common.Error(err)
	}
	a.logsApp.Info("successfully executed", "ssh-config", sshConfig, "command", command, "output", output)
	return common.JSON(map[string]interface{}{
		"output": output,
	})
}
