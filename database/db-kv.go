package database

type DBKV struct {
	AppID   string `json:"app_id" bson:"app_id"`
	Key     string `json:"key" bson:"key"`
	Value   string `json:"value" bson:"value"`
	Created int64  `json:"created" bson:"created"`
	Updated int64  `json:"updated" bson:"updated"`
}
