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

func (s *SmartHome) Handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		if agent, ok := s.agents[oauth.GetAgentUserIdFromHeader(r)]; ok {
			req := proto.IntentMessageRequest{}

			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
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
						http.Error(w, err.Error(), http.StatusInternalServerError)
					}
				}
			}

			if err := json.NewEncoder(io.MultiWriter(w, log.Writer())).Encode(res); err != nil {
				log.Println("error encoding sync response:", err)
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
			if ctx, ok := agent.Devices[d.ID]; ok {
				// execute all the things
				for _, e := range c.Execution {
					r := proto.CommandResponse{
						Ids:       ids,
						Status:    proto.CommandStatusError,
						ErrorCode: proto.ErrorCodeNotSupported,
					}
					// check all traits...
					for _, trait := range ctx.DeviceTraits() {
						// for all the commands...
						for _, cmd := range trait.TraitCommands() {
							// for the right command
							if e.Command == cmd.Name() {
								// and execute
								r.Status, r.ErrorCode = cmd.Execute(Context{Target: ctx}, e.Params)
							}
						}
					}
					responseBody.Commands = append(responseBody.Commands, r)
				}
			} else {
				responseBody.Commands = append(responseBody.Commands, proto.CommandResponse{
					Ids:       ids,
					Status:    proto.CommandStatusError,
					ErrorCode: proto.ErrorCodeDeviceNotFound,
				})
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
