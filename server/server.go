package server

import (
	"github.com/borerer/nlib/api"
	"github.com/borerer/nlib/configs"
	"github.com/borerer/nlib/logs"
	"github.com/borerer/nlib/utils"
	"go.uber.org/zap"
)

type Server struct {
	config *configs.ServerConfig
	api    *api.API
}

func NewServer(config *configs.ServerConfig) *Server {
	server := &Server{}
	server.config = config
	server.api = api.NewAPI(&config.API)
	return server
}

func (server *Server) Start() error {
	logs.Info("start server")
	if err := server.api.Start(); err != nil {
		return err
	}
	return nil
}

func (server *Server) Stop() error {
	logs.Info("stop server")
	if server.api != nil {
		if err := server.api.Stop(); err != nil {
			return err
		}
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
