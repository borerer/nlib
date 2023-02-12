package app

import (
	"sync"

	nlibshared "github.com/borerer/nlib-shared/go"
	"github.com/borerer/nlib/app/builtin"
	"github.com/borerer/nlib/app/builtin/kv"
	"github.com/borerer/nlib/app/common"
	"github.com/borerer/nlib/app/remote"
	"github.com/borerer/nlib/configs"
	"github.com/gorilla/websocket"
)

type AppManager struct {
	config      *configs.ServerConfig
	builtinApps sync.Map
}

func NewAppManager(config *configs.ServerConfig) *AppManager {
	m := &AppManager{
		config: config,
	}
	return m
}

func (m *AppManager) CallFunction(appID string, name string, req *nlibshared.Request) *nlibshared.Response {
	builtinAppRaw, ok := m.builtinApps.Load(appID)
	if ok {
		builtinApp, ok := builtinAppRaw.(*builtin.BuiltInApp)
		if !ok {
			return common.Err500
		}
		return builtinApp.CallFunction(name, req)
	}
	remoteAppRaw, ok := m.builtinApps.Load(appID)
	if ok {
		remoteApp, ok := remoteAppRaw.(*remote.RemoteApp)
		if !ok {
			return common.Err500
		}
		return remoteApp.CallFunction(name, req)
	}
	return common.Err404
}

func (m *AppManager) installBuiltinApps() {
	kvApp := kv.NewKVApp()
}

func (m *AppManager) AddBuiltinApp(appID string) {
	app := builtin.NewBuiltinApp(appID)
	m.builtinApps.Store(appID, app)
}

func (m *AppManager) AddRemoteApp(appID string, conn *websocket.Conn) error {
	app := remote.NewRemoteApp(appID)
	app.SetWebSocketConnection(conn)
	return app.ListenWebSocketMessages()
}

// func (m *AppManager) getRemoteApp(appID string) *App {
// 	var app *App
// 	raw, ok := m.apps.Load(appID)
// 	if ok {
// 		app, ok = raw.(*App)
// 		if !ok {
// 			logs.Warn("unexpected get nlib app error", zap.String("appID", appID))
// 			// fallback to create a new client instance
// 		}
// 	}
// 	if app == nil {

// 	}
// 	return app
// }

// func (m *AppManager) getBuiltinApp(appID string) *builtin.BuiltInApp {
// 	var app *builtin.BuiltInApp
// 	raw, ok := m.apps.Load(appID)
// 	if ok {
// 		app, ok = raw.(*App)
// 		if !ok {
// 			logs.Warn("unexpected get nlib app error", zap.String("appID", appID))
// 			// fallback to create a new client instance
// 		}
// 	}
// 	if app == nil {
// 		app = NewApp(appID)
// 		m.apps.Store(appID, app)
// 	}
// 	return app
// }
