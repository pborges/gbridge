package main

import (
	"net/http"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/pborges/gbridge"
	"log"
	"os"
)

var addr = ":8085"

func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	r := mux.NewRouter()

	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("UNKKNOWN:", r.RequestURI)
	})

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("/")
		fmt.Fprint(w, "Hello")
	})

	b := gbridge.Bridge{
		ClientId:     "123456", //as long as this matches the settings "Account linking" on actions console it works
		ClientSecret: "654321", //as long as this matches the settings "Account linking" on actions console it works
	}

	b.HandleDevice(NewSwitch("1", "alarm"), func(dev gbridge.Device, req gbridge.CommandRequest, res *gbridge.CommandResponse) {
		log.Printf("Exec Cmd: %+v\n", req)
		res.Status = gbridge.CommandStatusSuccess
		res.States.Online = true
		res.States.On = req.Params.On
		log.Printf("Exec Res: %+v\n", res)
	})

	r.HandleFunc("/oauth", b.HandleOauth)
	r.HandleFunc("/token", b.HandleToken)
	r.HandleFunc("/smarthome", b.HandleSmartHome)
	log.Println("Listening:", addr)
	log.Println(http.ListenAndServe(addr, r))
}

func NewSwitch(id string, name string) gbridge.Device {
	d := gbridge.Device{
		Id:     id,
		Type:   gbridge.DeviceTypeSwitch,
		Traits: []gbridge.DeviceTrait{gbridge.DeviceTraitOnOff},
		Name: gbridge.DeviceName{
			DefaultNames: []string{name},
			Name:         name,
			Nicknames:    []string{name},
		},
	}
	return d
}
