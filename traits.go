package gbridge

type Trait interface {
	TraitName() string
	TraitCommands() []Command
	TraitStates(Context) []State
	TraitAttributes() []Attribute
	ValidateTrait() error
}
