package app

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *App) appFunctionGetHandler(c *gin.Context) {
	appID := c.Param("app")
	client := app.GetNLIBClient(appID)
	funcName := c.Param("func")
	params := c.Query("params")
	res, err := client.CallFunction(funcName, params)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.String(http.StatusOK, res)
}

func (app *App) appFunctionPostHandler(c *gin.Context) {
	appID := c.Param("app")
	client := app.GetNLIBClient(appID)
	funcName := c.Param("func")
	defer c.Request.Body.Close()
	buf, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	params := string(buf)
	res, err := client.CallFunction(funcName, params)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.String(http.StatusOK, res)
}
