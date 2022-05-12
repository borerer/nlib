package models

type ResponseCheckin struct {
	Version string        `json:"version"`
	Mongo   *ServiceMongo `json:"mongo,omitempty"`
}
