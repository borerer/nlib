package app

import (
	"net/http"

	"github.com/borerer/nlib/logs"
	"github.com/borerer/nlib/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (app *App) logHandler(c *gin.Context) {
	var req models.LogRequest
	if err := c.BindJSON(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	logs.Info("reg", zap.Any("req", req))
	c.JSON(http.StatusOK, models.GeneralOK)
}
