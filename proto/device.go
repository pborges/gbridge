package proto

// Device represets a Google Smart Home Device
type Device struct {
	Id              string                 `json:"id"`
	Type            DeviceType             `json:"type"`
	Traits          []string               `json:"traits"`
	Name            DeviceName             `json:"name"`
	WillReportState bool                   `json:"willReportState"`
	Attributes      map[string]interface{} `json:"attributes,omitempty"`
	RoomHint        string                 `json:"roomHint,omitempty"`
	DeviceInfo      DeviceInfo             `json:"deviceInfo,omitempty"`
	CustomData      map[string]interface{} `json:"customData,omitempty"`
}

// DeviceName contains the names and nicknames of a device.
type DeviceName struct {
	Name         string   `json:"name"`
	DefaultNames []string `json:"defaultNames,omitempty"`
	Nicknames    []string `json:"nicknames,omitempty"`
}

// DeviceInfo holds the metadata (manufacturer,model,versions) of a Device
type DeviceInfo struct {
	Manufacturer string `json:"manufacturer,omitempty"`
	Model        string `json:"model,omitempty"`
	HwVersion    string `json:"hwVersion,omitempty"`
	SwVersion    string `json:"swVersion,omitempty"`
}
