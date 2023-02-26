package database

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	CollectionLogs = "logs"
)

func (mc *MongoClient) AddLogs(level string, message string, details map[string]interface{}) error {
	doc := DBLogs{
		Level:     level,
		Message:   message,
		Details:   details,
		Timestamp: time.Now().UnixMilli(),
	}
	if err := mc.InsertDocument(CollectionLogs, doc); err != nil {
		return err
	}
	return nil
}

func (mc *MongoClient) GetLogs(n int, skip int) ([]DBLogs, error) {
	var res []DBLogs
	opts := options.Find().SetSort(bson.D{{Key: "timestamp", Value: -1}}).SetSkip(int64(skip)).SetLimit(int64(n))
	err := mc.FindDocuments(CollectionLogs, emptyFilter, &res, opts)
	return res, err
}
