package gbridge

type Device interface {
	DeviceName() string
	DeviceTraits() []Trait
}

type BasicDevice struct {
	Name   string
	Traits []Trait
	Info   DeviceInfo
}

func (d BasicDevice) DeviceInfo() DeviceInfo {
	return d.Info
}

func (d BasicDevice) DeviceName() string {
	return d.Name
}
func (d BasicDevice) DeviceTraits() []Trait {
	return d.Traits
}
