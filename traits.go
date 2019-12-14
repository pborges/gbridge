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
	OpenDirection         []OpenCloseTraitDirection
	QueryOnlyOpenClose    bool
	OnExecuteChange       OpenCloseCommand
	OnStateHandler        func(Context) ([]OpenState, proto.ErrorCode)
}

// OpenCloseTraitDirection represents the different directions a device can open as a capitalized string defined by google
type OpenCloseTraitDirection string

const OpenCloseTraitDirectionNone OpenCloseTraitDirection = ""
const OpenCloseTraitDirectionUp OpenCloseTraitDirection = "UP"
const OpenCloseTraitDirectionDown OpenCloseTraitDirection = "DOWN"
const OpenCloseTraitDirectionLeft OpenCloseTraitDirection = "LEFT"
const OpenCloseTraitDirectionRight OpenCloseTraitDirection = "RIGHT"

// ValidateTrait checks if all required attributes and handlers are created/set
func (t OpenCloseTrait) ValidateTrait() error {
	if t.OnExecuteChange == nil {
		return errors.New("OnExecuteChange cannot be nil")
	}
	if t.OnStateHandler == nil {
		return errors.New("OnStateHandlers cannot be nil")
	}
	return nil
}

// TraitName returns the string how google defines this traits name
func (t OpenCloseTrait) TraitName() string {
	return "action.devices.traits.OpenClose"
}

type OpenState struct {
	OpenPercent   float64
	OpenDirection OpenCloseTraitDirection
}

// TraitStates parses the different state attributes and calls the corresponding handlers
func (t OpenCloseTrait) TraitStates(ctx Context) []State {
	onOffState := State{
		Name:  "on",
		Value: true,
	}

	handlerOpenState, err := t.OnStateHandler(ctx)

	// return current state
	curOpenState := State{
		Name:  "openState",
		Value: handlerOpenState,
		Error: err,
	}

	// check status handler
	return []State{onOffState, curOpenState}
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
			Name:  "queryOnlyOpenClose",
			Value: t.QueryOnlyOpenClose,
		},
	}

	// if optional Argument openDirection is set, add it to arguments.
	if len(t.OpenDirection) > 0 {
		openDirectionArg := Attribute{
			Name:  "openDirection",
			Value: t.OpenDirection,
		}
		atr = append(atr, openDirectionArg)
	}

	return atr
}

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
