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
	pendingResponseChannels sync.Map
}

func NewNLIBClient(appID string) *NLIBClient {
	return &NLIBClient{appID: appID}
}

func (c *NLIBClient) SetWebSocketConnection(conn *websocket.Conn) {
	c.connection = conn
}

func (c *NLIBClient) ListenWebSocketMessages() error {
	defer c.connection.Close()

	handleRequest := func(req *models.WebSocketMessage) error {
		switch req.SubType {
		case models.WebSocketSubTypeStart:
			logs.Info("websocket start", zap.String("appID", c.appID))
			res := &models.WebSocketMessage{
				MessageID:     uuid.NewString(),
				PairMessageID: req.MessageID,
				Type:          models.WebSocketTypeResponse,
				SubType:       models.WebSocketSubTypeStart,
				Timestamp:     time.Now().UnixMilli(),
				Payload:       nil,
			}
			_ = c.connection.WriteJSON(res)
		}
		return nil
	}

	handleResponse := func(res *models.WebSocketMessage) error {
		if chRaw, ok := c.pendingResponseChannels.Load(res.PairMessageID); ok {
			if ch, ok := chRaw.(chan *models.WebSocketMessage); ok {
				ch <- res
			}
		}
		return nil
	}

	for {
		var message models.WebSocketMessage
		if err := c.connection.ReadJSON(&message); err != nil {
			return err
		}
		switch message.Type {
		case models.WebSocketTypeRequest:
			handleRequest(&message)
		case models.WebSocketTypeResponse:
			handleResponse(&message)
		default:
			logs.Warn("unexpected websocket message", zap.Any("message", message))
		}
	}
	return nil
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
