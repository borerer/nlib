package kv

import (
	"encoding/base64"
	"encoding/json"
	"errors"

	nlibshared "github.com/borerer/nlib-shared/go"
	"github.com/borerer/nlib/app/builtin/kv/database"
	"github.com/borerer/nlib/app/common"
	"github.com/borerer/nlib/configs"
)

type KVApp struct {
	config      *configs.KVConfig
	mongoClient *database.MongoClient
}

func NewKVApp(config *configs.KVConfig) *KVApp {
	return &KVApp{
		config: config,
	}
}

func (kv *KVApp) Start() error {
	kv.mongoClient = database.NewMongoClient(&kv.config.Mongo)
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
	if name == "get" {
		return kv.getKey(req)
	} else if name == "set" {
		return kv.setKey(req)
	} else {
		return common.Err404
	}
}

func (kv *KVApp) getKey(req *nlibshared.Request) *nlibshared.Response {
	key := common.GetQuery(req, "key")
	val, err := kv.mongoClient.GetKey(key)
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
