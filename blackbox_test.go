package gbridge

import (
	"encoding/json"
	"github.com/pborges/gbridge/proto"
	"bytes"
	"strings"
	"testing"
)

type BasicLight struct {
	BasicDevice
	State bool
}

func (l BasicLight) Type() proto.DeviceType {
	return proto.DeviceTypeLight
}

func (l *BasicLight) DeviceTraits() []Trait {
	return []Trait{
		OnOffTrait{
			CommandOnlyOnOff: false,
			OnExecuteChange: func(ctx Context, state bool) proto.DeviceError {
				if l.Id == "456" {
					return proto.DeviceErrorTurnedOff
				}
				l.State = state
				return nil
			},
			OnStateHandler: func(ctx Context) (b bool, code proto.ErrorCode) {
				return l.State, nil
			},
		},
	}
}

func TestMultipleExecution(t *testing.T) {
	home := SmartHome{}
	home.RegisterDevice("test", &BasicLight{
		BasicDevice: BasicDevice{
			Id: "123",
			Name: proto.DeviceName{
				Name: "Light1",
			},
		},
	})

	home.RegisterDevice("test", &BasicLight{
		BasicDevice: BasicDevice{
			Id: "456",
			Name: proto.DeviceName{
				Name: "Light2",
			},
		},
	})

	res := home.decodeAndHandle(home.agents["test"],
		strings.NewReader(`{
					  "requestId": "ff36a3cc-ec34-11e6-b1a0-64510650abcf",
					  "inputs": [{
						"intent": "action.devices.EXECUTE",
						"payload": {
						  "commands": [{
							"devices": [{
							  "id": "123",
							  "customData": {
								"fooValue": 74,
								"barValue": true,
								"bazValue": "sheepdip"
							  }
							}, {
								  "id": "456",
								  "customData": {
									"fooValue": 36,
									"barValue": false,
									"bazValue": "moarsheep"
								  }
							}],
							"execution": [{
							  "command": "action.devices.commands.OnOff",
							  "params": {
								"on": true
							  }
							}]
						  }]
						}
					  }]
					}`))
	buf := bytes.NewBufferString("")
	json.NewEncoder(buf).Encode(res)

	if strings.TrimSpace(buf.String()) != `{"requestId":"ff36a3cc-ec34-11e6-b1a0-64510650abcf","payload":{"commands":[{"ids":["123","456"],"states":{"on":true,"online":true},"status":"SUCCESS"},{"ids":["123","456"],"states":{"online":true},"status":"ERROR","errorCode":"turnedOff"}]}}`{
		t.Error("unexpected response")
	}
}
