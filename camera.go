package gbridge

import (
	"github.com/pborges/gbridge/proto"
)

type Camera struct {
	ID                 string
	Name               string
	StreamURL          string
	NeedsAuthToken     bool
	NeedsDRMEncryption bool
	Protocol           CameraStreamProtocol
}

func (c Camera) DeviceTraits() []Trait {
	return []Trait{
		CameraStreamTrait{
			SupportedProtocols: []CameraStreamProtocol{c.Protocol},
			OnGetStream: func(ctx Context) (GetCameraStreamResponse, proto.DeviceError) {
				return GetCameraStreamResponse{
					StreamURL: c.StreamURL,
					Protocol:  c.Protocol,
				}, nil
			},
		},
	}
}

func (c Camera) DeviceId() string {
	return c.ID
}

func (c Camera) DeviceName() proto.DeviceName {
	return proto.DeviceName{
		Name: c.Name,
	}
}

func (c Camera) DeviceType() proto.DeviceType {
	return proto.DeviceTypeCamera
}
