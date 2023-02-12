package database

import (
	"context"

	"github.com/borerer/nlib/configs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoClient struct {
	config *configs.MongoConfig
	client *mongo.Client
}

func NewMongoClient(config *configs.MongoConfig) *MongoClient {
	return &MongoClient{
		config: config,
	}
}

func (mc *MongoClient) connect() error {
	if mc.client != nil {
		return nil
	}
	var err error
	mc.client, err = mongo.Connect(context.Background(), options.Client().ApplyURI(mc.config.URI))
	if err != nil {
		return err
	}
	return nil
}

func (mc *MongoClient) Start() error {
	if err := mc.connect(); err != nil {
		return err
	}
	return nil
}

func (mc *MongoClient) Stop() error {
	return nil
}

func (mc *MongoClient) InsertDocument(colName string, doc interface{}) error {
	col := mc.client.Database(mc.config.Database).Collection(colName)
	_, err := col.InsertOne(context.Background(), doc)
	if err != nil {
		return err
	}
	return nil
}

func (mc *MongoClient) UpdateDocument(colName string, filter interface{}, doc interface{}) error {
	col := mc.client.Database(mc.config.Database).Collection(colName)
	update := bson.D{{Key: "$set", Value: doc}}
	opt := options.Update().SetUpsert(true)
	_, err := col.UpdateOne(context.Background(), filter, update, opt)
	if err != nil {
		return err
	}
	return nil
}

func (mc *MongoClient) FindDocuments(colName string, filter interface{}, res interface{}) error {
	col := mc.client.Database(mc.config.Database).Collection(colName)
	cur, err := col.Find(context.Background(), filter)
	if err != nil {
		return err
	}
	err = cur.All(context.Background(), res)
	if err != nil {
		return err
	}
	return nil
}

func (mc *MongoClient) FindOneDocument(colName string, filter interface{}, res interface{}) error {
	col := mc.client.Database(mc.config.Database).Collection(colName)
	ret := col.FindOne(context.Background(), filter)
	if ret.Err() != nil {
		return ret.Err()
	}
	if err := ret.Decode(res); err != nil {
		return err
	}
	return nil
}
