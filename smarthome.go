package gbridge

import (
	"encoding/json"
	"errors"
	"github.com/pborges/gbridge/oauth"
	"github.com/pborges/gbridge/proto"
	"io"
	"log"
	"net/http"
	"sync"
)

type IntentAspect struct {
	Intent string
	Func   func()
}

//todo: naive concurrency scheme, this should be looked at and redone

type agentContext struct {
	AgentUserId string
	Devices     map[string]Device
}

type SmartHome struct {
	Log     *log.Logger
	agents  map[string]agentContext
	lock    sync.RWMutex
	Aspects []IntentAspect // these fire when the given intent is fired these should not block
}

func (s *SmartHome) decodeAndHandle(agentUserId string, r io.Reader) proto.IntentMessageResponse {
	req := proto.IntentMessageRequest{}

	if err := json.NewDecoder(r).Decode(&req); err != nil {
		if s.Log != nil {
			s.Log.Println("error decoding request:", err)
		}
	}

	res := proto.IntentMessageResponse{
		RequestId: req.RequestId,
	}

	s.lock.RLock()
	defer s.lock.RUnlock()

	for _, i := range req.Inputs {
		switch i.Intent {
		case "action.devices.SYNC":
			res.Payload = s.handleSyncIntent(agentUserId)
		case "action.devices.EXECUTE":
			requestBody := proto.ExecRequest{}
			if err := json.Unmarshal(i.Payload, &requestBody); err == nil {
				res.Payload = s.handleExecuteIntent(agentUserId, requestBody)
			} else {
				res.Payload = proto.ErrorResponse{
					Status:    proto.CommandStatusError,
					ErrorCode: proto.ErrorCodeProtocolError.Error(),
				}
			}
		case "action.devices.QUERY":
			requestBody := proto.QueryRequest{}
			if err := json.Unmarshal(i.Payload, &requestBody); err == nil {
				res.Payload = s.handleQueryIntent(requestBody)
			} else {
				res.Payload = proto.ErrorResponse{
					Status:    proto.CommandStatusError,
					ErrorCode: proto.ErrorCodeProtocolError.Error(),
				}
			}
		}

		for _, a := range s.Aspects {
			if a.Intent == i.Intent {
				a.Func()
			}
		}
	}
	return res
}

func (s *SmartHome) executeCommandForResponse(dev Device, ex proto.CommandRequest) proto.CommandResponse {

	// check all traits...
	for _, trait := range dev.DeviceTraits() {
		// for all the commands...
		for _, cmd := range trait.TraitCommands() {
			// for the right command
			if ex.Command == cmd.Name() {
				// and execute
				ctx := Context{Target: dev}
				res := cmd.Execute(ctx, ex.Params)
				proto.SetIds(&res, dev.DeviceId())
				return res
			}
		}
	}
	return proto.CommandResponse{
		ErrorCode: proto.ErrorCodeNotSupported,
	}
}

func (s *SmartHome) handleExecuteIntent(agentUserId string, req proto.ExecRequest) proto.ExecResponse {
	if s.Log != nil {
		o, _ := json.Marshal(req)
		s.Log.Printf("[%s] EXEC REQUEST: %s\n", agentUserId, string(o))
	}

	var ids []string
	for _, c := range req.Commands {
		for _, d := range c.Devices {
			ids = append(ids, d.ID)
		}
	}

	var resCount int
	var responseBody proto.ExecResponse
	responses := make(chan proto.CommandResponse)
	// for the correct agent
	if agent, ok := s.agents[agentUserId]; ok {
		// for each command
		for _, c := range req.Commands {
			// for each device
			for _, d := range c.Devices {
				if devCtx, ok := agent.Devices[d.ID]; ok {
					// execute all the things
					for _, e := range c.Execution {
						resCount++
						go func(s *SmartHome, ch chan proto.CommandResponse, dev Device, ex proto.CommandRequest) {
							ch <- s.executeCommandForResponse(dev, ex)
						}(s, responses, devCtx, e)
					}
				} else {
					responseBody.Commands = append(responseBody.Commands, proto.CommandResponse{
						ErrorCode: proto.ErrorCodeDeviceNotFound,
					})
				}
			}
		}
	}

	for len(responseBody.Commands) < resCount {
		responseBody.Commands = append(responseBody.Commands, <-responses)
	}

	if s.Log != nil {
		o, _ := json.Marshal(responseBody)
		s.Log.Printf("[%s] EXEC RESPONSE: %s\n", agentUserId, string(o))
	}

	return responseBody
}

func (s *SmartHome) encodeDeviceForSyncResponse(dev Device) proto.Device {
	devTraits := dev.DeviceTraits()
	traits := make([]string, 0, len(devTraits))
	attributes := make(map[string]interface{})

	for _, t := range dev.DeviceTraits() {
		traits = append(traits, t.TraitName())
		for _, a := range t.TraitAttributes() {
			attributes[a.Name] = a.Value
		}
	}

	var info proto.DeviceInfo
	if p, ok := dev.(DeviceInfoProvider); ok {
		info = p.DeviceInfo()
	}

	var roomHint string
	if p, ok := dev.(DeviceRoomHintProvider); ok {
		roomHint = p.DeviceRoomHint()
	}

	d := proto.Device{
		Id:         dev.DeviceId(),
		Type:       dev.DeviceType(),
		Traits:     traits,
		Name:       dev.DeviceName(),
		DeviceInfo: info,
		RoomHint:   roomHint,
		Attributes: attributes,

		//todo make interfaces for these types
		CustomData:      make(map[string]interface{}),
		WillReportState: false,
	}

	return d
}

func (s *SmartHome) handleSyncIntent(agentUserId string) proto.SyncResponse {
	if s.Log != nil {
		s.Log.Printf("[%s] SYNC REQUEST\n", agentUserId)
	}

	response := proto.SyncResponse{
		AgentUserId: agentUserId,
		Devices:     make([]proto.Device, 0),
	}

	if agent, ok := s.agents[agentUserId]; ok {
		for _, d := range agent.Devices {
			response.Devices = append(response.Devices, s.encodeDeviceForSyncResponse(d))
		}
	}

	if s.Log != nil {
		o, _ := json.Marshal(response)
		s.Log.Printf("[%s] SYNC RESPONSE: %s\n", agentUserId, string(o))
	}
	return response
}

func (s *SmartHome) handleQueryIntent(request proto.QueryRequest) proto.QueryResponse {
	if s.Log != nil {
		o, _ := json.Marshal(request)
		s.Log.Printf("QUERY REQUEST %s\n", string(o))
	}
	res := proto.QueryResponse{
		Devices: make(map[string]map[string]interface{}),
	}
	for _, r := range request.Devices {
		// look through all the agents for the requested device
		for _, a := range s.agents {
			if d, ok := a.Devices[r.ID]; ok {
				if _, ok := res.Devices[r.ID]; !ok {
					res.Devices[r.ID] = make(map[string]interface{})
				}
				ctx := Context{Target: d}
				for _, t := range d.DeviceTraits() {
					for _, s := range t.TraitStates(ctx) {
						res.Devices[d.DeviceId()][s.Name] = s.Value
						res.Devices[d.DeviceId()]["online"] = s.Error == nil
						if _, ok := res.Devices[d.DeviceId()]["status"]; ok &&
							res.Devices[d.DeviceId()]["status"] == proto.CommandStatusError {
							if s.Error != nil {
								res.Devices[d.DeviceId()]["status"] = proto.CommandStatusError
								res.Devices[d.DeviceId()]["errorCode"] = s.Error
							}
						}
					}
				}
			}
		}
	}
	if s.Log != nil {
		o, _ := json.Marshal(res)
		s.Log.Printf("QUERY RESPONSE %s\n", string(o))
	}
	return res
}

func (s *SmartHome) Handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		res := s.decodeAndHandle(oauth.GetAgentUserIdFromHeader(r), r.Body)
		if err := json.NewEncoder(w).Encode(res); err != nil {
			if s.Log != nil {
				s.Log.Println("error encoding response:", err)
			}
			http.Error(w, "error encoding response: "+err.Error(), http.StatusInternalServerError)
		}
	}
}

func (s *SmartHome) RegisterDevice(agentUserId string, dev Device) error {
	s.lock.RLock()
	defer s.lock.RUnlock()

	reducedTraits := make(map[string]Trait)

	// validate the device and its traits
	for _, t := range dev.DeviceTraits() {
		if _, ok := reducedTraits[t.TraitName()]; ok {
			return errors.New("duplicate trait found: " + t.TraitName())
		}
		if err := t.ValidateTrait(); err != nil {
			return err
		}
	}

	if s.agents == nil {
		s.agents = make(map[string]agentContext)
	}

	if _, ok := s.agents[agentUserId]; !ok {
		s.agents[agentUserId] = agentContext{
			AgentUserId: agentUserId,
			Devices:     make(map[string]Device),
		}
	}

	if _, ok := s.agents[agentUserId].Devices[dev.DeviceId()]; !ok {
		s.agents[agentUserId].Devices[dev.DeviceId()] = dev
	} else {
		return errors.New("device by this name already exists for this agent")
	}

	return nil
}

func (s *SmartHome) UnregisterDevice(agentUserId string, dev string) error {
	s.lock.RLock()
	defer s.lock.RUnlock()

	if _, ok := s.agents[agentUserId]; ok {
		if _, ok := s.agents[agentUserId].Devices[dev]; ok {
			delete(s.agents[agentUserId].Devices, dev)
		} else {
			return errors.New("unknown device")
		}
	} else {
		return errors.New("unknown agent")
	}
	return nil
}

func (s *SmartHome) ListDevices(agentUserId string) []Device {
	s.lock.RLock()
	defer s.lock.RUnlock()

	devs := make([]Device, 0)
	for _, dev := range s.agents[agentUserId].Devices {
		devs = append(devs, dev)
	}
	return devs
}
