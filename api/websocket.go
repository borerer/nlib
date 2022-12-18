package api

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

func (api *API) websocketHandler(c *gin.Context) {
	appID := c.Param("id")
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		abort500(c, err)
		return
	}
	if err := api.socketManager.StartConnection(appID, conn); err != nil {
		abort500(c, err)
		return
	}
}
