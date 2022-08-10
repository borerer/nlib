package models

type WebSocketCallFunctionReq struct {
	FuncName string `json:"func_name"`
	Params   string `json:"params"`
}
