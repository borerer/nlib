package socket

import (
	"sync"

	"github.com/borerer/nlib/logs"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type ClientsManager struct {
	clients sync.Map
}

func NewClientsManager() *ClientsManager {
	m := &ClientsManager{}
	return m
}

func (m *ClientsManager) StartConnection(clientID string, conn *websocket.Conn) error {
	client := m.getSocketClient(clientID)
	client.SetWebSocketConnection(conn)
	return client.ListenWebSocketMessages()
}

func (m *ClientsManager) CallFunction(clientID string, funcName string, params map[string]interface{}) (interface{}, error) {
	client := m.getSocketClient(clientID)
	return client.CallFunction(funcName, params)
}

func (m *ClientsManager) getSocketClient(clientID string) *Client {
	var client *Client
	clientRaw, ok := m.clients.Load(clientID)
	if ok {
		client, ok = clientRaw.(*Client)
		if !ok {
			logs.Warn("unexpected get nlib client error", zap.String("clientID", clientID))
			// fallback to create a new client instance
		}
	}
	if client == nil {
		client = NewClient(clientID)
		m.clients.Store(clientID, client)
	}
	return client
}
