package gbridge

import "github.com/pborges/gbridge/proto"

type Device interface {
	DeviceId() string
	DeviceName() proto.DeviceName
	DeviceType() proto.DeviceType
	DeviceTraits() []Trait
}

type DeviceInfoProvider interface {
	DeviceInfo() proto.DeviceInfo
}

type BasicDevice struct {
	Id     string
	Name   proto.DeviceName
	Type   proto.DeviceType
	Traits []Trait
	Info   proto.DeviceInfo
}

func (d BasicDevice) DeviceType() proto.DeviceType {
	return d.Type
}

func (d BasicDevice) DeviceInfo() proto.DeviceInfo {
	return d.Info
}

func (d BasicDevice) DeviceId() string {
	return d.Id
}

func (d BasicDevice) DeviceName() proto.DeviceName {
	return d.Name
}

func (d BasicDevice) DeviceTraits() []Trait {
	return d.Traits
}
