package server

import (
	"github.com/borerer/nlib/api"
	"github.com/borerer/nlib/app"
	"github.com/borerer/nlib/configs"
	"github.com/borerer/nlib/logs"
	"github.com/borerer/nlib/utils"
	"go.uber.org/zap"
)

type Server struct {
	config     *configs.ServerConfig
	api        *api.API
	appManager *app.AppManager
}

func NewServer(config *configs.ServerConfig) *Server {
	server := &Server{}
	server.config = config
	server.appManager = app.NewAppManager(&config.Builtin)
	server.api = api.NewAPI(&config.API, server.appManager)
	return server
}

func (server *Server) Start() error {
	logs.Info("start server")
	if err := server.appManager.Start(); err != nil {
		return err
	}
	if err := server.api.Start(); err != nil {
		return err
	}
	return nil
}

func (server *Server) Stop() error {
	logs.Info("stop server")
	if err := server.api.Stop(); err != nil {
		return err
	}
	if err := server.appManager.Stop(); err != nil {
		return err
	}
	return nil
}

func waitForever() {
	ch := make(chan bool)
	<-ch
}

func Run(config *configs.ServerConfig) {
	utils.Must(logs.Init(config))
	logs.Info("run", zap.Any("config", config))
	server := NewServer(config)
	utils.Must(server.Start())
	waitForever()
	server.Stop()
}
