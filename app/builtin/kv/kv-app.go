package kv

import (
	"encoding/base64"
	"encoding/json"
	"errors"

	nlibshared "github.com/borerer/nlib-shared/go"
	"github.com/borerer/nlib/app/builtin/kv/database"
	"github.com/borerer/nlib/app/common"
)

type KVApp struct {
	mongoURI    string
	mongoClient *database.MongoClient
}

func NewKVApp(mongoURI string) *KVApp {
	return &KVApp{
		mongoURI: mongoURI,
	}
}

func (kv *KVApp) Start() error {
	kv.mongoClient = database.NewMongoClient(kv.mongoURI)
	if err := kv.mongoClient.Start(); err != nil {
		return err
	}
	return nil
}

func (kv *KVApp) Stop() error {
	return nil
}

func (kv *KVApp) AppID() string {
	return "kv"
}

func (kv *KVApp) CallFunction(name string, req *nlibshared.Request) *nlibshared.Response {
	switch name {
	case "get":
		return kv.getKey(req)
	case "set":
		return kv.setKey(req)
	case "recent":
		return kv.getRecent(req)
	default:
		return common.Err404
	}
}

func (kv *KVApp) GetKey(key string) (string, error) {
	val, err := kv.mongoClient.GetKey(key)
	if err != nil {
		return "", err
	}
	return val, nil
}

func (kv *KVApp) getKey(req *nlibshared.Request) *nlibshared.Response {
	key := common.GetQuery(req, "key")
	val, err := kv.GetKey(key)
	if errors.Is(err, database.ErrNoDocuments) {
		return common.Err404
	} else if err != nil {
		return common.Error(err)
	}
	return common.Text(val)
}

func (kv *KVApp) setKeyGET(req *nlibshared.Request) *nlibshared.Response {
	key := common.GetQuery(req, "key")
	value := common.GetQuery(req, "value")
	err := kv.mongoClient.SetKey(key, value)
	if err != nil {
		return common.Error(err)
	}
	return common.Text("ok")
}

func (kv *KVApp) setKeyPOST(req *nlibshared.Request) *nlibshared.Response {
	parseKeyValue := func(req *nlibshared.Request) (string, string) {
		if req.PostData != nil && req.PostData.Text != nil {
			buf, err := base64.StdEncoding.DecodeString(*req.PostData.Text)
			if err == nil {
				var j map[string]interface{}
				err := json.Unmarshal(buf, &j)
				if err == nil {
					key := j["key"].(string)
					switch value := j["value"].(type) {
					case string:
						return key, value
					default:
						buf, err := json.Marshal(value)
						if err == nil {
							return key, string(buf)
						}
					}
				}
			}
		}
		return "", ""
	}

	key, value := parseKeyValue(req)
	err := kv.mongoClient.SetKey(key, value)
	if err != nil {
		return common.Error(err)
	}
	return common.Text("ok")
}

func (kv *KVApp) setKey(req *nlibshared.Request) *nlibshared.Response {
	if req.Method == "GET" {
		return kv.setKeyGET(req)
	} else if req.Method == "POST" || req.Method == "PUT" {
		return kv.setKeyPOST(req)
	}
	return common.Err405
}

func (kv *KVApp) getRecent(req *nlibshared.Request) *nlibshared.Response {
	skip := common.GetQueryAsInt(req, "skip", 0)
	limit := common.GetQueryAsInt(req, "limit", 10)
	res, err := kv.mongoClient.GetRecent(skip, limit)
	if err != nil {
		return common.Error(err)
	}
	return common.JSON(res)
}
