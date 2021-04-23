package gbridge

import (
	"github.com/pborges/gbridge/proto"
)

// BrightnessAbsoluteCommand is the basic set brightness functionality for any device that has integer brightness (https://developers.google.com/assistant/smarthome/traits/brightness.html)
type BrightnessAbsoluteCommand func(ctx Context, value int) proto.DeviceError

// Execute checks if the arguments from the Intent Request are correct and passes them to a user-defined type safe execute handler
func (t BrightnessAbsoluteCommand) Execute(ctx Context, args map[string]interface{}) proto.CommandResponse {
	res := proto.CommandResponse{
		Results: map[string]interface{}{},
	}
	if val, ok := args["brightness"]; ok {
		if state, ok := val.(float64); ok {
			res.ErrorCode = t(ctx, int(state))
		} else {
			res.ErrorCode = proto.ErrorCodeNotSupported
		}
	}
	return res
}

// Name returns the Identifier String of the Command
func (t BrightnessAbsoluteCommand) Name() string {
	return "action.devices.commands.BrightnessAbsolute"
}
