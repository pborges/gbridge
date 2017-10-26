package gbridge

import (
	"log"
	"net/http"
	"encoding/json"
	"os"
	"io"
)

func (b *Bridge) HandleSmartHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.RequestURI)

	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Authorization")

	req := IntentMessageRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("error decoding request:", err)
	}
	for _, i := range req.Inputs {
		log.Printf("Intent: %s-> %s\n", i.Intent, string(i.Payload))
		switch i.Intent {
		case "action.devices.SYNC":
			log.Println("SYNC")
			devices := []Device{}
			for _, v := range b.Devices {
				devices = append(devices, v)
			}
			if err := json.NewEncoder(io.MultiWriter(w, os.Stdout)).Encode(IntentMessageResponse{
				RequestId: req.RequestId,
				Payload: SyncResponse{
					AgentUserId: b.AgentUserId,
					Devices:     devices,
				},
			}); err != nil {
				log.Println("result error:", err)
			}
			return
		case "action.devices.EXECUTE":
			requestBody := ExecRequest{}
			if err := json.Unmarshal(i.Payload, &requestBody); err != nil {
				log.Println(err)
				return
			}
			log.Printf("EXEC: %+v\n", requestBody)
			ids := []string{}
			for _, c := range requestBody.Commands {
				for _, d := range c.Devices {
					ids = append(ids, d.ID)
				}
			}
			responseBody := ExecResponse{}
			for _, c := range requestBody.Commands {
				for _, d := range c.Devices {
					if fn, ok := b.fns[d.ID]; ok {
						for _, e := range c.Execution {
							r := CommandResponse{
								Ids:    ids,
								Status: CommandStatusError,
							}
							fn(b.Devices[d.ID], e, &r)
							if r.Status == CommandStatusError && r.ErrorCode == "" {
								r.ErrorCode = DeviceErrorUnknownError
							}
							responseBody.Commands = append(responseBody.Commands, r)
						}
					} else {
						responseBody.Commands = append(responseBody.Commands, CommandResponse{
							Ids:       ids,
							Status:    CommandStatusError,
							ErrorCode: DeviceErrorDeviceNotFound,
						})
					}
				}
			}

			if err := json.NewEncoder(io.MultiWriter(w, os.Stdout)).Encode(IntentMessageResponse{
				RequestId: req.RequestId,
				Payload:   responseBody,
			}); err != nil {
				log.Println(err)
			}
			return
		}
	}
}
