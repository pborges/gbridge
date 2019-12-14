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
type OpenCloseCommand func(ctx Context, openPercent float64, openDirection OpenCloseTraitDirection) proto.DeviceError

// Execute validates the request and calls the user defined handler for the device trait
func (t OpenCloseCommand) Execute(ctx Context, args map[string]interface{}) proto.DeviceError {
	// validate if our EXECUTE Request contains the "openPercent" object so the handler can actually handle something
	openDirection := OpenCloseTraitDirectionNone
	if argOpenPercent, ok := args["openPercent"]; ok {
		if openPercent, ok := argOpenPercent.(float64); ok {
			// if openDirection is specified, set it, otherwise lets use none
			if argDir, ok := args["openDirection"]; ok {
				if dir, ok := argDir.(string); ok {
					openDirection = OpenCloseTraitDirection(dir)
				}
			}
			return t(ctx, openPercent, openDirection)
		}
		return proto.ErrorCodeNotSupported
	}
	return proto.ErrorCodeProtocolError
}

// Name returns the intent for this Trait
func (t OpenCloseCommand) Name() string {
	return "action.devices.commands.OpenClose"
}
