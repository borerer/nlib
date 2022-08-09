package app

import (
	"fmt"
	"net/http"

	"github.com/borerer/nlib/models"
	"github.com/gin-gonic/gin"
)

func (app *App) addLog(appID string, msg string) error {
	colName := fmt.Sprintf("%s_logs", appID)
	doc := models.DBLogs{
		AppID:   appID,
		Message: msg,
	}
	if err := app.databaseManager.InsertDocument(colName, doc); err != nil {
		return err
	}
	return nil
}

func (app *App) addLogHandler(c *gin.Context) {
	appID := c.Query("app")
	msg := c.Query("msg")
	err := app.addLog(appID, msg)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, ResponseGeneralOK)
}
