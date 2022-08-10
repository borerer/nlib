package models

type WebSocketMessage struct {
	MessageID     string      `json:"message_id"`
	PairMessageID string      `json:"pair_message_id"`
	Type          string      `json:"type"`
	SubType       string      `json:"sub_type"`
	Timestamp     int64       `json:"timestamp"`
	Payload       interface{} `json:"payload"`
}

const (
	WebSocketTypeDefault  = "default"
	WebSocketTypeRequest  = "request"
	WebSocketTypeResponse = "response"
)

const (
	WebSocketSubTypeDefault          = "default"
	WebSocketSubTypeStart            = "start"
	WebSocketSubTypeRegisterFunction = "register_function"
	WebSocketSubTypeCallFunction     = "call_function"
)
