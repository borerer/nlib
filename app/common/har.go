package common

import (
	"encoding/json"
	"net/http"

	nlibshared "github.com/borerer/nlib-shared/go"
)

const (
	ContentTypeTextPlain       = "text/plain"
	ContentTypeApplicationJSON = "application/json"
)

func NewResponse(statusCode int, content string, contentType string) *nlibshared.Response {
	res := &nlibshared.Response{}
	res.Status = int64(statusCode)
	res.Content = nlibshared.Content{
		Text: &content,
	}
	res.Headers = append(res.Headers, nlibshared.Header{
		Name:  "Content-Type",
		Value: contentType,
	})
	return res
}

var Err404 = NewResponse(http.StatusNotFound, http.StatusText(http.StatusNotFound), "")
var Err405 = NewResponse(http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed), "")
var Err500 = NewResponse(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), "")

func Text(s string) *nlibshared.Response {
	return NewResponse(http.StatusOK, s, ContentTypeTextPlain)
}

func JSON(v interface{}) *nlibshared.Response {
	buf, err := json.Marshal(v)
	if err != nil {
		return Error(err)
	}
	return NewResponse(http.StatusOK, string(buf), ContentTypeApplicationJSON)
}

func Error(err error) *nlibshared.Response {
	return NewResponse(http.StatusInternalServerError, err.Error(), ContentTypeTextPlain)
}

func GetQuery(req *nlibshared.Request, key string) string {
	for _, query := range req.QueryString {
		if query.Name == key {
			return query.Value
		}
	}
	return ""
}

func GetHeader(req *nlibshared.Request, key string) string {
	for _, header := range req.Headers {
		if header.Name == key {
			return header.Value
		}
	}
	return ""
}

func GetCookie(req *nlibshared.Request, key string) string {
	for _, cookie := range req.Cookies {
		if cookie.Name == key {
			return cookie.Value
		}
	}
	return ""
}
