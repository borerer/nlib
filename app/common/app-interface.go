package common

import nlibshared "github.com/borerer/nlib-shared/go"

type AppInterface interface {
	Start() error
	Stop() error
	AppID() string
	CallFunction(name string, req *nlibshared.Request) *nlibshared.Response
}
