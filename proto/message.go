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

type QueryRequest struct {
	Devices []struct {
		ID string `json:"id"`
	} `json:"devices"`
}

type QueryResponse struct {
	Devices map[string]map[string]interface{} `json:"devices"`
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
	ErrorCode string        `json:"errorCode,omitempty"`
}

func SetIds(r *CommandResponse, ids ...string) {
	if r != nil {
		r.ids = ids
	}
}

type CommandResponse struct {
	ids       []string
	Results   map[string]interface{}
	ErrorCode ErrorCode
}

func (r CommandResponse) MarshalJSON() ([]byte, error) {
	res := make(map[string]interface{})
	res["status"] = CommandStatusSuccess
	res["ids"] = r.ids
	if r.ErrorCode != nil {
		res["errorCode"] = r.ErrorCode.Error()
		res["status"] = CommandStatusError
	}
	if r.Results != nil {
		for k, v := range r.Results {
			res[k] = v
		}
	}

	return json.Marshal(res)
}
