package socket

import (
	"sync"
	"time"

	"github.com/borerer/nlib/logs"
	"github.com/borerer/nlib/models"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
	"go.uber.org/zap"
)

type Client struct {
	clientID                string
	connection              *websocket.Conn
	messageCh               chan *models.WebSocketMessage
	closeSignalCh           chan bool
	pendingResponseChannels sync.Map
}

func NewClient(clientID string) *Client {
	c := &Client{clientID: clientID}
	return c
}

func (client *Client) SetWebSocketConnection(conn *websocket.Conn) {
	client.connection = conn
}

func (client *Client) handleRequest(message *models.WebSocketMessage) error {
	switch message.SubType {
	case models.WebSocketSubTypeStart:
		logs.Info("websocket start", zap.String("clientID", client.clientID))
		res := &models.WebSocketMessage{
			MessageID:     uuid.NewString(),
			PairMessageID: message.MessageID,
			Type:          models.WebSocketTypeResponse,
			SubType:       models.WebSocketSubTypeStart,
			Timestamp:     time.Now().UnixMilli(),
			Payload:       nil,
		}
		_ = client.connection.WriteJSON(res)
	}
	return nil
}

func (client *Client) handleResponse(message *models.WebSocketMessage) error {
	if chRaw, ok := client.pendingResponseChannels.Load(message.PairMessageID); ok {
		if ch, ok := chRaw.(chan *models.WebSocketMessage); ok {
			ch <- message
		}
	}
	return nil
}

func (client *Client) handleClose(code int, text string) error {
	logs.Info("websocket close", zap.String("clientID", client.clientID))
	client.connection = nil
	client.closeSignalCh <- true
	return nil
}

func (client *Client) handleMessage(message *models.WebSocketMessage) error {
	switch message.Type {
	case models.WebSocketTypeRequest:
		return client.handleRequest(message)
	case models.WebSocketTypeResponse:
		return client.handleResponse(message)
	default:
		logs.Warn("unexpected websocket message", zap.Any("message", message))
		return nil
	}
}

func (client *Client) readMessages() {
	for {
		var message models.WebSocketMessage
		if err := client.connection.ReadJSON(&message); err != nil {
			if client.connection == nil {
				// no-op
			} else {
				logs.Error("read websocket message error", zap.String("clientID", client.clientID), zap.Error(err))
			}
			return
		}
		client.messageCh <- &message
	}
}

func (client *Client) ListenWebSocketMessages() error {
	client.messageCh = make(chan *models.WebSocketMessage)
	client.closeSignalCh = make(chan bool)
	client.connection.SetCloseHandler(client.handleClose)

	go client.readMessages()

	for {
		select {
		case <-client.closeSignalCh:
			return nil
		case message := <-client.messageCh:
			client.handleMessage(message)
		}
	}
}

func (client *Client) SendWebSocketMessage(subType string, payload interface{}) (*models.WebSocketMessage, error) {
	message := &models.WebSocketMessage{
		MessageID: uuid.NewString(),
		Type:      models.WebSocketTypeRequest,
		SubType:   subType,
		Timestamp: time.Now().UnixMilli(),
		Payload:   payload,
	}
	ch := make(chan *models.WebSocketMessage, 1)
	client.pendingResponseChannels.Store(message.MessageID, ch)
	if err := client.connection.WriteJSON(message); err != nil {
		return nil, err
	}
	res := <-ch
	client.pendingResponseChannels.Delete(message.MessageID)
	return res, nil
}

func (client *Client) CallFunction(funcName string, params string) (string, error) {
	funcReq := &models.WebSocketCallFunctionReq{
		FuncName: funcName,
		Params:   params,
	}
	res, err := client.SendWebSocketMessage(models.WebSocketSubTypeCallFunction, funcReq)
	if err != nil {
		return "", err
	}
	var funcRes models.WebSocketCallFunctionRes
	if err = mapstructure.Decode(res.Payload, &funcRes); err != nil {
		return "", err
	}
	return funcRes.Response, nil
}
