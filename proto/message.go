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

type CommandStatus string

const (
	CommandStatusSuccess CommandStatus = "SUCCESS"
	CommandStatusError   CommandStatus = "ERROR"
)

type CommandResponse struct {
	Ids       []string        `json:"ids"`
	Status    CommandStatus   `json:"status"`
	States    json.RawMessage `json:"states"`
	ErrorCode ErrorCode       `json:"errorCode,omitempty"`
}

type ErrorCode string

func (e ErrorCode) Error() string {
	return string(e)
}

const (
	ErrorCodeNone            ErrorCode = ""
	ErrorCodeAuthExpired     ErrorCode = "authExpired"
	ErrorCodeAuthFailure     ErrorCode = "authFailure"
	ErrorCodeDeviceOffline   ErrorCode = "deviceOffline"
	ErrorCodeTimeout         ErrorCode = "timeout"
	ErrorCodeDeviceTurnedOff ErrorCode = "deviceTurnedOff"
	ErrorCodeDeviceNotFound  ErrorCode = "deviceNotFound"
	ErrorCodeValueOutofRange ErrorCode = "valueOutOfRange"
	ErrorCodeNotSupported    ErrorCode = "notSupported"
	ErrorCodeProtocolError   ErrorCode = "protocolError"
	ErrorCodeUnknown         ErrorCode = "unknownError"
)
