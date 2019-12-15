package traits

import "github.com/pborges/gbridge"

type Trait interface {
	TraitName() string
	TraitCommands() []Command
	TraitStates(gbridge.Context) []gbridge.State
	TraitAttributes() []gbridge.Attribute
	ValidateTrait() error
}
