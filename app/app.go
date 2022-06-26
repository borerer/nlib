package app

import (
	"context"
	"net/http"

	"github.com/borerer/nlib/configs"
	"github.com/borerer/nlib/utils"
	"github.com/gin-gonic/gin"
)

type App struct {
	config *configs.AppConfig

	ginRouter  *gin.Engine
	httpServer *http.Server
}

func NewApp(config *configs.AppConfig) *App {
	app := &App{
		config: config,
	}
	return app
}

func (app *App) Start() {
	app.startAPI()
}

func (app *App) Stop() {
	if app.httpServer != nil {
		utils.Must(app.httpServer.Shutdown(context.Background()))
	}
}
