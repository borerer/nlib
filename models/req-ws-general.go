package models

type ReqWSGeneral struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

const (
	WebSocketTypeDefault = "default"
	WebSocketTypeStart   = "start"
)
