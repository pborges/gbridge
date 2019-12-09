package proto

import (
	"encoding/json"
	"errors"
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

type CommandStatus string

const (
	CommandStatusSuccess CommandStatus = "SUCCESS"
	CommandStatusError   CommandStatus = "ERROR"
)

type CommandResponse struct {
	Ids       []string        `json:"ids"`
	States    json.RawMessage `json:"states"`
	Status    CommandStatus   `json:"status"`
	ErrorCode ErrorCode       `json:"errorCode,omitempty"`
}

type ErrorCode error

var (
	ErrorCodeAuthExpired     ErrorCode = errors.New("authExpired")
	ErrorCodeAuthFailure     ErrorCode = errors.New("authFailure")
	ErrorCodeDeviceOffline   ErrorCode = errors.New("deviceOffline")
	ErrorCodeTimeout         ErrorCode = errors.New("timeout")
	ErrorCodeDeviceTurnedOff ErrorCode = errors.New("deviceTurnedOff")
	ErrorCodeDeviceNotFound  ErrorCode = errors.New("deviceNotFound")
	ErrorCodeValueOutofRange ErrorCode = errors.New("valueOutOfRange")
	ErrorCodeNotSupported    ErrorCode = errors.New("notSupported")
	ErrorCodeProtocolError   ErrorCode = errors.New("protocolError")
	ErrorCodeUnknown         ErrorCode = errors.New("unknownError")
)
