package models

type ReqWSCallFunction struct {
	ID      string      `json:"id"`
	Func    string      `json:"func"`
	Payload interface{} `json:"payload"`
}
