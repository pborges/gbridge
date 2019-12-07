package gbridge

import "encoding/json"

type IntentMessageRequest struct {
	Inputs []struct {
		Intent  string          `json:"intent"`
		Payload json.RawMessage `json:"payload"`
	} `json:"inputs"`
	RequestId string `json:"requestId"`
}

type IntentMessageResponse struct {
	RequestId string      `json:"requestId"`
	Payload   interface{} `json:"payload"`
}
