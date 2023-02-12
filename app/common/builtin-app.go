package common

import (
	"sync"

	nlibshared "github.com/borerer/nlib-shared/go"
)

type BuiltInApp struct {
	appID               string
	registeredFunctions sync.Map
}

func NewBuiltinApp(appID string) *BuiltInApp {
	return &BuiltInApp{
		appID: appID,
	}
}

func (app *BuiltInApp) RegisterFunction(name string, f func(*nlibshared.Request) *nlibshared.Response) {
	app.registeredFunctions.Store(name, f)
}

func (app *BuiltInApp) CallFunction(name string, req *nlibshared.Request) *nlibshared.Response {
	funcRaw, ok := app.registeredFunctions.Load(name)
	if !ok {
		return Err404
	}
	f, ok := funcRaw.(func(*nlibshared.Request) *nlibshared.Response)
	if !ok {
		return Err500
	}
	return f(req)
}
