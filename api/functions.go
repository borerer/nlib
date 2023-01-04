package api

import (
	"encoding/base64"
	"io"

	nlibshared "github.com/borerer/nlib-shared/go"
	"github.com/borerer/nlib/api/helpers"
	"github.com/gin-gonic/gin"
)

func ginToHAR(c *gin.Context) *nlibshared.Request {
	var req nlibshared.Request
	req.Method = c.Request.Method
	for k, v := range c.Request.Header {
		if len(v) > 0 {
			req.Headers = append(req.Headers, nlibshared.Header{
				Name:  k,
				Value: v[0], // TODO, handle len > 1 cases
			})
		}
	}
	req.HeadersSize = int64(len(req.Headers))
	for _, cookie := range c.Request.Cookies() {
		req.Cookies = append(req.Cookies, nlibshared.Cookie{
			Domain:   &cookie.Domain,
			Expires:  &cookie.RawExpires,
			HTTPOnly: &cookie.HttpOnly,
			Name:     cookie.Name,
			Path:     &cookie.Path,
			Secure:   &cookie.Secure,
			Value:    cookie.Value,
		})
	}
	for k, v := range c.Request.URL.Query() {
		if len(v) > 0 {
			req.QueryString = append(req.QueryString, nlibshared.Query{
				Name:  k,
				Value: v[0], // TODO, handle len > 1 cases
			})
		}
	}
	req.URL = c.Request.URL.String()
	buf, err := io.ReadAll(c.Request.Body)
	if err == nil && len(buf) > 0 {
		s := string(buf)
		req.PostData = &nlibshared.PostData{
			Text: &s,
		}
	}
	return &req
}

func harToGin(c *gin.Context, res *nlibshared.Response) {
	for _, header := range res.Headers {
		c.Header(header.Name, header.Value)
	}
	if *res.Content.Encoding == "base64" {
		buf, err := base64.StdEncoding.DecodeString(*res.Content.Text)
		if err != nil {
			helpers.Abort500(c, err)
			return
		}
		c.Data(int(res.Status), "", buf)
	} else {
		c.String(int(res.Status), *res.Content.Text)
	}
}

func (api *API) appFunctionHandler(c *gin.Context) {
	appID := c.Param("id")
	funcName := c.Param("func")
	req := ginToHAR(c)
	res, err := api.appManager.CallFunction(appID, funcName, req)
	if err != nil {
		helpers.Abort500(c, err)
		return
	}
	harToGin(c, res)
}
