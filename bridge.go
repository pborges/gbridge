package gbridge

type DeviceHandlerFunc func(dev Device, cmd CommandRequest, res *CommandResponse)

type Bridge struct {
	ClientId     string
	ClientSecret string
	Devices      map[string]Device
	fns          map[string]DeviceHandlerFunc
}
