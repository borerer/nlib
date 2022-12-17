package api

import (
	"net/http"

	"github.com/borerer/nlib/logs"
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
	clientID := c.Query("id")
	logs.Info("websocket connected", zap.String("clientID", clientID))
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		abort500(c, err)
		return
	}
	if err := api.socketManager.StartConnection(clientID, conn); err != nil {
		abort500(c, err)
		return
	}
	c.JSON(http.StatusOK, ResponseGeneralOK)
}
