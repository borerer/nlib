package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// func (app *App) saveRegisteredFunctionToDatabase(appID string, req *models.WebSocketRegisterFunction) error {
// 	doc := models.DBAppFunction{
// 		AppID: appID,
// 		Func:  req.Func,
// 	}
// 	col := fmt.Sprintf("%s_functions", appID)
// 	if err := app.databaseManager.InsertDocument(col, doc); err != nil {
// 		return err
// 	}
// 	return nil
// }

func (app *App) appFunctionGetHandler(c *gin.Context) {
	appID := c.Param("app")
	client := app.GetNLIBClient(appID)
	funcName := c.Param("func")
	res, err := client.CallFunction(funcName, "")
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.String(http.StatusOK, res)
}

// func (app *App) appFunctionPostHandler(c *gin.Context) {
// 	appID := c.Param("app")
// 	funcName := c.Param("func")
// 	var payload interface{}
// 	if err := c.BindJSON(&payload); err != nil {
// 		c.AbortWithError(http.StatusBadRequest, err)
// 		return
// 	}

// 	clientRaw, ok := app.nlibClients.Load(appID)
// 	if !ok {
// 		c.AbortWithStatus(http.StatusNotFound)
// 		return
// 	}
// 	if client, ok := clientRaw.(*NLIBClient); ok {
// 		message := &models.WebSocketCallFunction{
// 			ID:      uuid.NewString(),
// 			Func:    funcName,
// 			Payload: payload,
// 		}
// 		client.connection.WriteJSON(message)
// 	}
// }
