package app

import (
	"net/http"

	"gitea.home.iloahz.com/iloahz/nlib/configs"
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

}
