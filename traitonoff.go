package gbridge

import (
	"errors"
	"github.com/pborges/gbridge/proto"
)

// OnOffTrait provides an implementation of the Smart Home OnOff Trait Schema from Google Smart Home Actions
type OnOffTrait struct {
	CommandOnlyOnOff bool
	OnExecuteChange  OnOffCommand
	OnStateHandler   func(Context) (bool, proto.ErrorCode)
}

func (t OnOffTrait) ValidateTrait() error {
	if t.OnExecuteChange == nil {
		return errors.New("OnExecuteChange cannot be nil")
	}
	if t.OnStateHandler == nil {
		return errors.New("OnStateHandler cannot be nil")
	}

	return nil
}
func (t OnOffTrait) TraitName() string {
	return "action.devices.traits.OnOff"
}

func (t OnOffTrait) TraitStates(ctx Context) []State {
	var onOffState State
	onOffState.Name = "on"
	onOffState.Value, onOffState.Error = t.OnStateHandler(ctx)

	return []State{onOffState}
}

func (t OnOffTrait) TraitCommands() []Command {
	return []Command{t.OnExecuteChange}
}

func (t OnOffTrait) TraitAttributes() []Attribute {
	return []Attribute{
		{
			Name:  "commandOnlyOnOff",
			Value: t.CommandOnlyOnOff,
		},
	}
}
