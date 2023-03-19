package app

import (
	"sync"

	nlibshared "github.com/borerer/nlib-shared/go"
	"github.com/borerer/nlib/app/builtin/echo"
	"github.com/borerer/nlib/app/builtin/files"
	"github.com/borerer/nlib/app/builtin/kv"
	"github.com/borerer/nlib/app/builtin/logs"
	"github.com/borerer/nlib/app/common"
	"github.com/borerer/nlib/app/remote"
	"github.com/borerer/nlib/configs"
	"github.com/gorilla/websocket"
)

type AppManager struct {
	// remote
	remoteApps sync.Map

	// builtin
	config      *configs.BuiltinConfig
	builtinApps map[string]common.AppInterface
}

func NewAppManager(config *configs.BuiltinConfig) *AppManager {
	m := &AppManager{
		config: config,
	}
	return m
}

func (m *AppManager) Start() error {
	echoApp := echo.NewEchoApp()
	kvApp := kv.NewKVApp(m.config.Mongo)
	logsApp := logs.NewLogsApp(m.config.Mongo)
	filesApp := files.NewFilesApp(m.config)
	m.builtinApps = map[string]common.AppInterface{
		echoApp.AppID():  echoApp,
		kvApp.AppID():    kvApp,
		logsApp.AppID():  logsApp,
		filesApp.AppID(): filesApp,
	}
	for _, app := range m.builtinApps {
		if err := app.Start(); err != nil {
			return err
		}
	}
	return nil
}

func (m *AppManager) Stop() error {
	for _, app := range m.builtinApps {
		if err := app.Stop(); err != nil {
			return err
		}
	}
	return nil
}

// the unified interface to call functions from both builtin and remote apps
func (m *AppManager) CallFunction(appID string, name string, req *nlibshared.Request) *nlibshared.Response {
	builtinApp, ok := m.builtinApps[appID]
	if ok {
		return builtinApp.CallFunction(name, req)
	}
	remoteApp, ok := m.getRemoteApp(appID)
	if ok {
		return remoteApp.CallFunction(name, req)
	}
	return common.Err404
}

func (m *AppManager) AddRemoteApp(appID string, conn *websocket.Conn) error {
	remoteApp := remote.NewRemoteApp(appID)
	remoteApp.SetWebSocketConnection(conn)
	m.remoteApps.Store(appID, remoteApp)
	return remoteApp.ListenWebSocketMessages()
}

func (m *AppManager) getRemoteApp(appID string) (*remote.RemoteApp, bool) {
	raw, ok := m.remoteApps.Load(appID)
	if ok {
		app, ok := raw.(*remote.RemoteApp)
		if ok {
			return app, true
		}
	}
	return nil, false
}
