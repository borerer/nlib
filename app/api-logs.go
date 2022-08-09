package app

import (
	"fmt"
	"net/http"
	"time"

	"github.com/borerer/nlib/models"
	"github.com/gin-gonic/gin"
)

func (app *App) addLog(appID string, message string, structuredMessage interface{}, level string) error {
	colName := fmt.Sprintf("%s_logs", appID)
	doc := models.DBLogs{
		AppID:             appID,
		Message:           message,
		StructuredMessage: structuredMessage,
		Level:             level,
		Timestamp:         time.Now().UnixMilli(),
	}
	if err := app.databaseManager.InsertDocument(colName, doc); err != nil {
		return err
	}
	return nil
}

func (app *App) addLogHandler(c *gin.Context) {
	appID := c.Query("app")
	message := c.Query("message")
	var structuredMessage interface{}
	_ = c.Bind(&structuredMessage)
	level := c.Query("level")
	err := app.addLog(appID, message, structuredMessage, level)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, ResponseGeneralOK)
}
