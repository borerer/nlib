package app

import (
	"sync"

	nlibshared "github.com/borerer/nlib-shared/go"
	"github.com/borerer/nlib/app/common"
	"github.com/borerer/nlib/app/remote"
	"github.com/borerer/nlib/configs"
	"github.com/gorilla/websocket"
)

type AppManager struct {
	// remote
	remoteApps sync.Map

	appManagerBuiltin *AppManagerBuiltin
}

func NewAppManager(config *configs.BuiltinConfig) *AppManager {
	m := &AppManager{
		appManagerBuiltin: NewAppManagerBuiltin(config),
	}
	return m
}

func (m *AppManager) Start() error {
	if err := m.appManagerBuiltin.Start(); err != nil {
		return err
	}
	return nil
}

func (m *AppManager) Stop() error {
	if err := m.appManagerBuiltin.Stop(); err != nil {
		return err
	}
	return nil
}

// the unified interface to call functions from both builtin and remote apps
func (m *AppManager) CallFunction(appID string, name string, req *nlibshared.Request) *nlibshared.Response {
	res, ok := m.appManagerBuiltin.CallFunction(appID, name, req)
	if ok {
		return res
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
