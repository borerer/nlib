package app

import (
	"net/http"

	"github.com/borerer/nlib/logs"
	"github.com/borerer/nlib/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const ()

func (app *App) functionRegisterHandler(c *gin.Context) {
	var req models.ReqFunctionRegister
	if err := c.BindJSON(&req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	logs.Info("function register", zap.Any("req", req))
	c.JSON(http.StatusOK, ResponseGeneralOK)
}

func (app *App) appFunctionGetHandler(c *gin.Context) {
	appID := c.Param("app")
	funcName := c.Param("func")

	clientRaw, ok := app.nlibClients.Load(appID)
	if !ok {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	if client, ok := clientRaw.(*NLIBClient); ok {
		message := &models.ReqWSCallFunction{
			ID:   uuid.NewString(),
			Func: funcName,
		}
		client.Connection.WriteJSON(message)
	}
}

func (app *App) appFunctionPostHandler(c *gin.Context) {
	appID := c.Param("app")
	funcName := c.Param("func")
	var payload interface{}
	if err := c.BindJSON(&payload); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	clientRaw, ok := app.nlibClients.Load(appID)
	if !ok {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	if client, ok := clientRaw.(*NLIBClient); ok {
		message := &models.ReqWSCallFunction{
			ID:      uuid.NewString(),
			Func:    funcName,
			Payload: payload,
		}
		client.Connection.WriteJSON(message)
	}
}
