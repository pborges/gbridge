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

func (s *SmartHome) handleSyncIntent(agent agentContext) proto.SyncResponse {
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

		dev := proto.Device{
			Id:         d.DeviceId(),
			Type:       d.DeviceType(),
			Traits:     traits,
			Name:       d.DeviceName(),
			DeviceInfo: info,

			//todo make interfaces for these types
			CustomData:      make(map[string]interface{}),
			WillReportState: false,
			Attributes:      struct{}{},
			RoomHint:        "",
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
