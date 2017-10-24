package gbridge

type DeviceType string

const (
	DeviceTypeCamera     DeviceType = "action.devices.types.CAMERA"
	DeviceTypeLight      DeviceType = "action.devices.types.LIGHT"
	DeviceTypeOutlet     DeviceType = "action.devices.types.OUTLET"
	DeviceTypeSwitch     DeviceType = "action.devices.types.SWITCH"
	DeviceTypeThermostat DeviceType = "action.devices.types.THERMOSTAT"
)

type DeviceTrait string

const (
	DeviceTraitBrightness       DeviceTrait = "action.devices.traits.Brightness"
	DeviceTraitCameraStream     DeviceTrait = "action.devices.traits.CameraStream"
	DeviceTraitColorSpectrum    DeviceTrait = "action.devices.traits.ColorSpectrum"
	DeviceTraitColorTemperature DeviceTrait = "action.devices.traits.ColorTemperature"
	DeviceTraitOnOff            DeviceTrait = "action.devices.traits.OnOff"
	DeviceTraitStartStop        DeviceTrait = "action.devices.traits.StartStop"
	DeviceTemperatureSettings   DeviceTrait = "action.devices.traits.TemperatureSettings"
	DeviceTraitToggles          DeviceTrait = "action.devices.traits.Toggles"
)

type DeviceError string

const (
	DeviceErrorAuthExpired     DeviceError = "authExpired"
	DeviceErrorAuthFailure     DeviceError = "authFailure"
	DeviceErrorDeviceOffline   DeviceError = "deviceOffline"
	DeviceErrorTimeout         DeviceError = "timeout"
	DeviceErrorDeviceTurnedOff DeviceError = "deviceTurnedOff"
	DeviceErrorDeviceNotFound  DeviceError = "deviceNotFound"
	DeviceErrorValueOutofRange DeviceError = "valueOutOfRange"
	DeviceErrorNotSupported    DeviceError = "notSupported"
	DeviceErrorProtocolError   DeviceError = "protocolError"
	DeviceErrorUnknownError    DeviceError = "unknownError"
)

type DeviceName struct {
	DefaultNames []string `json:"defaultNames"`
	Name         string   `json:"name"`
	Nicknames    []string `json:"nicknames"`
}

type DeviceInfo struct {
	Manufacturer string `json:"manufacturer"`
	Model        string `json:"model"`
	HwVersion    string `json:"hwVersion"`
	SwVersion    string `json:"swVersion"`
}

type Device struct {
	Id              string        `json:"id"`
	Type            DeviceType    `json:"type"`
	Traits          []DeviceTrait `json:"traits"`
	Name            DeviceName    `json:"name"`
	WillReportState bool          `json:"willReportState"`
	Attributes struct {
	} `json:"attributes,omitempty"`
	RoomHint   string                 `json:"roomHint,omitempty"`
	DeviceInfo *DeviceInfo             `json:"deviceInfo,omitempty"`
	CustomData map[string]interface{} `json:"customData,omitempty"`
}

func (b *Bridge) HandleDevice(d Device, fn DeviceHandlerFunc) {
	if b.Devices == nil {
		b.Devices = make(map[string]Device)
	}
	b.Devices[d.Id] = d
	if b.fns == nil {
		b.fns = make(map[string]DeviceHandlerFunc)
	}
	b.fns[d.Id] = fn
}
