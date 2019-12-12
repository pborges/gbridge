package proto

type DeviceType string

const (
	DeviceTypeCamera     DeviceType = "action.devices.types.CAMERA"
	DeviceTypeLight      DeviceType = "action.devices.types.LIGHT"
	DeviceTypeOutlet     DeviceType = "action.devices.types.OUTLET"
	DeviceTypeSwitch     DeviceType = "action.devices.types.SWITCH"
	DeviceTypeThermostat DeviceType = "action.devices.types.THERMOSTAT"
	DeviceTypeBlinds 	 DeviceType = "action.devices.types.BLINDS"
)
