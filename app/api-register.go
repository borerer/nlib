package app

import (
	"net/http"

	"github.com/borerer/nlib/constants"
	"github.com/borerer/nlib/logs"
	"github.com/borerer/nlib/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (app *App) registerHandler(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.BindJSON(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	logs.Info("reg", zap.Any("req", req))
	res := &models.RegisterResponse{
		Version: constants.Version,
	}
	c.JSON(http.StatusOK, res)
}
