package models

type WebSocketCallFunctionRes struct {
	FuncName string      `json:"func_name" mapstructure:"func_name"`
	Response interface{} `json:"response" mapstructure:"response"`
}
