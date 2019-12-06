package oauth

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Server struct {
	AuthProvider AuthenticationProvider
}

func (s Server) HandleAuth() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		clientId := r.URL.Query().Get("client_id")
		redirectUrl := r.URL.Query().Get("redirect_uri")
		state := r.URL.Query().Get("state")
		responseType := r.URL.Query().Get("response_type")

		if responseType != "code" {
			http.Error(w, "response_type must be code, got "+responseType, http.StatusInternalServerError)
			return
		}

		if authCode, err := s.AuthProvider.GenerateAuthCodeForClient(clientId); err == nil {
			http.Redirect(w, r, fmt.Sprintf("%s?code=%s&state=%s", redirectUrl, authCode, state), http.StatusFound)
		} else {
			http.Error(w, "unable to generate an authentication code, "+err.Error(), http.StatusBadRequest)
		}
	}
}

func (s Server) HandleToken() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		if r.Method == http.MethodPost {
			if err := r.ParseForm(); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			clientId := r.Form.Get("client_id")
			clientSecret := r.Form.Get("client_secret")
			grantType := r.Form.Get("grant_type")
			code := r.Form.Get("code")

			if grantType != "authorization_code" {
				http.Error(w, "grant_type must be authorization_code, got "+grantType, http.StatusInternalServerError)
				return
			}

			if token, err := s.AuthProvider.ValidateAuthCodeAndGenerateToken(clientId, clientSecret, code); err == nil {
				t := struct {
					TokenType string `json:"token_type"`
					Token
				}{
					TokenType: "bearer",
					Token:     token,
				}

				err := json.NewEncoder(w).Encode(t)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			} else {
				http.Error(w, "unable to validate client, "+err.Error(), http.StatusBadRequest)
			}
		}
	}
}
