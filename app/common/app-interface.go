package common

import nlibshared "github.com/borerer/nlib-shared/go"

type AppInterface interface {
	AppID() string
	CallFunction(name string, req *nlibshared.Request) *nlibshared.Response
}
