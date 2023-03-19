package database

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (mc *MongoClient) SetKey(key string, value string) error {
	var doc DBKV
	err := mc.FindOneDocument(FilterEquals("key", key), &doc)
	if err != nil {
		// for both ErrNoDocuments and other errors, try create a new one
		doc = DBKV{
			Key:     key,
			Value:   value,
			Created: time.Now().UnixMilli(),
			Updated: time.Now().UnixMilli(),
		}
	} else {
		doc.Value = value
		doc.Updated = time.Now().UnixMilli()
	}
	if err := mc.UpdateDocument(FilterEquals("key", key), doc); err != nil {
		return err
	}
	return nil
}

// returned error:
//  1. ErrNoDocuments
//  2. others
func (mc *MongoClient) GetKey(key string) (string, error) {
	var doc DBKV
	err := mc.FindOneDocument(FilterEquals("key", key), &doc)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return "", ErrNoDocuments
	} else if err != nil {
		return "", err
	}
	return doc.Value, nil
}

func (mc *MongoClient) GetRecent(skip int, limit int) ([]DBKV, error) {
	var res []DBKV
	err := mc.FindDocuments(NoFilter, &res, options.Find().SetSkip(int64(skip)).SetLimit(int64(limit)).SetSort(bson.D{{Key: "updated", Value: -1}}))
	if err != nil {
		return nil, err
	}
	return res, nil
}
