package database

import (
	"context"

	"github.com/borerer/nlib/configs"
	"github.com/borerer/nlib/logs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type DatabaseManager struct {
	config *configs.DatabaseConfig
	client *mongo.Client
}

func NewDatabaseManager(config *configs.DatabaseConfig) *DatabaseManager {
	return &DatabaseManager{
		config: config,
	}
}

func (dm *DatabaseManager) createMongoClient() error {
	if dm.client != nil {
		return nil
	}
	var err error
	logs.Info("create mongo client", zap.String("host", dm.config.Mongo))
	dm.client, err = mongo.Connect(context.Background(), options.Client().ApplyURI(dm.config.Mongo))
	if err != nil {
		return err
	}
	return nil
}

func (dm *DatabaseManager) Start() error {
	logs.Info("database manager start")
	if err := dm.createMongoClient(); err != nil {
		return err
	}
	return nil
}

func (dm *DatabaseManager) Stop() error {
	return nil
}

func (dm *DatabaseManager) InsertDocument(colName string, doc interface{}) error {
	col := dm.client.Database(dm.config.Database).Collection(colName)
	_, err := col.InsertOne(context.Background(), doc)
	if err != nil {
		return err
	}
	return nil
}

func (dm *DatabaseManager) UpdateDocument(colName string, filter interface{}, doc interface{}) error {
	col := dm.client.Database(dm.config.Database).Collection(colName)
	update := bson.D{{Key: "$set", Value: doc}}
	opt := options.Update().SetUpsert(true)
	_, err := col.UpdateOne(context.Background(), filter, update, opt)
	if err != nil {
		return err
	}
	return nil
}

func (dm *DatabaseManager) FindDocuments(colName string, filter interface{}, res interface{}) error {
	col := dm.client.Database(dm.config.Database).Collection(colName)
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
