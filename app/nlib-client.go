package app

import "github.com/gorilla/websocket"

type NLIBClient struct {
	AppID      string
	Connection *websocket.Conn
}
