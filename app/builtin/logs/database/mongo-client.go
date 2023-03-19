package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoClient struct {
	mongoURI string
	client   *mongo.Client
}

const (
	MongoDatabase   = "nlib"
	MongoCollection = "logs"
)

func NewMongoClient(mongoURI string) *MongoClient {
	return &MongoClient{
		mongoURI: mongoURI,
	}
}

func (mc *MongoClient) connect() error {
	if mc.client != nil {
		return nil
	}
	var err error
	mc.client, err = mongo.Connect(context.Background(), options.Client().ApplyURI(mc.mongoURI))
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

func (mc *MongoClient) InsertDocument(doc interface{}) error {
	col := mc.client.Database(MongoDatabase).Collection(MongoCollection)
	_, err := col.InsertOne(context.Background(), doc)
	if err != nil {
		return err
	}
	return nil
}

func (mc *MongoClient) UpdateDocument(filter interface{}, doc interface{}) error {
	col := mc.client.Database(MongoDatabase).Collection(MongoCollection)
	update := bson.D{{Key: "$set", Value: doc}}
	opt := options.Update().SetUpsert(true)
	_, err := col.UpdateOne(context.Background(), filter, update, opt)
	if err != nil {
		return err
	}
	return nil
}

func (mc *MongoClient) FindDocuments(filter interface{}, res interface{}, opts ...*options.FindOptions) error {
	col := mc.client.Database(MongoDatabase).Collection(MongoCollection)
	cur, err := col.Find(context.Background(), filter, opts...)
	if err != nil {
		return err
	}
	err = cur.All(context.Background(), res)
	if err != nil {
		return err
	}
	return nil
}

func (mc *MongoClient) FindOneDocument(filter interface{}, res interface{}) error {
	col := mc.client.Database(MongoDatabase).Collection(MongoCollection)
	ret := col.FindOne(context.Background(), filter)
	if ret.Err() != nil {
		return ret.Err()
	}
	if err := ret.Decode(res); err != nil {
		return err
	}
	return nil
}
