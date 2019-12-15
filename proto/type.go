package proto

type DeviceType string

// A list of possible device types can be found here 
// https://developers.google.com/assistant/smarthome/guides
const (
	DeviceTypeCamera     DeviceType = "action.devices.types.CAMERA"
	DeviceTypeLight      DeviceType = "action.devices.types.LIGHT"
	DeviceTypeFan      DeviceType = "action.devices.types.FAN"
	DeviceTypeOutlet     DeviceType = "action.devices.types.OUTLET"
	DeviceTypeSwitch     DeviceType = "action.devices.types.SWITCH"
	DeviceTypeThermostat DeviceType = "action.devices.types.THERMOSTAT"
	DeviceTypeBlinds 	 DeviceType = "action.devices.types.BLINDS"
)
