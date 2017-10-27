package gbridge

type ExecHandlerFunc func(dev Device, cmd CommandRequest, res *CommandResponse)
type QueryHandlerFunc func(dev Device, res *DeviceState)

type Bridge struct {
	AgentUserId  string
	ClientId     string
	ClientSecret string
	Devices      map[string]*DeviceContext
}

type DeviceContext struct {
	Device Device
	Exec   ExecHandlerFunc
	Query  QueryHandlerFunc
}
