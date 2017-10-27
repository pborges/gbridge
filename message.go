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
	Command string `json:"command"`
	Params struct {
		On         bool `json:"on"`
		Brightness int  `json:"brightness"`
	} `json:"params"`
}

type CommandStatus string

const (
	CommandStatusSuccess CommandStatus = "SUCCESS"
	CommandStatusError   CommandStatus = "ERROR"
)

type DeviceState struct {
	On         bool `json:"on"`
	Online     bool `json:"online"`
	Brightness int  `json:"brightness"`
}

type CommandResponse struct {
	Ids       []string      `json:"ids"`
	Status    CommandStatus `json:"status"`
	States    DeviceState   `json:"states"`
	ErrorCode DeviceError   `json:"errorCode"`
}

type QueryRequest struct {
	Devices []struct {
		ID string `json:"id"`
	} `json:"devices"`
}

type QueryResponse struct {
	Devices map[string]DeviceState `json:"devices"`
}
