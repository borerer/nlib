package app

import (
	"fmt"
	"net/http"

	"github.com/borerer/nlib/database"
	"github.com/borerer/nlib/models"
	"github.com/gin-gonic/gin"
)

func (app *App) getKey(appID string, key string) (string, error) {
	var res []models.DBKeyValue
	colName := fmt.Sprintf("%s_kv", appID)
	if err := app.databaseManager.FindDocuments(colName, database.FilterEquals("key", key), &res); err != nil {
		return "", err
	}
	if len(res) == 0 {
		return "", nil
	}
	doc := res[0]
	return doc.Value, nil
}

func (app *App) setKey(appID string, key string, value string) error {
	colName := fmt.Sprintf("%s_kv", appID)
	doc := models.DBKeyValue{
		Key:   key,
		Value: value,
	}
	if err := app.databaseManager.UpdateDocument(colName, database.FilterEquals("key", key), doc); err != nil {
		return err
	}
	return nil
}

func (app *App) getKeyValueHandler(c *gin.Context) {
	appID := c.Query("app")
	key := c.Query("key")
	val, err := app.getKey(appID, key)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.String(http.StatusOK, val)
}

func (app *App) setKeyValueHandler(c *gin.Context) {
	appID := c.Query("app")
	key := c.Query("key")
	value := c.Query("value")
	err := app.setKey(appID, key, value)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, ResponseGeneralOK)
}
