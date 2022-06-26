package database

import (
	"context"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	database = "nlib"
	timeout  = time.Second * 15
)

type MongoClient struct {
	mc     *mongo.Client
	db     *mongo.Database
	host   string
	lock   sync.Mutex
	ctx    context.Context
	cancel context.CancelFunc
}

func NewMongoClient(host string) *MongoClient {
	return &MongoClient{
		host: host,
	}
}

func (c *MongoClient) maybeInit() error {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.mc != nil {
		return nil
	}
	c.ctx, c.cancel = context.WithTimeout(context.Background(), timeout)
	var err error
	c.mc, err = mongo.Connect(c.ctx, options.Client().ApplyURI(c.host))
	if err != nil {
		return err
	}
	c.db = c.mc.Database(database)
	return nil
}

func (c *MongoClient) Insert(colName string) error {
	if err := c.maybeInit(); err != nil {
		return err
	}
	c.lock.Lock()
	defer c.lock.Unlock()
	col := c.db.Collection(colName)
	return nil
}

func (c *MongoClient) Stop() error {
	if c.cancel != nil {
		c.cancel()
	}
	return nil
}
