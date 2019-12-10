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

func (s *SmartHome) decodeAndHandle(agent agentContext, r io.Reader) proto.IntentMessageResponse {
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
		log.Printf("[%s] Intent: %s-> %s\n", agent.AgentUserId, i.Intent, string(i.Payload))
		switch i.Intent {
		case "action.devices.SYNC":
			res.Payload = s.handleSyncIntent(agent)
		case "action.devices.EXECUTE":
			requestBody := proto.ExecRequest{}
			if err := json.Unmarshal(i.Payload, &requestBody); err == nil {
				res.Payload = s.handleExecuteIntent(agent, requestBody)
			} else {
				res.Payload = proto.ErrorResponse{
					Status:    proto.CommandStatusError,
					ErrorCode: proto.ErrorCodeProtocolError.Error(),
				}
			}
		case "action.devices.QUERY":
			//todo handle query
		}
	}
	return res
}

func (s *SmartHome) Handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		if agent, ok := s.agents[oauth.GetAgentUserIdFromHeader(r)]; ok {
			res := s.decodeAndHandle(agent, r.Body)
			if err := json.NewEncoder(io.MultiWriter(w, log.Writer())).Encode(res); err != nil {
				log.Println("error encoding sync response:", err)
				http.Error(w, "error encoding sync response: "+err.Error(), http.StatusInternalServerError)
			}
		} else {
			http.Error(w, "unable to find agentUserId", http.StatusInternalServerError)
		}
	}
}

func (s *SmartHome) handleExecuteIntent(agent agentContext, req proto.ExecRequest) proto.ExecResponse {
	log.Printf("[%s] EXEC: %+v\n", agent.AgentUserId, req)

	var ids []string
	for _, c := range req.Commands {
		for _, d := range c.Devices {
			ids = append(ids, d.ID)
		}
	}
	responseBody := proto.ExecResponse{}
	// for each command
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
								if err := cmd.Execute(ctx, e.Params);err == nil {
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
	return responseBody
}

func (s *SmartHome) handleSyncIntent(agent agentContext) proto.SyncResponse {
	log.Printf("[%s] SYNC\n", agent.AgentUserId)

	devices := make([]proto.Device, 0, len(agent.Devices))
	for _, d := range agent.Devices {

		devTraits := d.DeviceTraits()
		traits := make([]string, 0, len(devTraits))
		for _, t := range d.DeviceTraits() {
			traits = append(traits, t.TraitName())
		}

		var info proto.DeviceInfo
		if p, ok := d.(DeviceInfoProvider); ok {
			info = p.DeviceInfo()
		}

		var roomHint string
		if p, ok := d.(DeviceRoomHintProvider); ok {
			roomHint = p.DeviceRoomHint()
		}

		dev := proto.Device{
			Id:         d.DeviceId(),
			Type:       d.DeviceType(),
			Traits:     traits,
			Name:       d.DeviceName(),
			DeviceInfo: info,
			RoomHint:   roomHint,

			//todo make interfaces for these types
			CustomData:      make(map[string]interface{}),
			WillReportState: false,
			Attributes:      struct{}{},
		}
		devices = append(devices, dev)
	}
	return proto.SyncResponse{
		AgentUserId: agent.AgentUserId,
		Devices:     devices,
	}
}

func (s *SmartHome) RegisterDevice(agentUserId string, dev Device) error {
	s.lock.RLock()
	defer s.lock.RUnlock()

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

type agentContext struct {
	AgentUserId string
	Devices     map[string]Device
}
