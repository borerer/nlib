package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (api *API) appFunctionGetHandler(c *gin.Context) {
	clientID := c.Param("id")
	funcName := c.Param("func")
	params := queryToMap(c)
	res, err := api.socketManager.CallFunction(clientID, funcName, params)
	if err != nil {
		abort500(c, err)
		return
	}
	c.JSON(http.StatusOK, res)
}

func (api *API) appFunctionPostHandler(c *gin.Context) {
	clientID := c.Param("id")
	funcName := c.Param("func")
	var params map[string]interface{}
	c.BindJSON(&params)
	res, err := api.socketManager.CallFunction(clientID, funcName, params)
	if err != nil {
		abort500(c, err)
		return
	}
	c.JSON(http.StatusOK, res)
}
