package models

type GeneralResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

var GeneralOK = &GeneralResponse{
	Status: "ok",
}

func GeneralError(message string) *GeneralResponse {
	return &GeneralResponse{
		Status:  "error",
		Message: message,
	}
}
