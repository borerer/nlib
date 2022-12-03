package database

import (
	"time"
)

const (
	CollectionKV = "kv"
)

func (mc *MongoClient) SetKey(appID string, key string, value string) error {
	var doc DBKV
	err := mc.FindOneDocument(CollectionKV, FilterEquals("key", key), &doc)
	if err != nil {
		// for both ErrNoDocuments and other errors, try create a new one
		doc = DBKV{
			AppID:   appID,
			Key:     key,
			Value:   value,
			Created: time.Now().UnixMilli(),
			Updated: time.Now().UnixMilli(),
		}
	} else {
		doc.Value = value
		doc.Updated = time.Now().UnixMilli()
	}
	if err := mc.UpdateDocument(CollectionKV, FilterEquals("key", key), doc); err != nil {
		return err
	}
	return nil
}

// returned error:
//  1. mongo.ErrNoDocuments
//  2. others
func (mc *MongoClient) GetKey(appID string, key string) (string, error) {
	var doc DBKV
	if err := mc.FindOneDocument(CollectionKV, FilterEquals("key", key), &doc); err != nil {
		return "", err
	}
	return doc.Value, nil
}
