package gbridge

import (
	"encoding/json"
	"errors"
	"github.com/pborges/gbridge/oauth"
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

		log.Println("smarthome", r.Header.Get(oauth.AgentUserIdHeader))

		req := IntentMessageRequest{}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Println("error decoding request:", err)
		}

		s.lock.RLock()
		defer s.lock.RUnlock()

		for _, i := range req.Inputs {
			log.Printf("Intent: %s-> %s\n", i.Intent, string(i.Payload))

			switch i.Intent {
			case "action.agents.SYNC":
				log.Println("SYNC")

				//	var devices []Device
				//	for _, d := range s.agents {
				//		devices = append(devices, d)
				//	}
				//	if err := json.NewEncoder(io.MultiWriter(w, os.Stdout)).Encode(IntentMessageResponse{
				//		RequestId: req.RequestId,
				//		Payload: SyncResponse{
				//			AgentUserId: b.AgentUserId,
				//			Devices:     devices,
				//		},
				//	}); err != nil {
				//		log.Println("result error:", err)
				//	}
				//	return
			}
		}
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

	if _, ok := s.agents[agentUserId].Devices[dev.DeviceName()]; ok {
		return errors.New("device by this name already exists for this agent")
	}

	return nil
}

type agentContext struct {
	AgentUserId string
	Devices     map[string]Device
}
