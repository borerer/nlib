package ssh

import (
	nlibshared "github.com/borerer/nlib-shared/go"
	"github.com/borerer/nlib/app/common"
)

type SSHApp struct {
}

func NewSSHApp() *SSHApp {
	return &SSHApp{}
}

func (a *SSHApp) AppID() string {
	return "ssh"
}

func (a *SSHApp) CallFunction(name string, req *nlibshared.Request) *nlibshared.Response {
	switch name {
	case "exec":
		return a.exec(req)
	}
	return common.Err404
}

func (a *SSHApp) exec(req *nlibshared.Request) *nlibshared.Response {
	// var sshConfig SSHConfig
	return nil
}
