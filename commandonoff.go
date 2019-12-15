package gbridge

import (
	"github.com/pborges/gbridge/proto"
)

// OnOffCommand is the basic on and off functionality for any device that has binary on and off (https://developers.google.com/assistant/smarthome/traits/onoff.html)
type OnOffCommand func(ctx Context, state bool) proto.DeviceError

// Execute checks if the arguments from the Intent Request are correct and passes them to a user-defined type safe execute handler
func (t OnOffCommand) Execute(ctx Context, args map[string]interface{}) proto.DeviceError {
	if val, ok := args["on"]; ok {
		if state, ok := val.(bool); ok {
			return t(ctx, state)
		}
		return proto.ErrorCodeNotSupported
	}
	return proto.ErrorCodeProtocolError
}

// Name returns the Identifier String of the Command
func (t OnOffCommand) Name() string {
	return "action.devices.commands.OnOff"
}
