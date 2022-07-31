package app

import (
	"fmt"
	"net/http"
	"path"
	"time"

	"github.com/borerer/nlib/engines"
	"github.com/borerer/nlib/logs"
	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (app *App) getAppEngine(c *gin.Context) {
	appID := c.Query("app")
	fileBaseDir := path.Join(app.config.FileDir, appID)
	c.Set("app-engine", engines.NewAppEngine(appID, fileBaseDir))
	c.Next()
}

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

	r.GET("/api/file/get", app.getAppEngine, app.getObjectHandler)
	r.PUT("/api/file/put", app.getAppEngine, app.putObjectHandler)
	r.POST("/api/log", app.logHandler)

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
