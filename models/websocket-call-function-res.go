package models

type WebSocketCallFunctionRes struct {
	FuncName string `json:"func_name" mapstructure:"func_name"`
	Response string `json:"response" mapstructure:"response"`
}
