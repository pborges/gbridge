package gbridge

import (
	fmt "fmt"
	"log"
	"net/http"
)

func (b *Bridge) HandleOauth(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.RequestURI)
	clientId := r.URL.Query().Get("client_id")
	redirectUrl := r.URL.Query().Get("redirect_uri")
	state := r.URL.Query().Get("state")
	responseType := r.URL.Query().Get("response_type")
	authCode := r.URL.Query().Get("code")

	log.Println("clientId:", clientId)
	log.Println("redirectUrl:", redirectUrl)
	log.Println("state:", state)
	log.Println("responseType:", responseType)

	if responseType != "code" {
		http.Error(w, "response_type must be code, got "+responseType, http.StatusInternalServerError)
		return
	}

	if authCode != "" {
		log.Println("Authcode successful:", authCode)
		http.Redirect(w, r, fmt.Sprintf("%s?code=%s&state=%s", redirectUrl, authCode, state), http.StatusFound)
		return
	}

	// TODO: check to see if theres a user in session otherwise redirect to login and generate an authCode
	if authCode != "" {
		log.Println("Authcode successful:", authCode)
		http.Redirect(w, r, fmt.Sprintf("%s?code=%s&state=%s", redirectUrl, authCode, state), http.StatusFound)
		return
	}

	http.Error(w, "something went wrong", http.StatusBadRequest)
}
