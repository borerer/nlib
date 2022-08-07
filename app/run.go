package app

import (
	"github.com/borerer/nlib/configs"
	"github.com/borerer/nlib/logs"
	"github.com/borerer/nlib/utils"
	"go.uber.org/zap"
)

func waitingForever() {
	ch := make(chan bool)
	<-ch
}

func Run(config *configs.AppConfig) {
	utils.Must(logs.Init(config))
	logs.Info("run", zap.Any("config", config))
	app := NewApp(config)
	app.Start()
	waitingForever()
}
