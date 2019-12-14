package oauth

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type AuthenticationProvider interface {
	GenerateAuthCode(clientId string, agentUserId string) (string, error)
	GenerateToken(clientId string, clientSecret string, authCode string) (Token, error)
	GetAgentUserIDForToken(token string) (agentUserId string, refresh bool, err error)
}
