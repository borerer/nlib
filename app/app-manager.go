package app

import (
	nlibshared "github.com/borerer/nlib-shared/go"
	"github.com/borerer/nlib/app/builtin/kv"
	"github.com/borerer/nlib/app/common"
	"github.com/borerer/nlib/app/remote"
	"github.com/borerer/nlib/configs"
	"github.com/gorilla/websocket"
)

type AppManager struct {
	config *configs.BuiltinConfig
	kvApp  *kv.KVApp
}

func NewAppManager(config *configs.BuiltinConfig) *AppManager {
	m := &AppManager{
		config: config,
	}
	return m
}

func (m *AppManager) Start() error {
	m.kvApp = kv.NewKVApp(&m.config.KV)
	if err := m.kvApp.Start(); err != nil {
		return err
	}
	return nil
}

func (m *AppManager) Stop() error {
	return nil
}

// the unified interface to call functions from both builtin and remote apps
func (m *AppManager) CallFunction(appID string, name string, req *nlibshared.Request) *nlibshared.Response {
	switch appID {
	case m.kvApp.AppID():
		return m.kvApp.CallFunction(name, req)
	}
	// builtinAppRaw, ok := m.builtinApps.Load(appID)
	// if ok {
	// 	builtinApp, ok := builtinAppRaw.(*builtin.BuiltInApp)
	// 	if !ok {
	// 		return common.Err500
	// 	}
	// 	return builtinApp.CallFunction(name, req)
	// }
	// remoteAppRaw, ok := m.builtinApps.Load(appID)
	// if ok {
	// 	remoteApp, ok := remoteAppRaw.(*remote.RemoteApp)
	// 	if !ok {
	// 		return common.Err500
	// 	}
	// 	return remoteApp.CallFunction(name, req)
	// }
	return common.Err404
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
