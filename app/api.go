package app

import (
	"fmt"
	"net/http"
	"time"

	"gitea.home.iloahz.com/iloahz/nlib/logs"
	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (app *App) createRouter() error {
	r := gin.New()
	zapLogger := logs.GetZapLogger()
	r.Use(ginzap.Ginzap(zapLogger, time.RFC3339, true))
	r.Use(ginzap.RecoveryWithZap(zapLogger, true))
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	r.GET("/", app.getHomepageHandler)

	r.POST("/api/checkin", app.checkinHandler)

	app.ginRouter = r
	return nil
}

func (app *App) startAPI() error {
	if err := app.createRouter(); err != nil {
		return err
	}
	listenAddr := app.getListenAddr()
	app.httpServer = &http.Server{
		Addr:    listenAddr,
		Handler: app.ginRouter,
	}
	go func() {
		logs.Info("listen and server api", zap.String("listen", listenAddr))
		if err := app.httpServer.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				logs.Error("api listen and serve error", zap.Error(err))
			}
		}
	}()
	return nil
}

func (app *App) getListenAddr() string {
	addr := app.config.Addr
	port := app.config.Port
	return fmt.Sprintf("%s:%s", addr, port)
}
