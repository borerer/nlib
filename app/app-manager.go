package app

import (
	"sync"

	nlibshared "github.com/borerer/nlib-shared/go"
	"github.com/borerer/nlib/logs"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type AppManager struct {
	apps sync.Map
}

func NewAppManager() *AppManager {
	m := &AppManager{}
	return m
}

func (m *AppManager) StartConnection(appID string, conn *websocket.Conn) error {
	client := m.getApp(appID)
	client.SetWebSocketConnection(conn)
	return client.ListenWebSocketMessages()
}

func (m *AppManager) FunctionUseHAR(appID string, name string) bool {
	app := m.getApp(appID)
	return app.FunctionUseHAR(name)
}

func (m *AppManager) CallSimpleFunction(appID string, name string, req *nlibshared.SimpleFunctionIn) (nlibshared.SimpleFunctionOut, error) {
	app := m.getApp(appID)
	return app.CallSimpleFunction(name, req)
}

func (m *AppManager) CallHARFunction(appID string, name string, req *nlibshared.HARFunctionIn) (*nlibshared.HARFunctionOut, error) {
	app := m.getApp(appID)
	return app.CallHARFunction(name, req)
}

func (m *AppManager) getApp(appID string) *App {
	var app *App
	raw, ok := m.apps.Load(appID)
	if ok {
		app, ok = raw.(*App)
		if !ok {
			logs.Warn("unexpected get nlib app error", zap.String("appID", appID))
			// fallback to create a new client instance
		}
	}
	if app == nil {
		app = NewApp(appID)
		m.apps.Store(appID, app)
	}
	return app
}
