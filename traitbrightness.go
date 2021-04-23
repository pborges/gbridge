package gbridge

import (
	"errors"
	"github.com/pborges/gbridge/proto"
)

// BrightnessTrait provides an implementation of the Smart Home BrightnessTrait Trait Schema from Google Smart Home Actions
type BrightnessTrait struct {
	CommandOnlyBrightness bool
	OnBrightnessChange    BrightnessAbsoluteCommand
	OnStateHandler        func(Context) (int, proto.ErrorCode)
}

func (t BrightnessTrait) ValidateTrait() error {
	if t.OnBrightnessChange == nil {
		return errors.New("OnBrightnessChange cannot be nil")
	}
	if t.OnStateHandler == nil {
		return errors.New("OnStateHandler cannot be nil")
	}

	return nil
}
func (t BrightnessTrait) TraitName() string {
	return "action.devices.traits.Brightness"
}

func (t BrightnessTrait) TraitStates(ctx Context) []State {
	var onOffState State
	onOffState.Name = "brightness"
	onOffState.Value, onOffState.Error = t.OnStateHandler(ctx)

	return []State{onOffState}
}

func (t BrightnessTrait) TraitCommands() []Command {
	return []Command{t.OnBrightnessChange}
}

func (t BrightnessTrait) TraitAttributes() []Attribute {
	return []Attribute{
		{
			Name:  "commandOnlyBrightness",
			Value: t.CommandOnlyBrightness,
		},
	}
}
