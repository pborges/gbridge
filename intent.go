package gbridge

const IntentSync = "action.devices.SYNC"
const IntentExecute = "action.devices.EXECUTE"
const IntentQuery = "action.devices.QUERY"

type IntentAspect struct {
	Intent string
	Func   func()
}
