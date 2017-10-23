# G Bridge
## library for Actions on Google Home Control
**WARNING** work in progress

porting from [actionssdk-smart-home-nodejs](https://github.com/actions-on-google/actionssdk-smart-home-nodejs)

### Useful links
* [https://developers.google.com/actions/smarthome/](https://developers.google.com/actions/smarthome/)
* [https://github.com/actions-on-google/actionssdk-smart-home-nodejs](https://github.com/actions-on-google/actionssdk-smart-home-nodejs)
* [https://console.actions.google.com/](https://console.actions.google.com/)

### Instructions
1. Create an action.json file
```
    {
      "actions": [{
        "name": "actions.devices",
        "deviceControl": {
        },
        "fulfillment": {
          "conversationName": "automation"
        }
      }],
      "conversations": {
        "automation" :
        {
          "name": "automation",
          "url": "https://<YOUR URL>/smarthome"
        }
      }
    }
```

2. Create a project on [https://console.actions.google.com/](https://console.actions.google.com/)

3. Click **Use Actions SDK**

4. Use the [gActions](https://developers.google.com/actions/tools/gactions-cli) CLI to run the command given with the 'action.json' file as your Action Package.

5. Click **Okay**.

6. Click **ADD** under **App information**.

7. Give your App some information like an invocation name, some description, and some policy and contact info.

8. Click **Save**.

9. Click **Add** under **Account Linking**.

10. Select **Authorization Code** for Grant Type.

11. Under **Client Information**, enter the client ID and secret from earlier.

12. The Authorization URL is the hosted URL of your app with '/oauth' as the path, e.g. https:/<YOUR URL>/oauth

13. The Token URL is the hosted URL of your app with '/token' as the path, e.g. https://<YOUR URL>/token

14. Click **TEST DRAFT**

15. Open the **Google Home App** on your phone and navigate to the **Home Control** section

15. Click the 3 dots and go to **Manage Accounts**

16. Select your project from the **Add New** section


### Example
```
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


```
