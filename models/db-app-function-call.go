package models

type DBAppFunctionCall struct {
	FuncID string `json:"func_id" bson:"func_id"`
	// CallInfo
	Request        interface{} `json:"request" bson:"request"`
	Response       interface{} `json:"response" bson:"response"`
	StartTimestamp int64       `json:"start_timestamp" bson:"start_timestamp"`
	EndTimestamp   int64       `json:"end_timestamp" bson:"end_timestamp"`
	Status         string      `json:"status" bson:"status"`
}
