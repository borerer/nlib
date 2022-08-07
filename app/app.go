package app

import (
	"context"
	"net/http"

	"github.com/borerer/nlib/configs"
	"github.com/borerer/nlib/file"
	"github.com/borerer/nlib/logs"
	"github.com/borerer/nlib/utils"
	"github.com/gin-gonic/gin"
)

type App struct {
	config *configs.AppConfig

	ginRouter  *gin.Engine
	httpServer *http.Server

	fileHelper file.FileHelper
}

func NewApp(config *configs.AppConfig) *App {
	app := &App{
		config:     config,
		fileHelper: file.NewFileHelper(&config.File),
	}
	return app
}

func (app *App) Start() error {
	logs.Info("app start")
	if err := app.startAPI(); err != nil {
		return err
	}
	return nil
}

func (app *App) Stop() {
	logs.Info("app stop")
	if app.httpServer != nil {
		utils.Must(app.httpServer.Shutdown(context.Background()))
	}
}
