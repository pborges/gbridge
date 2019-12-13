package proto

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

type DeviceName struct {
	Name         string   `json:"name"`
	DefaultNames []string `json:"defaultNames"`
	Nicknames    []string `json:"nicknames"`
}

type DeviceInfo struct {
	Manufacturer string `json:"manufacturer"`
	Model        string `json:"model"`
	HwVersion    string `json:"hwVersion"`
	SwVersion    string `json:"swVersion"`
}
