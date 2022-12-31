package helpers

import (
	"encoding/json"
	"io"
	"net/http"

	nlibshared "github.com/borerer/nlib-shared/go"
	"github.com/borerer/nlib/logs"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Abort500(c *gin.Context, err error) {
	logs.GetZapLogger().Error("abort 500", zap.Error(err))
	c.AbortWithError(http.StatusInternalServerError, err)
}

func Abort404(c *gin.Context, err error) {
	logs.GetZapLogger().Error("abort 404", zap.Error(err))
	c.AbortWithError(http.StatusNotFound, err)
}

func Any200(c *gin.Context, v interface{}) {
	switch t := v.(type) {
	case string:
		c.String(http.StatusOK, t)
	default:
		c.JSON(http.StatusOK, t)
	}
}

func QueryToMap(c *gin.Context) map[string]interface{} {
	ret := map[string]interface{}{}
	for k, v := range c.Request.URL.Query() {
		if len(v) == 1 {
			ret[k] = v[0]
		} else {
			ret[k] = v
		}
	}
	return ret
}

func BodyToMap(c *gin.Context) map[string]interface{} {
	var ret map[string]interface{}
	buf, err := io.ReadAll(c.Request.Body)
	if err == nil {
		json.Unmarshal(buf, &ret)
	}
	return ret
}

func GinToHARRequest(c *gin.Context) *nlibshared.Request {
	var req nlibshared.Request
	req.Method = c.Request.Method
	for k, v := range c.Request.URL.Query() {
		if len(v) > 0 {
			req.QueryString = append(req.QueryString, nlibshared.Query{
				Name:  k,
				Value: v[0], // TODO, handle len > 1 cases
			})
		}
	}
	req.URL = c.Request.URL.String()
	return &req
}
