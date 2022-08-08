package app

import (
	"fmt"
	"net/http"

	"github.com/borerer/nlib/database"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func getColName(appID string) string {
	return fmt.Sprintf("%s_kv", appID)
}

func (app *App) getKey(appID string, key string) (string, error) {
	var res []bson.M
	colName := getColName(appID)
	if err := app.databaseManager.FindDocuments(colName, database.FilterEquals("key", key), &res); err != nil {
		return "", err
	}
	if len(res) == 0 {
		return "", nil
	}
	doc := res[0]
	if val, ok := doc["value"]; ok {
		if valStr, ok := val.(string); ok {
			return valStr, nil
		}
	}
	return "", nil
}

func (app *App) setKey(appID string, key string, value string) error {
	colName := getColName(appID)
	doc := bson.M{
		"key":   key,
		"value": value,
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
