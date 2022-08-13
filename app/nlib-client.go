package app

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

type NLIBClient struct {
	appID                   string
	connection              *websocket.Conn
	messageCh               chan *models.WebSocketMessage
	closeSignalCh           chan bool
	pendingResponseChannels sync.Map
}

func NewNLIBClient(appID string) *NLIBClient {
	c := &NLIBClient{appID: appID}
	return c
}

func (c *NLIBClient) SetWebSocketConnection(conn *websocket.Conn) {
	c.connection = conn
}

func (c *NLIBClient) handleRequest(message *models.WebSocketMessage) error {
	switch message.SubType {
	case models.WebSocketSubTypeStart:
		logs.Info("websocket start", zap.String("appID", c.appID))
		res := &models.WebSocketMessage{
			MessageID:     uuid.NewString(),
			PairMessageID: message.MessageID,
			Type:          models.WebSocketTypeResponse,
			SubType:       models.WebSocketSubTypeStart,
			Timestamp:     time.Now().UnixMilli(),
			Payload:       nil,
		}
		_ = c.connection.WriteJSON(res)
	}
	return nil
}

func (c *NLIBClient) handleResponse(message *models.WebSocketMessage) error {
	if chRaw, ok := c.pendingResponseChannels.Load(message.PairMessageID); ok {
		if ch, ok := chRaw.(chan *models.WebSocketMessage); ok {
			ch <- message
		}
	}
	return nil
}

func (c *NLIBClient) handleClose(code int, text string) error {
	logs.Info("websocket close", zap.String("appID", c.appID))
	c.connection = nil
	c.closeSignalCh <- true
	return nil
}

func (c *NLIBClient) handleMessage(message *models.WebSocketMessage) error {
	switch message.Type {
	case models.WebSocketTypeRequest:
		return c.handleRequest(message)
	case models.WebSocketTypeResponse:
		return c.handleResponse(message)
	default:
		logs.Warn("unexpected websocket message", zap.Any("message", message))
		return nil
	}
}

func (c *NLIBClient) readMessages() {
	for {
		var message models.WebSocketMessage
		if err := c.connection.ReadJSON(&message); err != nil {
			if c.connection == nil {
				// no-op
			} else {
				logs.Error("read websocket message error", zap.String("appID", c.appID), zap.Error(err))
			}
			return
		}
		c.messageCh <- &message
	}
}

func (c *NLIBClient) ListenWebSocketMessages() error {
	c.messageCh = make(chan *models.WebSocketMessage)
	c.closeSignalCh = make(chan bool)
	c.connection.SetCloseHandler(c.handleClose)

	go c.readMessages()

	for {
		select {
		case <-c.closeSignalCh:
			return nil
		case message := <-c.messageCh:
			c.handleMessage(message)
		}
	}
}

func (c *NLIBClient) SendWebSocketMessage(subType string, payload interface{}) (*models.WebSocketMessage, error) {
	message := &models.WebSocketMessage{
		MessageID: uuid.NewString(),
		Type:      models.WebSocketTypeRequest,
		SubType:   subType,
		Timestamp: time.Now().UnixMilli(),
		Payload:   payload,
	}
	ch := make(chan *models.WebSocketMessage, 1)
	c.pendingResponseChannels.Store(message.MessageID, ch)
	if err := c.connection.WriteJSON(message); err != nil {
		return nil, err
	}
	res := <-ch
	c.pendingResponseChannels.Delete(message.MessageID)
	return res, nil
}

func (c *NLIBClient) CallFunction(funcName string, params string) (string, error) {
	funcReq := &models.WebSocketCallFunctionReq{
		FuncName: funcName,
		Params:   params,
	}
	res, err := c.SendWebSocketMessage(models.WebSocketSubTypeCallFunction, funcReq)
	if err != nil {
		return "", err
	}
	var funcRes models.WebSocketCallFunctionRes
	if err = mapstructure.Decode(res.Payload, &funcRes); err != nil {
		return "", err
	}
	return funcRes.Response, nil
}
