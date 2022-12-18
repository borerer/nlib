package api

import (
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
	any200(c, res)
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
	any200(c, res)
}
