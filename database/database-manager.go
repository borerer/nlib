package database

import (
	"fmt"

	"github.com/borerer/nlib/configs"
	"github.com/borerer/nlib/logs"
)

type DatabaseManager struct {
	config *configs.DatabaseConfig
	client *MongoClient
}

func NewDatabaseManager(config *configs.DatabaseConfig) *DatabaseManager {
	return &DatabaseManager{
		config: config,
	}
}

func (dm *DatabaseManager) Start() error {
	logs.Info("database manager start")
	dm.client = NewMongoClient(dm.config.Mongo, dm.config.Database)
	if err := dm.client.Init(); err != nil {
		return err
	}
	return nil
}

func (dm *DatabaseManager) Stop() error {
	return nil
}

func (dm *DatabaseManager) GetKey(appID string, key string) (string, error) {
	colName := fmt.Sprintf("%s_kv", appID)
	return dm.client.Get(colName, key)
}
