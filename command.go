package gbridge

import (
	"github.com/pborges/gbridge/proto"
)

type Command interface {
	Name() string
	// passing maps in functions feels dirty
	Execute(ctx Context, args map[string]interface{}) proto.DeviceError
}

type OnOffCommand func(ctx Context, state bool) proto.DeviceError

func (t OnOffCommand) Execute(ctx Context, args map[string]interface{}) proto.DeviceError {
	if val, ok := args["on"]; ok {
		if state, ok := val.(bool); ok {
			return t(ctx, state)
		}
		return proto.ErrorCodeNotSupported
	}
	return proto.ErrorCodeProtocolError
}

func (t OnOffCommand) Name() string {
	return "action.devices.commands.OnOff"
}

// OpenCloseCommand defines how a function should handle this specific trait
type OpenCloseCommand func(ctx Context, params interface{}) proto.DeviceError

// Execute validates the request and calls the user defined handler for the device trait
func (t OpenCloseCommand) Execute(ctx Context, args map[string]interface{}) proto.DeviceError {
	// validate if our EXECUTE Request contains the "openPercent" object so the handler can actually handle something
	if val, ok := args["openPercent"]; ok {
		return t(ctx, val) 
	}
	return proto.ErrorCodeProtocolError
}

// Name returns the intent for this Trait
func (t OpenCloseCommand) Name() string {
	return "action.devices.commands.OpenClose"
}
