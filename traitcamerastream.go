package gbridge

import (
	"errors"
	"github.com/pborges/gbridge/proto"
)

type CameraStreamProtocol string

const CameraStreamProtocolHLS = CameraStreamProtocol("hls")
const CameraStreamProtocolDash = CameraStreamProtocol("dash")
const CameraStreamProtocolSmoothStream = CameraStreamProtocol("smooth_stream")
const CameraStreamProtocolProgressiveMP4 = CameraStreamProtocol("progressive_mp4")

type GetCameraStreamResponse struct {
	StreamURL string
	AuthToken string
	Protocol  CameraStreamProtocol
}

type GetCameraStreamCommand func(ctx Context) (GetCameraStreamResponse, proto.DeviceError)

// Execute checks if the arguments from the Intent Request are correct and passes them to a user-defined type safe execute handler
func (t GetCameraStreamCommand) Execute(ctx Context, args map[string]interface{}) proto.CommandResponse {
	res := proto.CommandResponse{
		Results:   make(map[string]interface{}),
		ErrorCode: proto.ErrorCodeProtocolError,
	}
	var r GetCameraStreamResponse
	r, res.ErrorCode = t(ctx)
	if res.ErrorCode == nil {
		res.Results["cameraStreamAccessUrl"] = r.StreamURL
		res.Results["cameraStreamProtocol"] = r.Protocol
		if r.AuthToken != "" {
			res.Results["cameraStreamAuthToken"] = r.AuthToken
		}
	}
	return res
}

// Name returns the Identifier String of the Command
func (t GetCameraStreamCommand) Name() string {
	return "action.devices.commands.GetCameraStream"
}

// https://developers.google.com/assistant/smarthome/traits/camerastream
type CameraStreamTrait struct {
	OnGetStream        GetCameraStreamCommand
	SupportedProtocols []CameraStreamProtocol
	NeedsAuthToken     bool
	NeedsDRMEncryption bool
}

func (t CameraStreamTrait) TraitName() string {
	return "action.devices.traits.CameraStream"
}

func (t CameraStreamTrait) TraitCommands() []Command {
	return []Command{t.OnGetStream}
}

func (t CameraStreamTrait) TraitStates(Context) []State {
	return []State{}
}

func (t CameraStreamTrait) TraitAttributes() []Attribute {
	return []Attribute{
		{
			Name:  "cameraStreamSupportedProtocols",
			Value: t.SupportedProtocols,
		},
		{
			Name:  "cameraStreamNeedAuthToken",
			Value: t.NeedsAuthToken,
		},
		{
			Name:  "cameraStreamNeedDrmEncryption",
			Value: t.NeedsDRMEncryption,
		},
	}
}
func (t CameraStreamTrait) ValidateTrait() error {
	if len(t.SupportedProtocols) <= 0 {
		return errors.New("SupportedProtocols cannot be empty")
	}
	return nil
}
