package models

type APIAddLogRequest struct {
	Level   string      `json:"level"`
	Message string      `json:"message"`
	Details interface{} `json:"details"`
}
