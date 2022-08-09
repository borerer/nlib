package app

import (
	"errors"
	"net/http"

	"github.com/borerer/nlib/logs"
	"github.com/borerer/nlib/models"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
	"go.uber.org/zap"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func (app *App) handleWebSocketMessage(appID string, conn *websocket.Conn) error {
	var general models.ReqWSGeneral
	if err := conn.ReadJSON(&general); err != nil {
		return err
	}
	switch general.Type {
	case models.WebSocketTypeStart:
		var message models.ReqWSStart
		if err := mapstructure.Decode(general.Payload, &message); err != nil {
			return err
		}
		logs.Info("websocket start", zap.Any("appID", appID), zap.Any("message", message))
		if _, ok := app.nlibClients.Load(appID); ok {
			return errors.New("TODO")
		}
		client := &NLIBClient{
			AppID:      appID,
			Connection: conn,
		}
		app.nlibClients.Store(appID, client)
	}
	return nil
}

func (app *App) websocketHandler(c *gin.Context) {
	appID := c.Query("app")
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer ws.Close()
	for {
		if err = app.handleWebSocketMessage(appID, ws); err != nil {
			logs.Error("web socket message error", zap.Error(err))
		}
	}
}
