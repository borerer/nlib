package models

type WebSocketCallFunctionRes struct {
	FuncName string                 `json:"func_name" mapstructure:"func_name"`
	Response map[string]interface{} `json:"response" mapstructure:"response"`
}
