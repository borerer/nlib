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
	config          *configs.AppConfig
	ginRouter       *gin.Engine
	httpServer      *http.Server
	fileHelper      file.FileHelper
	databaseManager *database.DatabaseManager
	nlibClients     sync.Map
}

func NewApp(config *configs.AppConfig) *App {
	app := &App{
		config: config,
	}
	return app
}

func (app *App) Start() error {
	logs.Info("app start")
	app.fileHelper = file.NewFileHelper(&app.config.File)
	if err := app.fileHelper.Start(); err != nil {
		return err
	}
	app.databaseManager = database.NewDatabaseManager(&app.config.Database)
	if err := app.databaseManager.Start(); err != nil {
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
	if app.databaseManager != nil {
		if err := app.databaseManager.Stop(); err != nil {
			return err
		}
	}
	if app.fileHelper != nil {
		if err := app.fileHelper.Stop(); err != nil {
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
