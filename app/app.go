package app

import "gitea.home.iloahz.com/iloahz/nlib/configs"

type App struct {
	config *configs.AppConfig
}

func NewApp(config *configs.AppConfig) *App {
	app := &App{
		config: config,
	}
	return app
}

func (app *App) Start() {

}

func (app *App) Stop() {

}
