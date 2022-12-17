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

	r.GET("/api/file/get", api.getFileHandler)
	r.PUT("/api/file/put", api.putFileHandler)
	r.DELETE("/api/file/delete", api.deleteFileHandler)
	r.GET("/api/file/stats", api.fileStatsHandler)
	r.GET("/api/file/list", api.listFolderHandler)

	r.GET("/api/kv/get", api.getKeyValueHandler)
	r.PUT("/api/kv/set", api.setKeyValueHandler)

	r.GET("/api/db/:id")
	r.PUT("/api/db/:id")

	r.GET("/api/ws", api.websocketHandler)

	r.GET("/api/logs", api.addLogGetHandler)
	r.POST("/api/logs", api.addLogPostHandler)

	r.GET("/api/remote/:id/:func", api.appFunctionGetHandler)
	r.POST("/api/remote/:id/:func", api.appFunctionPostHandler)

	api.ginRouter = r
	return nil
}
