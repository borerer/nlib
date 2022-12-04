package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/borerer/nlib/models"
	"github.com/gin-gonic/gin"
)

func (api *API) addLog(appID string, level string, message string, details interface{}) error {
	colName := fmt.Sprintf("%s_logs", appID)
	doc := models.DBLogs{
		AppID:     appID,
		Level:     level,
		Message:   message,
		Details:   details,
		Timestamp: time.Now().UnixMilli(),
	}
	if err := api.mongoClient.InsertDocument(colName, doc); err != nil {
		return err
	}
	return nil
}

func (api *API) addLogGetHandler(c *gin.Context) {
	appID := c.Query("app")
	level := "info"
	message := ""
	details := map[string]interface{}{}
	for k, v := range c.Request.URL.Query() {
		if k == "level" {
			if len(v) > 0 {
				level = v[0]
			}
		} else if k == "message" {
			if len(v) > 0 {
				message = v[0]
			}
		} else {
			if len(v) > 1 {
				details[k] = v
			} else if len(v) == 1 {
				details[k] = v[0]
			}
		}
	}
	err := api.addLog(appID, level, message, details)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, ResponseGeneralOK)
}

func (api *API) addLogPostHandler(c *gin.Context) {
	appID := c.Query("app")
	var req models.APIAddLogRequest
	err := c.BindJSON(&req)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	err = api.addLog(appID, req.Level, req.Message, req.Details)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, ResponseGeneralOK)
}
