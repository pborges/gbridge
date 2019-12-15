package gbridge

import (
	"github.com/pborges/gbridge/proto"
	"github.com/pborges/gbridge/traits"
)

// Device represents a Smart Home Device with a Name, ID, and other metadata as defined by google [https://developers.google.com/assistant/smarthome/concepts/devices-traits]
type Device interface {
	DeviceId() string
	DeviceName() proto.DeviceName
	DeviceType() proto.DeviceType
	DeviceTraits() []traits.Trait
}

// DeviceInfoProvider returns the metadata about a device.
type DeviceInfoProvider interface {
	DeviceInfo() proto.DeviceInfo
}

type DeviceRoomHintProvider interface {
	DeviceRoomHint() string
}

// BasicDevice represents a simple Smart Home Device
type BasicDevice struct {
	Id       string
	Name     proto.DeviceName
	Type     proto.DeviceType
	Traits   []traits.Trait
	Info     proto.DeviceInfo
	RoomHint string
}

// DeviceRoomHint returns in which Room this device should probally be.
func (d BasicDevice) DeviceRoomHint() string {
	return d.RoomHint
}

// DeviceType returns which kind of device this is.
func (d BasicDevice) DeviceType() proto.DeviceType {
	return d.Type
}

// DeviceInfo returns Information about the device
func (d BasicDevice) DeviceInfo() proto.DeviceInfo {
	return d.Info
}

// DeviceId returns the Identifier of the device
func (d BasicDevice) DeviceId() string {
	return d.Id
}

// DeviceName returns the Name
func (d BasicDevice) DeviceName() proto.DeviceName {
	return d.Name
}

// DeviceTraits returns what functionality (trait) the device provides
func (d BasicDevice) DeviceTraits() []traits.Trait {
	return d.Traits
}
