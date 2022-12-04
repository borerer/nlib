package api

import (
	"net/http"

	"github.com/borerer/nlib/logs"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func abort500(c *gin.Context, err error) {
	logs.GetZapLogger().Error("abort 500", zap.Error(err))
	c.AbortWithError(http.StatusInternalServerError, err)
}

func abort404(c *gin.Context, err error) {
	logs.GetZapLogger().Error("abort 404", zap.Error(err))
	c.AbortWithError(http.StatusNotFound, err)
}
