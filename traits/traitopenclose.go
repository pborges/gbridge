package traits

import (
	"errors"
	"github.com/pborges/gbridge"
	"github.com/pborges/gbridge/proto"
)

// MultiDirectionOpenCloseTrait provides an implementation of the Smart Home OpenClose Trait Schema from Google Smart Home Actions
type MultiDirectionOpenCloseTrait struct {
	DiscreteOnlyOpenClose bool                      // Defaults to false. When set to true, this indicates that the device must either be fully open or fully closed (that is, it does not support values between 0% and 100%). An example of such a device may be a valve.
	OpenDirection         []OpenCloseTraitDirection //Array of strings. Optional. Required if the device supports opening and closing in more than one direction. Comma-separated list of directions in which this device can be opened. Valid options include: UP, DOWN, LEFT, RIGHT, IN, and OUT. For example, top-down bottom-up blinds may open either up or down.
	QueryOnlyOpenClose    bool                      // Indicates if the device can only be queried for state information and cannot be controlled. Sensors that can only report open state should set this field to true.
	OnExecuteChange       DirectionalOpenCloseCommand
	OnStateHandler        func(gbridge.Context) ([]OpenState, proto.ErrorCode)
}

// OpenCloseTraitDirection represents the different directions a device can open as a capitalized string defined by google
type OpenCloseTraitDirection string

const OpenCloseTraitDirectionNone OpenCloseTraitDirection = ""
const OpenCloseTraitDirectionUp OpenCloseTraitDirection = "UP"
const OpenCloseTraitDirectionDown OpenCloseTraitDirection = "DOWN"
const OpenCloseTraitDirectionLeft OpenCloseTraitDirection = "LEFT"
const OpenCloseTraitDirectionRight OpenCloseTraitDirection = "RIGHT"

// ValidateTrait checks if all required attributes and handlers are created/set
func (t MultiDirectionOpenCloseTrait) ValidateTrait() error {
	if t.OnExecuteChange == nil {
		return errors.New("OnExecuteChange cannot be nil")
	}
	if t.OnStateHandler == nil {
		return errors.New("OnStateHandlers cannot be nil")
	}
	return nil
}

// TraitName returns the string how google defines this traits name
func (t MultiDirectionOpenCloseTrait) TraitName() string {
	return "action.devices.traits.OpenClose"
}

type OpenState struct {
	OpenPercent   float64
	OpenDirection OpenCloseTraitDirection
}

// TraitStates parses the different state attributes and calls the corresponding handlers
func (t MultiDirectionOpenCloseTrait) TraitStates(ctx gbridge.Context) []gbridge.State {
	onOffState := gbridge.State{
		Name:  "on",
		Value: true,
	}

	handlerOpenState, err := t.OnStateHandler(ctx)

	// return current state
	curOpenState := gbridge.State{
		Name:  "openState",
		Value: handlerOpenState,
		Error: err,
	}

	// check status handler
	return []gbridge.State{onOffState, curOpenState}
}

func (t MultiDirectionOpenCloseTrait) TraitCommands() []Command {
	return []Command{t.OnExecuteChange}
}

// TraitAttributes defines all Attributes of the MultiDirectionOpenCloseTrait
func (t MultiDirectionOpenCloseTrait) TraitAttributes() []gbridge.Attribute {
	atr := []gbridge.Attribute{
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
		openDirectionArg := gbridge.Attribute{
			Name:  "openDirection",
			Value: t.OpenDirection,
		}
		atr = append(atr, openDirectionArg)
	}

	return atr
}

// MultiDirectionOpenCloseTrait provides an implementation of the Smart Home OpenClose Trait Schema from Google Smart Home Actions
type OpenCloseTrait struct {
	DiscreteOnlyOpenClose bool //Defaults to false. When set to true, this indicates that the device must either be fully open or fully closed (that is, it does not support values between 0% and 100%). An example of such a device may be a valve.
	QueryOnlyOpenClose    bool
	OnExecuteChange       OpenCloseCommand
	OnStateHandler        func(gbridge.Context) (float64, proto.ErrorCode)
}

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

// TraitStates parses the different state attributes and calls the corresponding handlers
func (t OpenCloseTrait) TraitStates(ctx gbridge.Context) []gbridge.State {
	onOffState := gbridge.State{
		Name:  "on",
		Value: true,
	}

	handlerOpenState, err := t.OnStateHandler(ctx)

	// return current state
	curOpenState := gbridge.State{
		Name:  "openState",
		Value: handlerOpenState,
		Error: err,
	}

	// check status handler
	return []gbridge.State{onOffState, curOpenState}
}

func (t OpenCloseTrait) TraitCommands() []Command {
	return []Command{t.OnExecuteChange}
}

// TraitAttributes defines all Attributes of the OpenCloseTrait
func (t OpenCloseTrait) TraitAttributes() []gbridge.Attribute {
	atr := []gbridge.Attribute{
		{
			Name:  "discreteOnlyOpenClose",
			Value: t.DiscreteOnlyOpenClose,
		},
		{
			Name:  "queryOnlyOpenClose",
			Value: t.QueryOnlyOpenClose,
		},
	}

	return atr
}
