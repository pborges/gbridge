package proto

import (
	"encoding/json"
)

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

type SyncResponse struct {
	AgentUserId string   `json:"agentUserId"`
	Devices     []Device `json:"devices"`
}

type ExecRequest struct {
	Commands []struct {
		Devices []struct {
			ID string `json:"id"`
		} `json:"devices"`
		Execution []CommandRequest `json:"execution"`
	} `json:"commands"`
}

type ExecResponse struct {
	Commands []CommandResponse `json:"commands"`
}

type CommandRequest struct {
	Command string                 `json:"command"`
	Params  map[string]interface{} `json:"params"`
}

type ErrorResponse struct {
	Status    CommandStatus `json:"status"`
	ErrorCode ErrorCode     `json:"errorCode,omitempty"`
}

type CommandResponse struct {
	Ids       []string        `json:"ids"`
	States    json.RawMessage `json:"states"`
	Status    CommandStatus   `json:"status"`
	ErrorCode ErrorCode       `json:"errorCode,omitempty"`
}
