package models

type RegisterResponse struct {
	Version string        `json:"version"`
	Mongo   *ServiceMongo `json:"mongo,omitempty"`
}
