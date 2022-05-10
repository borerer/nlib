package app

import (
	"gitea.home.iloahz.com/iloahz/nlib/configs"
	"gitea.home.iloahz.com/iloahz/nlib/logs"
	"gitea.home.iloahz.com/iloahz/nlib/utils"
	"go.uber.org/zap"
)

func waitingForever() {
	ch := make(chan bool)
	<-ch
}

func Run(config *configs.AppConfig) {
	utils.Must(logs.Init(config))
	logs.Info("run app", zap.Any("config", config))
	app := NewApp(config)
	app.Start()
	waitingForever()
}
