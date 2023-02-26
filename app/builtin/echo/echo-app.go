package echo

import (
	nlibshared "github.com/borerer/nlib-shared/go"
	"github.com/borerer/nlib/app/common"
)

type EchoApp struct {
}

func NewEchoApp() *EchoApp {
	return &EchoApp{}
}

func (a *EchoApp) AppID() string {
	return "echo"
}

func (a *EchoApp) CallFunction(name string, req *nlibshared.Request) *nlibshared.Response {
	m := map[string]interface{}{
		"name": name,
		"req":  req,
	}
	return common.JSON(m)
}
