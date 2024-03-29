package api

import (
	"time"

	"github.com/borerer/nlib/logs"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	ginprometheus "github.com/zsais/go-gin-prometheus"
)

func (api *API) createRouter() error {
	r := gin.New()
	p := ginprometheus.NewPrometheus("nlib")
	p.Use(r)
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

	r.Use(static.Serve("/", static.LocalFile("ui", false)))

	r.GET("/api/app/:id/ws", api.websocketHandler)
	r.Any("/api/app/:id/:func", api.appFunctionHandler)

	api.ginRouter = r
	return nil
}
