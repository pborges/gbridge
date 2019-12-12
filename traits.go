package gbridge

import (
	"errors"
	"github.com/pborges/gbridge/proto"
)

type Trait interface {
	TraitName() string
	TraitCommands() []Command
	TraitStates(Context) []State
	TraitAttributes() []Attribute
	ValidateTrait() error
}

// Provided second implementation for blinds
type OpenCloseTrait struct {
	DiscreteOnlyOpenClose bool
	// openDirection []string	// not implementing in first concept
	QueryOnlyOpenClose bool
	OnExecuteChange  OpenCloseCommand
	OnStateHandler   func(Context) (float64, string, proto.ErrorCode, proto.ErrorCode)
}

func (t OpenCloseTrait) ValidateTrait() error {
	if t.OnExecuteChange == nil {
		return errors.New("OnExecuteChange cannot be nil")
	}
	if t.OnStateHandler == nil {
		return errors.New("OnStateHandler cannot be nil")
	}

	return nil
}

func (t OpenCloseTrait) TraitName() string {
	return "action.devices.traits.OpenClose"
}

func (t OpenCloseTrait) TraitStates(ctx Context) []State {
	var openPercentState State
	openPercentState.Name = "openPercent"

	// optional 
	var openDirectionState State
	openDirectionState.Name = "openDirection"

	// check status handler
	openPercentState.Value, openDirectionState.Value, openPercentState.Error, openDirectionState.Error = t.OnStateHandler(ctx)
	return []State{openPercentState,openDirectionState}
}

func (t OpenCloseTrait) TraitCommands() []Command {
	return []Command{t.OnExecuteChange}
}

// Provided Impl, but users SHOULD be able to make their own Traits easy enough by copypasta
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
