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
	return "action.agents.commands.OnOff"
}
