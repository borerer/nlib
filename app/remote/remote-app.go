package remote

import (
	"io"
	"strings"
	"sync"
	"time"

	nlibshared "github.com/borerer/nlib-shared/go"
	"github.com/borerer/nlib/app/common"
	"github.com/borerer/nlib/logs"
	"github.com/borerer/nlib/utils"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type RemoteApp struct {
	appID                   string
	connection              *websocket.Conn
	messageCh               chan *nlibshared.WebSocketMessage
	closeSignalCh           chan bool
	pendingResponseChannels sync.Map
	registeredFunctions     sync.Map
}

func NewRemoteApp(appID string) *RemoteApp {
	c := &RemoteApp{appID: appID}
	return c
}

func (app *RemoteApp) SetWebSocketConnection(conn *websocket.Conn) {
	app.connection = conn
}

func (app *RemoteApp) handleRequest(message *nlibshared.WebSocketMessage) error {
	switch message.SubType {
	case nlibshared.WebSocketMessageSubTypeRegisterFunction:
		var payloadReq nlibshared.PayloadRegisterFunctionRequest
		if err := utils.DecodeStruct(message.Payload, &payloadReq); err != nil {
			return err
		}
		logs.Info("register function", zap.String("appID", app.appID), zap.String("func", payloadReq.Name))
		app.registeredFunctions.Store(payloadReq.Name, &payloadReq)
		res := &nlibshared.WebSocketMessage{
			MessageID:     uuid.NewString(),
			PairMessageID: message.MessageID,
			Type:          nlibshared.WebSocketMessageTypeResponse,
			SubType:       nlibshared.WebSocketMessageSubTypeRegisterFunction,
			Timestamp:     time.Now().UnixMilli(),
			Payload: nlibshared.PayloadRegisterFunctionResponse{
				Name: payloadReq.Name,
			},
		}
		if err := app.sendMessage(res); err != nil {
			return err
		}
	default:
		logs.Warn("unexpected websocket message", zap.Any("message", message))
	}
	return nil
}

func (app *RemoteApp) handleResponse(message *nlibshared.WebSocketMessage) error {
	if chRaw, ok := app.pendingResponseChannels.LoadAndDelete(message.PairMessageID); ok {
		if ch, ok := chRaw.(chan *nlibshared.WebSocketMessage); ok {
			ch <- message
		}
	}
	return nil
}

func (app *RemoteApp) socketCloseHandler(code int, text string) error {
	logs.Info("websocket close handler called", zap.String("appID", app.appID))
	app.handleClose()
	return nil
}

func (app *RemoteApp) handleClose() {
	logs.Info("handle websocket close", zap.String("appID", app.appID))
	app.connection = nil
	app.closeSignalCh <- true
}

func (app *RemoteApp) receiveMessage(message *nlibshared.WebSocketMessage) error {
	logs.Debug("receive message", zap.String("appID", app.appID), zap.Any("message", message))
	switch message.Type {
	case nlibshared.WebSocketMessageTypeRequest:
		return app.handleRequest(message)
	case nlibshared.WebSocketMessageTypeResponse:
		return app.handleResponse(message)
	default:
		logs.Warn("unexpected websocket message", zap.Any("message", message))
		return nil
	}
}

func (app *RemoteApp) readMessages() {
	for {
		var message nlibshared.WebSocketMessage
		if err := app.connection.ReadJSON(&message); err != nil {
			if app.connection == nil {
				// no-op
			} else {
				// errUnexpectedEOF
				if strings.Contains(err.Error(), io.ErrUnexpectedEOF.Error()) {
					app.handleClose()
				} else {
					logs.Error("read websocket message error", zap.String("appID", app.appID), zap.Error(err))
				}
			}
			return
		}
		app.messageCh <- &message
	}
}

func (app *RemoteApp) ListenWebSocketMessages() error {
	logs.Info("websocket connected", zap.String("appID", app.appID))
	app.messageCh = make(chan *nlibshared.WebSocketMessage)
	app.closeSignalCh = make(chan bool)
	app.connection.SetCloseHandler(app.socketCloseHandler)

	go app.readMessages()

	for {
		select {
		case <-app.closeSignalCh:
			return nil
		case message := <-app.messageCh:
			app.receiveMessage(message)
		}
	}
}

func (app *RemoteApp) sendMessage(message *nlibshared.WebSocketMessage) error {
	logs.Debug("send message", zap.String("appID", app.appID), zap.Any("message", message))
	if err := app.connection.WriteJSON(message); err != nil {
		return err
	}
	return nil
}

func (app *RemoteApp) SendWebSocketMessage(subType string, payload interface{}) (*nlibshared.WebSocketMessage, error) {
	message := &nlibshared.WebSocketMessage{
		MessageID: uuid.NewString(),
		Type:      nlibshared.WebSocketMessageTypeRequest,
		SubType:   subType,
		Timestamp: time.Now().UnixMilli(),
		Payload:   payload,
	}
	ch := make(chan *nlibshared.WebSocketMessage, 1)
	app.pendingResponseChannels.Store(message.MessageID, ch)
	if err := app.sendMessage(message); err != nil {
		return nil, err
	}
	res := <-ch
	return res, nil
}

func (app *RemoteApp) CallFunction(name string, req *nlibshared.Request) *nlibshared.Response {
	funcReq := &nlibshared.PayloadCallFunctionRequest{
		Name:    name,
		Request: *req,
	}
	res, err := app.SendWebSocketMessage(nlibshared.WebSocketMessageSubTypeCallFunction, funcReq)
	if err != nil {
		return common.Error(err)
	}
	var funcRes nlibshared.PayloadCallFunctionResponse
	if err = utils.DecodeStruct(res.Payload, &funcRes); err != nil {
		return common.Error(err)
	}
	return &funcRes.Response
}
