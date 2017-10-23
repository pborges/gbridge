package gbridge

import (
	"log"
	"net/http"
	"encoding/json"
)

type Token struct {
	TokenType    string `json:"token_type"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (b *Bridge) HandleToken(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.RequestURI)

	clientId := r.FormValue("client_id")
	clientSecret := r.FormValue("client_secret")
	grantType := r.FormValue("grant_type")
	code := r.FormValue("code")

	log.Println("clientId:", clientId)
	log.Println("clientSecret:", clientSecret)
	log.Println("grantType:", grantType)
	log.Println("code:", code)

	if clientId == "" || clientSecret == "" {
		log.Println("missing ID or Secret")
		http.Redirect(w, r, "missing required parameter", http.StatusBadRequest)
		return
	}

	if clientId != b.ClientId && clientSecret != b.ClientSecret {
		log.Println("wrong ID or Secret")
		http.Redirect(w, r, "incorrect client data", http.StatusBadRequest)
		return
	}

	if grantType == "authorization_code" {
		// TODO: use real access and refresh tokens
		// and validate the user blah blah blah
		// seems any static value will work although its not secure
		// https://github.com/actions-on-google/actionssdk-smart-home-nodejs/blob/7620e79f309331e112b1e6a113274aa135957210/smart-home-provider/cloud/auth-provider.js#L295
		t := Token{
			TokenType:    "bearer",
			AccessToken:  "db283094-b74b-11e7-abc4-cec278b6b50a",
			RefreshToken: "e2e649c4-b74b-11e7-abc4-cec278b6b50a",
		}
		log.Printf("authorization_code: %+v\n", t)
		err := json.NewEncoder(w).Encode(t)
		if err != nil {
			log.Println(err)
		}
	} else if grantType == "refresh_token" {
		// TODO: actually fill this in...
	} else {
		log.Println("grant type not supported")
		http.Redirect(w, r, "grant_type "+grantType+"is not supported", http.StatusBadRequest)
		return
	}
}
