package logs

import (
	"encoding/base64"
	"encoding/json"
	"strconv"

	nlibshared "github.com/borerer/nlib-shared/go"
	"github.com/borerer/nlib/app/builtin/logs/database"
	"github.com/borerer/nlib/app/common"
	"github.com/borerer/nlib/configs"
)

type LogsApp struct {
	config      *configs.LogsConfig
	mongoClient *database.MongoClient
}

func NewLogsApp(config *configs.LogsConfig) *LogsApp {
	return &LogsApp{
		config: config,
	}
}

func (app *LogsApp) Start() error {
	app.mongoClient = database.NewMongoClient(&app.config.Mongo)
	if err := app.mongoClient.Start(); err != nil {
		return err
	}
	return nil
}

func (app *LogsApp) Stop() error {
	return nil
}

func (app *LogsApp) AppID() string {
	return "logs"
}

func (app *LogsApp) CallFunction(name string, req *nlibshared.Request) *nlibshared.Response {
	switch name {
	case "log":
		return app.log(req)
	case "debug":
		return app.debug(req)
	case "info":
		return app.info(req)
	case "warn":
		return app.warn(req)
	case "error":
		return app.error_(req)
	case "get":
		return app.get(req)
	default:
		return common.Err404
	}
}

func getQuery(req *nlibshared.Request, key string) string {
	for _, query := range req.QueryString {
		if query.Name == key {
			return query.Value
		}
	}
	return ""
}

func getQueryAsInt(req *nlibshared.Request, key string) int {
	val := getQuery(req, key)
	i, err := strconv.Atoi(val)
	if err != nil {
		return 0
	}
	return i
}

func (app *LogsApp) logGET(req *nlibshared.Request) *nlibshared.Response {
	level := getQuery(req, "level")
	if len(level) == 0 {
		level = "info"
	}
	message := getQuery(req, "message")
	err := app.mongoClient.AddLogs(level, message, nil)
	if err != nil {
		return common.Error(err)
	}
	return common.Text("ok")
}

func (app *LogsApp) logPOST(req *nlibshared.Request) *nlibshared.Response {
	parseLog := func(req *nlibshared.Request) (string, string, map[string]interface{}) {
		if req.PostData != nil && req.PostData.Text != nil {
			buf, err := base64.StdEncoding.DecodeString(*req.PostData.Text)
			if err == nil {
				var j map[string]interface{}
				err := json.Unmarshal(buf, &j)
				if err == nil {
					level := "info"
					levelRaw, ok := j["level"]
					if ok {
						level = levelRaw.(string)
					}
					message := j["message"].(string)
					detailsRaw, ok := j["details"]
					if !ok {
						return level, message, nil
					}
					details := detailsRaw.(map[string]interface{})
					return level, message, details
				}
			}
		}
		return "", "", nil
	}
	level, message, details := parseLog(req)
	err := app.mongoClient.AddLogs(level, message, details)
	if err != nil {
		return common.Error(err)
	}
	return common.Text("ok")
}

func (app *LogsApp) log(req *nlibshared.Request) *nlibshared.Response {
	if req.Method == "GET" {
		return app.logGET(req)
	} else if req.Method == "POST" {
		return app.logPOST(req)
	}
	return common.Err405
}

func (app *LogsApp) debug(req *nlibshared.Request) *nlibshared.Response {
	req.QueryString = append(req.QueryString, nlibshared.Query{Name: "level", Value: "debug"})
	return app.log(req)
}

func (app *LogsApp) info(req *nlibshared.Request) *nlibshared.Response {
	req.QueryString = append(req.QueryString, nlibshared.Query{Name: "level", Value: "info"})
	return app.log(req)
}

func (app *LogsApp) warn(req *nlibshared.Request) *nlibshared.Response {
	req.QueryString = append(req.QueryString, nlibshared.Query{Name: "level", Value: "warn"})
	return app.log(req)
}

func (app *LogsApp) error_(req *nlibshared.Request) *nlibshared.Response {
	req.QueryString = append(req.QueryString, nlibshared.Query{Name: "level", Value: "error"})
	return app.log(req)
}

func (app *LogsApp) get(req *nlibshared.Request) *nlibshared.Response {
	n := getQueryAsInt(req, "n")
	if n == 0 {
		n = 100
	}
	skip := getQueryAsInt(req, "skip")
	logs, err := app.mongoClient.GetLogs(n, skip)
	if err != nil {
		return common.Error(err)
	}
	return common.JSON(logs)
}
