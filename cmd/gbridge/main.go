package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/pborges/gbridge"
	"github.com/pborges/gbridge/oauth"
	"github.com/pborges/gbridge/proto"
)

var loginPage = `
<!DOCTYPE html>
<html>
	<head>
		<meta name="viewport" content="width=device-width, initial-scale=1">
		<style>
			* {
				box-sizing: border-box;
			}
			input[type=text], input[type=password], select, textarea {
				width: 100%;
				padding: 12px;
				border: 1px solid #ccc;
				border-radius: 4px;
				resize: vertical;
			}
			label {
				padding: 12px 12px 12px 0;
				display: inline-block;
			}
			input[type=submit] {
				background-color: #4CAF50;
				color: white;
				padding: 12px 20px;
				border: none;
				border-radius: 4px;
				cursor: pointer;
				float: right;
			}
			input[type=submit]:hover {
				background-color: #45a049;
			}
			.container {
				border-radius: 5px;
				background-color: #f2f2f2;
				padding: 20px;
			}
			.col-25 {
				float: left;
				width: 25%;
				margin-top: 6px;
			}
			.col-75 {
				float: left;
				width: 75%;
				margin-top: 6px;
			}
			.row:after {
				content: "";
				display: table;
				clear: both;
			}
			@media screen and (max-width: 600px) {
				.col-25, .col-75, input[type=submit] {
					width: 100%;
					margin-top: 0;
				}
			}
		</style>
	</head>
	<body>
		<center>
			<h1>GBridge</h1>
		</center>
		<div class="container">
		<form method="POST">
			<div class="row">
				<div class="col-25">
					<label for="agentUserId">AgentUserId</label> 
				</div>
				<div class="col-75">
					<input type="text" id="agentUserId" name="agentUserId">
				</div>
			</div>
			<div class="row">
				<div class="col-25">
					<label for="password">Password</label> 
				</div>
				<div class="col-75">
					<input type="password" id="password" name="password">
				</div>
			</div>
			<div class="row">
				&nbsp;
			</div>
			<div class="row">
				<input type="submit" value="Submit">
			</div>
		</form>
		</div>
	</body>
</html>
`

func logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("HTTP", r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}

func main() {
	log.Println("Server started")
	mux := http.NewServeMux()
	authProvider := oauth.SimpleAuthenticationProvider{}
	smartHome := gbridge.SmartHome{}

	// configure the oauth server
	oauthServer := oauth.Server{
		AuthenticationProvider: &authProvider,
		AgentUserLoginHandler: func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet {
				log.Println("display login page")
				fmt.Fprint(w, loginPage)
			} else if r.Method == http.MethodPost {
				agentUserId := r.FormValue("agentUserId")
				password := r.FormValue("password")
				//todo: validate agentUserId and password...
				if agentUserId != "" && password != "" {
					log.Println("register agent", agentUserId)
					authProvider.RegisterAgent(agentUserId)

					//This Handler must set up the agentUserId Header or the oauth cannot continue
					oauth.SetAgentUserIdHeader(r, agentUserId)
				} else {
					http.Error(w, "invalid agentUserId", http.StatusInternalServerError)
				}
			}
		},
	}

	// register clients
	authProvider.RegisterClient("123456", "654321")

	// register devices
	if err := smartHome.RegisterDevice("pborges", gbridge.BasicDevice{
		Id: "1234567890",
		Name: proto.DeviceName{
			Name: "Light1",
		},
		Type: proto.DeviceTypeLight,
		Traits: []gbridge.Trait{
			gbridge.OnOffTrait{
				CommandOnlyOnOff: false,
				OnExecuteChange: func(ctx gbridge.Context, state bool) proto.DeviceError {
					log.Println("turn", ctx.Target.DeviceName(), "device", state)
					return nil
				},
				OnStateHandler: func(ctx gbridge.Context) (bool, proto.ErrorCode) {
					log.Println("query state of", ctx.Target.DeviceName())
					return false, nil
				},
			}},
		Info: proto.DeviceInfo{
			HwVersion: "1.0",
		},
		}); err != nil {
		log.Fatal(err)
	}

	// register devices
	if err := smartHome.RegisterDevice("pborges", gbridge.BasicDevice{
		Id: "1234567891",
		Name: proto.DeviceName{
			Name: "Blind1",
		},
		Type: proto.DeviceTypeBlinds,
		Traits: []gbridge.Trait{
			gbridge.OpenCloseTrait{
				DiscreteOnlyOpenClose: true,
				OpenDirection: []gbridge.OpenCloseTraitDirection{gbridge.OpenCloseTraitDirectionUp, gbridge.OpenCloseTraitDirectionDown},
				QueryOnlyOpenClose: false,
				OnExecuteChange: func(ctx gbridge.Context, params interface{}) proto.DeviceError {
					log.Println("Percent of", ctx.Target.DeviceName(), "should be set to", params)
					return nil
				},
				OnStateHandler: func(ctx gbridge.Context) (gbridge.OpenState, proto.ErrorCode) {
					log.Println("query state of", ctx.Target.DeviceName())
					curOpenState := gbridge.OpenState{100.0,gbridge.OpenCloseTraitDirectionUp, nil}
					return curOpenState, nil
				},
			}},
		Info: proto.DeviceInfo{
			HwVersion: "1.0",
		},
	}); err != nil {
		log.Fatal(err)
	}

	// set up the http endpoints
	mux.HandleFunc("/oauth", oauthServer.HandleAuth())
	mux.HandleFunc("/token", oauthServer.HandleToken())
	mux.HandleFunc("/smarthome", oauthServer.Authenticate(smartHome.Handle()))

	// serve!
	log.Fatal(http.ListenAndServe(":8085", logger(mux)))
}
