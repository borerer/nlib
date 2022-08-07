package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *App) getKeyValueHandler(c *gin.Context) {
	appID := c.Query("app")
	key := c.Query("key")
	val, err := app.databaseManager.GetKey(appID, key)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.String(http.StatusOK, val)
}
