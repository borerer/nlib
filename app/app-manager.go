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
	config   *configs.BuiltinConfig
	echoApp  *echo.EchoApp
	kvApp    *kv.KVApp
	logsApp  *logs.LogsApp
	filesApp *files.FilesApp
}

func NewAppManager(config *configs.BuiltinConfig) *AppManager {
	m := &AppManager{
		config: config,
	}
	return m
}

func (m *AppManager) Start() error {
	if m.config.Echo.Enabled {
		m.echoApp = echo.NewEchoApp()
	}
	if m.config.KV.Enabled {
		m.kvApp = kv.NewKVApp(&m.config.KV)
		if err := m.kvApp.Start(); err != nil {
			return err
		}
	}
	if m.config.Logs.Enabled {
		m.logsApp = logs.NewLogsApp(&m.config.Logs)
		if err := m.logsApp.Start(); err != nil {
			return err
		}
	}
	if m.config.Files.Enabled {
		m.filesApp = files.NewFilesApp(&m.config.Files)
		if err := m.filesApp.Start(); err != nil {
			return err
		}
	}
	return nil
}

func (m *AppManager) Stop() error {
	if m.kvApp != nil {
		if err := m.kvApp.Stop(); err != nil {
			return err
		}
	}
	if m.logsApp != nil {
		if err := m.logsApp.Stop(); err != nil {
			return err
		}
	}
	if m.filesApp != nil {
		if err := m.filesApp.Stop(); err != nil {
			return err
		}
	}
	return nil
}

// the unified interface to call functions from both builtin and remote apps
func (m *AppManager) CallFunction(appID string, name string, req *nlibshared.Request) *nlibshared.Response {
	switch appID {
	case m.echoApp.AppID():
		return m.echoApp.CallFunction(name, req)
	case m.kvApp.AppID():
		return m.kvApp.CallFunction(name, req)
	case m.logsApp.AppID():
		return m.logsApp.CallFunction(name, req)
	case m.filesApp.AppID():
		return m.filesApp.CallFunction(name, req)
	}
	remoteApp, ok := m.getRemoteApp(appID)
	if !ok {
		return common.Err404
	}
	return remoteApp.CallFunction(name, req)
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
