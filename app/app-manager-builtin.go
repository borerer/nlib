package app

import (
	nlibshared "github.com/borerer/nlib-shared/go"
	"github.com/borerer/nlib/app/builtin/echo"
	"github.com/borerer/nlib/app/builtin/files"
	"github.com/borerer/nlib/app/builtin/kv"
	"github.com/borerer/nlib/app/builtin/logs"
	"github.com/borerer/nlib/app/builtin/ssh"
	"github.com/borerer/nlib/app/common"
	"github.com/borerer/nlib/configs"
)

type AppManagerBuiltin struct {
	config *configs.BuiltinConfig
	apps   map[string]common.AppInterface
}

func NewAppManagerBuiltin(config *configs.BuiltinConfig) *AppManagerBuiltin {
	m := &AppManagerBuiltin{
		config: config,
	}
	return m
}

func (m *AppManagerBuiltin) Start() error {
	echoApp := echo.NewEchoApp()
	kvApp := kv.NewKVApp(m.config.Mongo)
	logsApp := logs.NewLogsApp(m.config.Mongo)
	filesApp := files.NewFilesApp(m.config)
	sshApp := ssh.NewSSHApp(kvApp, logsApp)
	m.apps = map[string]common.AppInterface{
		echoApp.AppID():  echoApp,
		kvApp.AppID():    kvApp,
		logsApp.AppID():  logsApp,
		filesApp.AppID(): filesApp,
		sshApp.AppID():   sshApp,
	}
	for _, app := range m.apps {
		if err := app.Start(); err != nil {
			return err
		}
	}
	return nil
}

func (m *AppManagerBuiltin) Stop() error {
	for _, app := range m.apps {
		if err := app.Stop(); err != nil {
			return err
		}
	}
	return nil
}

func (m *AppManagerBuiltin) CallFunction(appID string, name string, req *nlibshared.Request) (*nlibshared.Response, bool) {
	app, ok := m.apps[appID]
	if ok {
		return app.CallFunction(name, req), true
	}
	return nil, false
}
