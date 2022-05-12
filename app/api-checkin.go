package app

import (
	"net/http"

	"gitea.home.iloahz.com/iloahz/nlib/constants"
	"gitea.home.iloahz.com/iloahz/nlib/logs"
	"gitea.home.iloahz.com/iloahz/nlib/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (app *App) checkinHandler(c *gin.Context) {
	var req models.RequestCheckin
	if err := c.BindJSON(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	logs.Info("reg", zap.Any("req", req))
	res := &models.ResponseCheckin{
		Version: constants.Version,
	}
	c.JSON(http.StatusOK, res)
}
