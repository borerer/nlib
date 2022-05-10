package utils

import (
	"encoding/json"
	"fmt"
)

func ToJSON(v interface{}) []byte {
	buf, err := json.Marshal(v)
	if err != nil {
		fmt.Printf("marshal json error: %v\n", err)
		return nil
	}
	return buf
}

func ToJSONString(v interface{}) string {
	return string(ToJSON(v))
}
