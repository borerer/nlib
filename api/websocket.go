package api

import (
	"net/http"

	"github.com/borerer/nlib/logs"
	"github.com/borerer/nlib/socket"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func (api *API) websocketHandler(c *gin.Context) {
	appID := c.Query("app")
	logs.Info("websocket connected", zap.String("appID", appID))
	client := api.getApp(appID)
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		abort500(c, err)
		return
	}
	client.SetWebSocketConnection(conn)
	if err := client.ListenWebSocketMessages(); err != nil {
		abort500(c, err)
		return
	}
	c.JSON(http.StatusOK, ResponseGeneralOK)
}

func (api *API) getApp(appID string) *socket.App {
	var app *socket.App
	appRaw, ok := api.clients.Load(appID)
	if ok {
		app, ok = appRaw.(*socket.App)
		if !ok {
			logs.Warn("unexpected get nlib client error", zap.String("appID", appID))
			// fallback to create a new client instance
		}
	}
	if app == nil {
		app = socket.NewApp(appID)
		api.clients.Store(appID, app)
	}
	return app
}
