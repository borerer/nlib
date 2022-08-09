package models

type DBLogs struct {
	AppID             string      `json:"app_id" bson:"app_id"`
	Message           string      `json:"message" bson:"message"`
	StructuredMessage interface{} `json:"structured_message" bson:"structured_message"`
}
