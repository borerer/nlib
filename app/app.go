package app

import (
	"io"
	"strings"
	"sync"
	"time"

	nlibshared "github.com/borerer/nlib-shared/go"
	"github.com/borerer/nlib/logs"
	"github.com/borerer/nlib/utils"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type App struct {
	appID                   string
	connection              *websocket.Conn
	messageCh               chan *nlibshared.WebSocketMessage
	closeSignalCh           chan bool
	pendingResponseChannels sync.Map
	registeredFunctions     sync.Map
}

func NewApp(appID string) *App {
	c := &App{appID: appID}
	return c
}

func (app *App) SetWebSocketConnection(conn *websocket.Conn) {
	app.connection = conn
}

func (app *App) handleRequest(message *nlibshared.WebSocketMessage) error {
	switch message.SubType {
	case nlibshared.WebSocketMessageSubTypeRegisterFunction:
		var payloadReq nlibshared.PayloadRegisterFunctionRequest
		if err := utils.DecodeStruct(message.Payload, &payloadReq); err != nil {
			return err
		}
		logs.Info("register function", zap.String("appID", app.appID), zap.String("func", payloadReq.Name), zap.Bool("useHAR", payloadReq.UseHAR))
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

func (app *App) handleResponse(message *nlibshared.WebSocketMessage) error {
	if chRaw, ok := app.pendingResponseChannels.LoadAndDelete(message.PairMessageID); ok {
		if ch, ok := chRaw.(chan *nlibshared.WebSocketMessage); ok {
			ch <- message
		}
	}
	return nil
}

func (app *App) socketCloseHandler(code int, text string) error {
	logs.Info("websocket close handler called", zap.String("appID", app.appID))
	app.handleClose()
	return nil
}

func (app *App) handleClose() {
	logs.Info("handle websocket close", zap.String("appID", app.appID))
	app.connection = nil
	app.closeSignalCh <- true
}

func (app *App) receiveMessage(message *nlibshared.WebSocketMessage) error {
	logs.Debug("receive message", zap.Any("message", message))
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

func (app *App) readMessages() {
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

func (app *App) ListenWebSocketMessages() error {
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

func (app *App) sendMessage(message *nlibshared.WebSocketMessage) error {
	logs.Debug("send message", zap.Any("message", message))
	if err := app.connection.WriteJSON(message); err != nil {
		return err
	}
	return nil
}

func (app *App) SendWebSocketMessage(subType string, payload interface{}) (*nlibshared.WebSocketMessage, error) {
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

func (app *App) FunctionUseHAR(name string) bool {
	raw, ok := app.registeredFunctions.Load(name)
	if !ok {
		return false
	}
	req, ok := raw.(*nlibshared.PayloadRegisterFunctionRequest)
	if !ok {
		return false
	}
	return req.UseHAR
}

func (app *App) CallSimpleFunction(name string, req *nlibshared.SimpleFunctionIn) (nlibshared.SimpleFunctionOut, error) {
	funcReq := &nlibshared.PayloadCallFunctionRequest{
		Name:    name,
		Request: req,
	}
	res, err := app.SendWebSocketMessage(nlibshared.WebSocketMessageSubTypeCallFunction, funcReq)
	if err != nil {
		return nil, err
	}
	var funcRes nlibshared.PayloadCallFunctionResponse
	if err = utils.DecodeStruct(res.Payload, &funcRes); err != nil {
		return nil, err
	}
	var out nlibshared.SimpleFunctionOut
	if err = utils.DecodeStruct(funcRes.Response, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (app *App) CallHARFunction(name string, req *nlibshared.HARFunctionIn) (*nlibshared.HARFunctionOut, error) {
	funcReq := &nlibshared.PayloadCallFunctionRequest{
		Name:    name,
		UseHAR:  true,
		Request: req,
	}
	res, err := app.SendWebSocketMessage(nlibshared.WebSocketMessageSubTypeCallFunction, funcReq)
	if err != nil {
		return nil, err
	}
	var funcRes nlibshared.PayloadCallFunctionResponse
	if err = utils.DecodeStruct(res.Payload, &funcRes); err != nil {
		return nil, err
	}
	var out nlibshared.HARFunctionOut
	if err = utils.DecodeStruct(funcRes.Response, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
