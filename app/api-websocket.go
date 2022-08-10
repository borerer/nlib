package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func (app *App) websocketHandler(c *gin.Context) {
	appID := c.Query("app")
	client := app.GetNLIBClient(appID)
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	client.SetWebSocketConnection(conn)
	if err := client.ListenWebSocketMessages(); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, ResponseGeneralOK)
}
