package models

type WebSocketCallFunctionReq struct {
	FuncName string `json:"func_name" mapstructure:"func_name"`
	Params   string `json:"params" mapstructure:"params"`
}
