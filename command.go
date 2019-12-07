package gbridge

import "errors"

type Command interface {
	Name() string
	// passing maps in functions feels dirty
	Execute(ctx Context, args map[string]interface{}) error
}

type OnOffCommand func(ctx Context, state bool) error

func (t OnOffCommand) Execute(ctx Context, args map[string]interface{}) error {
	if val, ok := args["on"]; ok {
		if state, ok := val.(bool); ok {
			return t(ctx, state)
		}
		return errors.New("argument 'on' should be a bool")
	}
	return errors.New("missing argument 'on'")
}

func (t OnOffCommand) Name() string {
	return "action.agents.commands.OnOff"
}
