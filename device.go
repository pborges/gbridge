package gbridge

import "github.com/pborges/gbridge/proto"

type Device interface {
	DeviceId() string
	DeviceName() proto.DeviceName
	DeviceType() proto.DeviceType
	DeviceTraits() []Trait
	DeviceAttributes() []Attribute
}

type DeviceInfoProvider interface {
	DeviceInfo() proto.DeviceInfo
}

type DeviceRoomHintProvider interface {
	DeviceRoomHint() string
}

type BasicDevice struct {
	Id         string
	Name       proto.DeviceName
	Type       proto.DeviceType
	Traits     []Trait
	Attributes []Attribute
	Info       proto.DeviceInfo
	RoomHint   string
}

func (d BasicDevice) DeviceRoomHint() string {
	return d.RoomHint
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

func (d BasicDevice) DeviceAttributes() []Attribute {
	return d.Attributes
}
