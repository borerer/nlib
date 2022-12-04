package api

import (
	"errors"
	"net/http"

	"github.com/borerer/nlib/database"
	"github.com/gin-gonic/gin"
)

func (api *API) getKeyValueHandler(c *gin.Context) {
	appID := c.Query("app")
	key := c.Query("key")
	val, err := api.mongoClient.GetKey(appID, key)
	if errors.Is(err, database.ErrNoDocuments) {
		abort404(c, err)
		return
	} else if err != nil {
		abort500(c, err)
		return
	}
	c.String(http.StatusOK, val)
}

func (api *API) setKeyValueHandler(c *gin.Context) {
	appID := c.Query("app")
	key := c.Query("key")
	value := c.Query("value")
	err := api.mongoClient.SetKey(appID, key, value)
	if err != nil {
		abort500(c, err)
		return
	}
	c.JSON(http.StatusOK, ResponseGeneralOK)
}
