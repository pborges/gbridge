package main

import (
	"fmt"
	"github.com/pborges/gbridge"
	"github.com/pborges/gbridge/oauth"
	"github.com/pborges/gbridge/proto"
	"log"
	"net/http"
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
			}"
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
		//for name, headers := range r.Header {
		//	for _, h := range headers {
		//		log.Println("\tHEADER:", name, "->", h)
		//	}
		//}
		next.ServeHTTP(w, r)
	})
}

func main() {
	mux := http.NewServeMux()

	authProvider := oauth.SimpleAuthenticationProvider{}
	authProvider.RegisterClient("123456", "654321")

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

	smartHome := gbridge.SmartHome{}

	if err := smartHome.RegisterDevice("pborges", gbridge.BasicDevice{
		Id: "1234567890",
		Name: proto.DeviceName{
			Name: "Light1",
		},
		Type: proto.DeviceTypeLight,
		Traits: []gbridge.Trait{
			gbridge.OnOffTrait{
				CommandOnlyOnOff: false,
				OnExecuteChange: func(ctx gbridge.Context, state bool) (proto.CommandStatus, proto.ErrorCode) {
					log.Println("turn", ctx.Target.DeviceName(), "device", state)
					return proto.CommandStatusSuccess, proto.ErrorCodeNone
				},
				OnStateHandler: func(ctx gbridge.Context) (bool, error) {
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

	mux.HandleFunc("/oauth", oauthServer.HandleAuth())
	mux.HandleFunc("/token", oauthServer.HandleToken())
	mux.HandleFunc("/smarthome", oauthServer.Authenticate(smartHome.Handle()))

	log.Fatal(http.ListenAndServe(":8085", logger(mux)))
}