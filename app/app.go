package app

import (
	"context"
	"net/http"
	"sync"

	"github.com/borerer/nlib/configs"
	"github.com/borerer/nlib/database"
	"github.com/borerer/nlib/file"
	"github.com/borerer/nlib/logs"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type App struct {
	config      *configs.AppConfig
	ginRouter   *gin.Engine
	httpServer  *http.Server
	minioClient *file.MinioClient
	mongoClient *database.MongoClient
	nlibClients sync.Map
}

func NewApp(config *configs.AppConfig) *App {
	app := &App{
		config:      config,
		minioClient: file.NewMinioClient(&config.Minio),
		mongoClient: database.NewMongoClient(&config.Mongo),
	}
	return app
}

func (app *App) Start() error {
	logs.Info("app start")
	if err := app.minioClient.Start(); err != nil {
		return err
	}
	if err := app.mongoClient.Start(); err != nil {
		return err
	}
	if err := app.startAPI(); err != nil {
		return err
	}
	return nil
}

func (app *App) Stop() error {
	logs.Info("app stop")
	if app.httpServer != nil {
		if err := app.httpServer.Shutdown(context.Background()); err != nil {
			return err
		}
	}
	if app.mongoClient != nil {
		if err := app.mongoClient.Stop(); err != nil {
			return err
		}
	}
	if app.minioClient != nil {
		if err := app.minioClient.Stop(); err != nil {
			return err
		}
	}
	return nil
}

func (app *App) GetNLIBClient(appID string) *NLIBClient {
	var client *NLIBClient
	clientRaw, ok := app.nlibClients.Load(appID)
	if ok {
		client, ok = clientRaw.(*NLIBClient)
		if !ok {
			logs.Warn("unexpected get nlib client error", zap.String("appID", appID))
			// fallback to create a new client instance
		}
	}
	if client == nil {
		client = NewNLIBClient(appID)
		app.nlibClients.Store(appID, client)
	}
	return client
}
