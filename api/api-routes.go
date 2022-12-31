package api

import (
	"time"

	"github.com/borerer/nlib/logs"
	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
)

func (api *API) createRouter() error {
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

	r.GET("/", api.getHomepageHandler)

	r.GET("/api/app/:id/ws", api.websocketHandler)
	r.Any("/api/app/:id/:func", api.appFunctionHandler)

	api.ginRouter = r
	return nil
}
