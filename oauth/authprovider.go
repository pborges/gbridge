package oauth

import (
	"errors"
	"sync"
)

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type AuthenticationProvider interface {
	GenerateAuthCode(clientId string, agentUserId string) (string, error)
	GenerateToken(clientId string, clientSecret string, authCode string) (Token, error)
	GetAgentUserIdForToken(token string) (agentUserId string, refresh bool, err error)
}

// a very simple authentication provider

type SimpleAuthenticationProvider struct {
	clients map[string]simpleAuthProviderClient // key: clientId
	agents  map[string]simpleAuthProviderAgent  // key: agentUserId
	lock    sync.Mutex
}

type simpleAuthProviderClient struct {
	ID        string
	Secret    string
	AuthCodes map[string]string
}

type simpleAuthProviderAgent struct {
	ID    string
	Token Token
}

func (m *SimpleAuthenticationProvider) RegisterClient(clientId string, clientSecret string) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if m.clients == nil {
		m.clients = make(map[string]simpleAuthProviderClient)
	}
	m.clients[clientId] = simpleAuthProviderClient{
		ID:        clientId,
		Secret:    clientSecret,
		AuthCodes: make(map[string]string),
	}
}

func (m *SimpleAuthenticationProvider) RegisterAgent(agentUserId string) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if m.agents == nil {
		m.agents = make(map[string]simpleAuthProviderAgent)
	}
	agent := simpleAuthProviderAgent{
		ID: agentUserId,
		Token: Token{
			AccessToken:  GenerateRandomString(36),
			RefreshToken: GenerateRandomString(36),
		},
	}
	m.agents[agent.ID] = agent
}

func (m *SimpleAuthenticationProvider) GenerateAuthCode(clientId string, agentUserId string) (string, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	// do we have a valid client?
	if client, ok := m.clients[clientId]; ok {
		// do we have a valid agent?
		if agent, ok := m.agents[agentUserId]; ok {
			authCode := GenerateRandomString(36)
			client.AuthCodes[authCode] = agent.ID
			return authCode, nil
		} else {
			return "", errors.New("unknown agent")
		}
	}

	return "", errors.New("unknown client")
}

func (m *SimpleAuthenticationProvider) GenerateToken(clientId string, clientSecret string, authCode string) (Token, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	// do we have a valid client?
	if client, ok := m.clients[clientId]; ok {
		// with a valid secret
		if client.Secret != clientSecret {
			return Token{}, errors.New("invalid client secret")
		}

		// do we have a valid agent?
		if agent, ok := client.AuthCodes[authCode]; ok {
			return m.agents[agent].Token, nil
		} else {
			return Token{}, errors.New("unknown agent")
		}
	}

	return Token{}, errors.New("unknown client")
}

func (m *SimpleAuthenticationProvider) GetAgentUserIDForToken(token string) (agentUserId string, refresh bool, err error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	for _, a := range m.agents {
		if a.Token.AccessToken == token {
			return a.ID, false, nil
		} else if a.Token.RefreshToken == token {
			return a.ID, true, nil
		}
	}
	return "", false, errors.New("invalid token")
}
