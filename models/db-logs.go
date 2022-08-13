package models

type DBLogs struct {
	AppID     string      `json:"app_id" bson:"app_id"`
	Level     string      `json:"level" bson:"level"`
	Message   string      `json:"message" bson:"message"`
	Details   interface{} `json:"details" bson:"details"`
	Timestamp int64       `json:"timestamp" bson:"timestamp"`
}
