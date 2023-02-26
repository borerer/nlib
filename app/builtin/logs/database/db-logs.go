package database

type DBLogs struct {
	Level     string                 `json:"level" bson:"level"`
	Message   string                 `json:"message" bson:"message"`
	Details   map[string]interface{} `json:"details" bson:"details"`
	Timestamp int64                  `json:"timestamp" bson:"timestamp"`
}
