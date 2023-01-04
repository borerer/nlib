package api

import (
	nlibshared "github.com/borerer/nlib-shared/go"
	"github.com/borerer/nlib/api/helpers"
	"github.com/gin-gonic/gin"
)

func getFunctionIn(c *gin.Context) *nlibshared.Request {
	return helpers.GinToHARRequest(c)
}

func (api *API) appFunctionHandler(c *gin.Context) {
	appID := c.Param("id")
	funcName := c.Param("func")
	input := getFunctionIn(c)
	output, err := api.appManager.CallFunction(appID, funcName, input)
	if err != nil {
		helpers.Abort500(c, err)
	}
	c.String(int(output.Status), *output.Content.Text)
}
