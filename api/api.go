package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/borerer/nlib/app"
	"github.com/borerer/nlib/configs"
	"github.com/borerer/nlib/logs"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type API struct {
	config     *configs.APIConfig
	appManager *app.AppManager
	ginRouter  *gin.Engine
	httpServer *http.Server
}

func NewAPI(config *configs.APIConfig, appManager *app.AppManager) *API {
	api := &API{}
	api.config = config
	api.appManager = appManager
	return api
}

func (api *API) Start() error {
	if err := api.appManager.Start(); err != nil {
		return err
	}
	if err := api.createRouter(); err != nil {
		return err
	}
	listenAddr := api.getListenAddr()
	api.httpServer = &http.Server{
		Addr:    listenAddr,
		Handler: api.ginRouter,
	}
	go func() {
		logs.Info("listen and server api", zap.String("listen", listenAddr))
		if err := api.httpServer.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				logs.Error("api listen and serve error", zap.Error(err))
			}
		}
	}()
	return nil
}

func (api *API) Stop() error {
	if api.httpServer != nil {
		if err := api.httpServer.Shutdown(context.Background()); err != nil {
			return err
		}
	}
	if err := api.appManager.Stop(); err != nil {
		return err
	}
	return nil
}

var (
	ResponseGeneralOK = gin.H{
		"status": "ok",
	}
)

func (a *API) getListenAddr() string {
	addr := a.config.Addr
	port := a.config.Port
	return fmt.Sprintf("%s:%s", addr, port)
}
