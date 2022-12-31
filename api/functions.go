package api

import (
	nlibshared "github.com/borerer/nlib-shared/go"
	"github.com/borerer/nlib/api/helpers"
	"github.com/gin-gonic/gin"
)

func getSimpleFunctionIn(c *gin.Context) *nlibshared.SimpleFunctionIn {
	input := helpers.BodyToMap(c)
	if len(input) == 0 {
		input = helpers.QueryToMap(c)
	}
	return &input
}

func getHARFunctionIn(c *gin.Context) *nlibshared.HARFunctionIn {
	return helpers.GinToHARRequest(c)
}

func (api *API) appFunctionHandler(c *gin.Context) {
	appID := c.Param("id")
	funcName := c.Param("func")
	useHAR := api.appManager.FunctionUseHAR(appID, funcName)
	if useHAR {
		input := getHARFunctionIn(c)
		output, err := api.appManager.CallHARFunction(appID, funcName, input)
		if err != nil {
			helpers.Abort500(c, err)
		}
		c.String(int(output.Status), *output.Content.Text)
	} else {
		input := getSimpleFunctionIn(c)
		output, err := api.appManager.CallSimpleFunction(appID, funcName, input)
		if err != nil {
			helpers.Abort500(c, err)
		}
		helpers.Any200(c, output)
	}
}
