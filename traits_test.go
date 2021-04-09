package gbridge

import (
	"github.com/pborges/gbridge/proto"
	"testing"
)

func TestDuplicateTraits(t *testing.T) {
	dev := BasicDevice{
		Id:   "",
		Name: proto.DeviceName{},
		Type: "",
		Traits: []Trait{
			OnOffTrait{
				CommandOnlyOnOff: false,
				OnExecuteChange:  nil,
				OnStateHandler:   nil,
			},
			OnOffTrait{
				CommandOnlyOnOff: false,
				OnExecuteChange:  nil,
				OnStateHandler:   nil,
			},
		},
		Info:     proto.DeviceInfo{},
		RoomHint: "",
	}

	home := SmartHome{}

	if err := home.RegisterDevice("test", dev); err == nil {
		t.Error("You cannot register two traits with the same name")
	}
}
