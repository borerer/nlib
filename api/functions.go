package api

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (api *API) appFunctionGetHandler(c *gin.Context) {
	clientID := c.Param("id")
	funcName := c.Param("func")
	params := c.Query("params")
	res, err := api.socketManager.CallFunction(clientID, funcName, params)
	if err != nil {
		abort500(c, err)
		return
	}
	c.String(http.StatusOK, res)
}

func (api *API) appFunctionPostHandler(c *gin.Context) {
	clientID := c.Param("id")
	funcName := c.Param("func")
	defer c.Request.Body.Close()
	buf, err := io.ReadAll(c.Request.Body)
	if err != nil {
		abort500(c, err)
		return
	}
	params := string(buf)
	res, err := api.socketManager.CallFunction(clientID, funcName, params)
	if err != nil {
		abort500(c, err)
		return
	}
	c.String(http.StatusOK, res)
}
