package gbridge

type Device interface {
	DeviceName() string
	DeviceTraits() []Trait
}

type BaseDevice struct {
	Name   string
	Traits []Trait
	Info   DeviceInfo
}

func (d BaseDevice) DeviceInfo() DeviceInfo {
	return d.Info
}

func (d BaseDevice) DeviceName() string {
	return d.Name
}
func (d BaseDevice) DeviceTraits() []Trait {
	return d.Traits
}
