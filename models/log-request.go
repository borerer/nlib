package models

type LogRequest struct {
	ID      string      `json:"id"`
	Level   string      `json:"level"`
	Message string      `json:"message"`
	Body    interface{} `json:"body"`
}
