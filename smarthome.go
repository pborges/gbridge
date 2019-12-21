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

//todo: naive concurrency scheme, this should be looked at and redone

type SmartHome struct {
	agents map[string]agentContext
	lock   sync.RWMutex
}

func (s *SmartHome) decodeAndHandle(agentUserId string, r io.Reader) proto.IntentMessageResponse {
	req := proto.IntentMessageRequest{}

	if err := json.NewDecoder(r).Decode(&req); err != nil {
		log.Println("error decoding request:", err)
	}

	res := proto.IntentMessageResponse{
		RequestId: req.RequestId,
	}

	s.lock.RLock()
	defer s.lock.RUnlock()

	for _, i := range req.Inputs {
		log.Printf("[%s] Intent: %s-> %s\n", agentUserId, i.Intent, string(i.Payload))
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
	}
	return res
}

func (s *SmartHome) Handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		res := s.decodeAndHandle(oauth.GetAgentUserIdFromHeader(r), r.Body)
		if err := json.NewEncoder(io.MultiWriter(w, log.Writer())).Encode(res); err != nil {
			log.Println("error encoding response:", err)
			http.Error(w, "error encoding response: "+err.Error(), http.StatusInternalServerError)
		}
	}
}

func (s *SmartHome) handleExecuteIntent(agentUserId string, req proto.ExecRequest) proto.ExecResponse {
	log.Printf("[%s] EXEC: %+v\n", agentUserId, req)

	var ids []string
	for _, c := range req.Commands {
		for _, d := range c.Devices {
			ids = append(ids, d.ID)
		}
	}
	responseBody := proto.ExecResponse{}
	// for each command
	if agent, ok := s.agents[agentUserId]; ok {
		for _, c := range req.Commands {
			// for each device
			for _, d := range c.Devices {
				if devCtx, ok := agent.Devices[d.ID]; ok {
					// execute all the things
					for _, e := range c.Execution {
						r := proto.CommandResponse{
							Ids:       ids,
							ErrorCode: proto.ErrorCodeNotSupported.Error(),
							States:    make(map[string]interface{}),
						}
						r.States["online"] = true

						// check all traits...
						for _, trait := range devCtx.DeviceTraits() {
							// for all the commands...
							for _, cmd := range trait.TraitCommands() {
								// for the right command
								if e.Command == cmd.Name() {
									// and execute
									ctx := Context{Target: devCtx}
									if err := cmd.Execute(ctx, e.Params); err == nil {
										r.ErrorCode = ""
										for _, s := range trait.TraitStates(ctx) {
											if s.Error == nil {
												r.Status = proto.CommandStatusSuccess
												r.States[s.Name] = s.Value
											} else {
												r.Status = proto.CommandStatusError
												r.ErrorCode = s.Error.Error()
												r.States["online"] = false
												break
											}
										}
									} else {
										r.Status = proto.CommandStatusError
										r.ErrorCode = err.Error()
									}
								}
							}
						}
						responseBody.Commands = append(responseBody.Commands, r)
					}
				}
			}
		}
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
	log.Printf("[%s] SYNC\n", agentUserId)

	devices := make([]proto.Device, 0)
	if agent, ok := s.agents[agentUserId]; ok {
		for _, d := range agent.Devices {
			devices = append(devices, s.encodeDeviceForSyncResponse(d))
		}
	}

	return proto.SyncResponse{
		AgentUserId: agentUserId,
		Devices:     devices,
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

func (s *SmartHome) handleQueryIntent(request proto.QueryRequest) proto.QueryResponse {
	log.Printf("QUERY %+v\n", request)
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

	return res
}

type agentContext struct {
	AgentUserId string
	Devices     map[string]Device
}
