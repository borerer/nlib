package database

import (
	"context"
	"time"

	"github.com/borerer/nlib/logs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	timeout = time.Second * 15
)

type MongoClient struct {
	client   *mongo.Client
	db       *mongo.Database
	host     string
	database string
}

func NewMongoClient(host string, database string) *MongoClient {
	return &MongoClient{
		host:     host,
		database: database,
	}
}

func (c *MongoClient) Init() error {
	if c.client != nil {
		return nil
	}
	var err error
	logs.Info(c.host)
	c.client, err = mongo.Connect(context.Background(), options.Client().ApplyURI(c.host))
	if err != nil {
		return err
	}
	c.db = c.client.Database(c.database)
	return nil
}

func (c *MongoClient) Insert(colName string) error {
	// col := c.db.Collection(colName)
	return nil
}

func (c *MongoClient) Get(colName string, documentID string) (string, error) {
	col := c.db.Collection(colName)
	res := col.FindOne(context.Background(), bson.D{{"_id", documentID}})
	if res.Err() != nil {
		return "", res.Err()
	}
	if raw, err := res.DecodeBytes(); err != nil {
		return "", err
	} else {
		return raw.String(), nil
	}
}

func (c *MongoClient) Stop() error {
	return nil
}
