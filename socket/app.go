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

type App struct {
	appID                   string
	connection              *websocket.Conn
	messageCh               chan *models.WebSocketMessage
	closeSignalCh           chan bool
	pendingResponseChannels sync.Map
}

func NewApp(appID string) *App {
	c := &App{appID: appID}
	return c
}

func (app *App) SetWebSocketConnection(conn *websocket.Conn) {
	app.connection = conn
}

func (app *App) handleRequest(message *models.WebSocketMessage) error {
	switch message.SubType {
	case models.WebSocketSubTypeStart:
		logs.Info("websocket start", zap.String("appID", app.appID))
		res := &models.WebSocketMessage{
			MessageID:     uuid.NewString(),
			PairMessageID: message.MessageID,
			Type:          models.WebSocketTypeResponse,
			SubType:       models.WebSocketSubTypeStart,
			Timestamp:     time.Now().UnixMilli(),
			Payload:       nil,
		}
		_ = app.connection.WriteJSON(res)
	}
	return nil
}

func (app *App) handleResponse(message *models.WebSocketMessage) error {
	if chRaw, ok := app.pendingResponseChannels.Load(message.PairMessageID); ok {
		if ch, ok := chRaw.(chan *models.WebSocketMessage); ok {
			ch <- message
		}
	}
	return nil
}

func (app *App) handleClose(code int, text string) error {
	logs.Info("websocket close", zap.String("appID", app.appID))
	app.connection = nil
	app.closeSignalCh <- true
	return nil
}

func (app *App) handleMessage(message *models.WebSocketMessage) error {
	switch message.Type {
	case models.WebSocketTypeRequest:
		return app.handleRequest(message)
	case models.WebSocketTypeResponse:
		return app.handleResponse(message)
	default:
		logs.Warn("unexpected websocket message", zap.Any("message", message))
		return nil
	}
}

func (app *App) readMessages() {
	for {
		var message models.WebSocketMessage
		if err := app.connection.ReadJSON(&message); err != nil {
			if app.connection == nil {
				// no-op
			} else {
				logs.Error("read websocket message error", zap.String("appID", app.appID), zap.Error(err))
			}
			return
		}
		app.messageCh <- &message
	}
}

func (app *App) ListenWebSocketMessages() error {
	app.messageCh = make(chan *models.WebSocketMessage)
	app.closeSignalCh = make(chan bool)
	app.connection.SetCloseHandler(app.handleClose)

	go app.readMessages()

	for {
		select {
		case <-app.closeSignalCh:
			return nil
		case message := <-app.messageCh:
			app.handleMessage(message)
		}
	}
}

func (app *App) SendWebSocketMessage(subType string, payload interface{}) (*models.WebSocketMessage, error) {
	message := &models.WebSocketMessage{
		MessageID: uuid.NewString(),
		Type:      models.WebSocketTypeRequest,
		SubType:   subType,
		Timestamp: time.Now().UnixMilli(),
		Payload:   payload,
	}
	ch := make(chan *models.WebSocketMessage, 1)
	app.pendingResponseChannels.Store(message.MessageID, ch)
	if err := app.connection.WriteJSON(message); err != nil {
		return nil, err
	}
	res := <-ch
	app.pendingResponseChannels.Delete(message.MessageID)
	return res, nil
}

func (app *App) CallFunction(funcName string, params string) (string, error) {
	funcReq := &models.WebSocketCallFunctionReq{
		FuncName: funcName,
		Params:   params,
	}
	res, err := app.SendWebSocketMessage(models.WebSocketSubTypeCallFunction, funcReq)
	if err != nil {
		return "", err
	}
	var funcRes models.WebSocketCallFunctionRes
	if err = mapstructure.Decode(res.Payload, &funcRes); err != nil {
		return "", err
	}
	return funcRes.Response, nil
}
