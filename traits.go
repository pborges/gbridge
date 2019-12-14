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

// OpenCloseTrait provides an implementation of the Smart Home OpenClose Trait Schema from Google Smart Home Actions
type OpenCloseTrait struct {
	DiscreteOnlyOpenClose bool
	OpenDirection []OpenCloseTraitDirection
	QueryOnlyOpenClose bool
	OnExecuteChange  OpenCloseCommand
	OnDirectionStateHandler   func(Context) (string, proto.ErrorCode)
	OnPercentStateHandler func(Context) (float64, proto.ErrorCode)
}

type OpenCloseTraitDirection string 
const OpenCloseTraitDirectionUp OpenCloseTraitDirection = "UP"
const OpenCloseTraitDirectionDown OpenCloseTraitDirection = "DOWN"
const OpenCloseTraitDirectionLeft OpenCloseTraitDirection = "LEFT"
const OpenCloseTraitDirectionRight OpenCloseTraitDirection = "RIGHT"

// ValidateTrait checks if all required attributes and handlers are created/set
func (t OpenCloseTrait) ValidateTrait() error {
	if t.OnExecuteChange == nil {
		return errors.New("OnExecuteChange cannot be nil")
	}
	if t.OnDirectionStateHandler == nil || t.OnPercentStateHandler == nil {
		return errors.New("Both OnStateHandlers must be filled cannot be nil")
	}

	return nil
}

// TraitName returns the string how google defines this traits name
func (t OpenCloseTrait) TraitName() string {
	return "action.devices.traits.OpenClose"
}

// TraitStates parses the diffrent state attributes and calls the corresponding handlers
func (t OpenCloseTrait) TraitStates(ctx Context) []State {
	var openPercentState State
	openPercentState.Name = "openPercent"
	openPercentState.Value, openPercentState.Error = t.OnPercentStateHandler(ctx)

	// optional 
	var openDirectionState State
	openDirectionState.Name = "openDirection"
	openDirectionState.Value, openDirectionState.Error = t.OnDirectionStateHandler(ctx) 

	// check status handler
	return []State{openPercentState, openDirectionState}
}


func (t OpenCloseTrait) TraitCommands() []Command {
	return []Command{t.OnExecuteChange}
}

// TraitAttributes defines all Attributes of the OpenCloseTrait
func (t OpenCloseTrait) TraitAttributes() []Attribute {
	atr := []Attribute{
		{
			Name:  "discreteOnlyOpenClose",
			Value: t.DiscreteOnlyOpenClose,
		},
		{
			Name: "queryOnlyOpenClose",
			Value: t.QueryOnlyOpenClose,
		},
	}
	
	// if optional Argument openDirection is set, add it to arguments. 
	if len(t.OpenDirection) > 0 {
		openDirectionArg := Attribute{
			Name: "openDirection",
			Value: t.OpenDirection,
		}
		atr = append(atr, openDirectionArg)
	}

	return atr
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
