package models

type ServiceMongo struct {
	Uri  string `json:"uri"`
	User string `json:"user"`
	Pass string `json:"pass"`
	Host string `json:"host"`
	Port string `json:"port"`
}
