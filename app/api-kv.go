package app

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func (app *App) getKeyValueHandler(c *gin.Context) {
	appID := c.Query("app")
	key := c.Query("key")
	val, err := app.mongoClient.GetKey(appID, key)
	if errors.Is(err, mongo.ErrNoDocuments) {
		c.AbortWithStatus(http.StatusNotFound)
		return
	} else if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.String(http.StatusOK, val)
}

func (app *App) setKeyValueHandler(c *gin.Context) {
	appID := c.Query("app")
	key := c.Query("key")
	value := c.Query("value")
	err := app.mongoClient.SetKey(appID, key, value)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, ResponseGeneralOK)
}
