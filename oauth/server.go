package oauth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const AgentUserIdHeader = "agentUserId"

func SetAgentUserIdHeader(r *http.Request, agentUserId string) {
	r.Header.Set(AgentUserIdHeader, agentUserId)
}

func GetAgentUserIdFromHeader(r *http.Request) string {
	return r.Header.Get(AgentUserIdHeader)
}

type Server struct {
	AuthenticationProvider AuthenticationProvider
	AgentUserLoginHandler  http.HandlerFunc
}

func (s Server) Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accessToken := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		if agentUserId, _, err := s.AuthenticationProvider.GetAgentUserIdForToken(accessToken); err == nil {
			SetAgentUserIdHeader(r, agentUserId)
			next(w, r)
		} else {
			http.Error(w, err.Error(), http.StatusForbidden)
		}
	}
}

func (s Server) HandleAuth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		// let the caller display the login page, and handle the post back before we take over
		s.AgentUserLoginHandler(w, r)

		if r.Method == http.MethodPost {
			// see if we have a agent set up in the headers
			if agentUserId := GetAgentUserIdFromHeader(r); agentUserId != "" {
				clientId := r.URL.Query().Get("client_id")
				redirectUrl := r.URL.Query().Get("redirect_uri")
				state := r.URL.Query().Get("state")
				responseType := r.URL.Query().Get("response_type")

				if responseType != "code" {
					http.Error(w, "response_type must be code, got "+responseType, http.StatusInternalServerError)
					return
				}

				if authCode, err := s.AuthenticationProvider.GenerateAuthCode(clientId, agentUserId); err == nil {
					url := fmt.Sprintf("%s?code=%s&state=%s", redirectUrl, authCode, state)
					http.Redirect(w, r, url, http.StatusFound)
				} else {
					http.Error(w, "unable to generate an authentication code, "+err.Error(), http.StatusBadRequest)
				}
			}
		}
	}
}

func (s Server) HandleToken() http.HandlerFunc {
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
			authCode := r.Form.Get("code")

			if grantType != "authorization_code" {
				http.Error(w, "grant_type must be authorization_code, got "+grantType, http.StatusInternalServerError)
				return
			}

			if token, err := s.AuthenticationProvider.GenerateToken(clientId, clientSecret, authCode); err == nil {
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
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	}
}
