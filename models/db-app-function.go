package models

type DBAppFunction struct {
	AppID string `json:"app_id" bson:"app_id"`
	Func  string `json:"func" bson:"func"`
}
