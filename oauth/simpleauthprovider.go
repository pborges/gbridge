package oauth

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

// a very simple authentication provider
type SimpleAuthenticationProvider struct {
	db   simpleAuthData
	lock sync.Mutex
	file string
}

type simpleAuthData struct {
	Clients []simpleAuthProviderClient `json:"clients"` // key: clientId
	Agents  []simpleAuthProviderAgent  `json:"agents"`  // key: agentUserId
}

func (db *simpleAuthData) putAgent(agent simpleAuthProviderAgent) {
	for idx, a := range db.Agents {
		if agent.ID == a.ID {
			db.Agents[idx] = agent
			return
		}
	}
	db.Agents = append(db.Agents, agent)
}

func (db *simpleAuthData) findAgent(id string) (simpleAuthProviderAgent, error) {
	for _, a := range db.Agents {
		if a.ID == id {
			return a, nil
		}
	}
	return simpleAuthProviderAgent{}, errors.New("unable to find agent")
}

func (db *simpleAuthData) putClient(client simpleAuthProviderClient) {
	for idx, a := range db.Clients {
		if client.ID == a.ID {
			db.Clients[idx] = client
			return
		}
	}
	db.Clients = append(db.Clients, client)
}

func (db *simpleAuthData) findClient(id string) (simpleAuthProviderClient, error) {
	for _, a := range db.Clients {
		if a.ID == id {
			return a, nil
		}
	}
	return simpleAuthProviderClient{}, errors.New("unable to find client")
}

type simpleAuthProviderClient struct {
	ID        string                       `json:"id"`
	Secret    string                       `json:"secret"`
	authCodes []simpleAuthProviderAuthCode // dont write these to database
}

func (a *simpleAuthProviderClient) putAuthCode(agentUserId, authCode string) {
	for _, a := range a.authCodes {
		if a.AgentUserId == agentUserId && a.AuthCode == authCode {
			return
		}
	}
	a.authCodes = append(a.authCodes, simpleAuthProviderAuthCode{AgentUserId: agentUserId, AuthCode: authCode})
}

func (a *simpleAuthProviderClient) findAuthCode(authCode string) string {
	for _, a := range a.authCodes {
		if a.AuthCode == authCode {
			return a.AgentUserId
		}
	}
	return ""
}

type simpleAuthProviderAgent struct {
	ID       string `json:"id"`
	Password string `json:"password"` // todo hash this
	Token    Token  `json:"token"`    // todo support more then one token, or store token on the "link" between agent and client
}

type simpleAuthProviderAuthCode struct {
	AgentUserId string `json:"agentUserId"`
	AuthCode    string `json:"authCode"`
}

// Init attempts to write a file to disk, or load an existing file the bool returned indicates if data was loaded
// i.e an existing database
func (m *SimpleAuthenticationProvider) Init(file string) (bool, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.file = file
	if _, err := os.Stat(m.file); os.IsNotExist(err) {
		return false, nil
	}

	if err := m.withFile(false, func(fd *os.File) error {
		return json.NewDecoder(fd).Decode(&m.db)
	}); err != nil {
		return false, err
	}

	return true, nil
}

// Save truncates the underlying file without locking
func (m *SimpleAuthenticationProvider) save() error {
	return m.withFile(true, func(fd *os.File) error {
		return json.NewEncoder(fd).Encode(m.db)
	})
}

// Save truncates the underlying file
func (m *SimpleAuthenticationProvider) Save() error {
	m.lock.Lock()
	defer m.lock.Unlock()
	return m.save()
}

// Delete the underlying file
func (m *SimpleAuthenticationProvider) Delete() error {
	if m.file == "" {
		return errors.New("file has not been set")
	}
	return os.Remove(m.file)
}

func (m *SimpleAuthenticationProvider) withFile(write bool, fn func(fd *os.File) error) error {
	if m.file == "" {
		return errors.New("file has not been set")
	}

	flags := os.O_RDONLY
	if write {
		flags = os.O_CREATE | os.O_WRONLY | os.O_TRUNC
	}

	fd, err := os.OpenFile(m.file, flags, 0755)
	if err != nil {
		return err
	}
	defer fd.Close()

	return fn(fd)
}

func (m *SimpleAuthenticationProvider) RegisterClient(clientId string, clientSecret string) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.db.putClient(simpleAuthProviderClient{
		ID:     clientId,
		Secret: clientSecret,
	})
	return m.save()
}

func (m *SimpleAuthenticationProvider) RegisterAgent(agentUserId string, password string) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	agent := simpleAuthProviderAgent{
		ID:       agentUserId,
		Password: password,
		Token: Token{
			AccessToken:  GenerateRandomString(36),
			RefreshToken: GenerateRandomString(36),
		},
	}
	m.db.putAgent(agent)
	return m.save()
}

func (m *SimpleAuthenticationProvider) ValidateAgent(agentUserId string, password string) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	if agent, err := m.db.findAgent(agentUserId); err == nil {
		if agent.Password == password {
			return nil
		} else {
			return errors.New("invalid credentials")
		}
	} else {
		return err
	}
}

func (m *SimpleAuthenticationProvider) GenerateAuthCode(clientId string, agentUserId string) (string, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	// do we have a valid client?
	if client, err := m.db.findClient(clientId); err == nil {
		// do we have a valid agent?
		if agent, err := m.db.findAgent(agentUserId); err == nil {
			authCode := GenerateRandomString(36)
			client.putAuthCode(agent.ID, authCode)
			m.db.putClient(client)
			return authCode, nil
		} else {
			return "", err
		}
	}

	return "", errors.New("unknown client")
}

func (m *SimpleAuthenticationProvider) GenerateToken(clientId string, clientSecret string, authCode string) (Token, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	// do we have a valid client?
	if client, err := m.db.findClient(clientId); err == nil {
		// with a valid secret
		if client.Secret != clientSecret {
			return Token{}, errors.New("invalid client secret")
		}

		// do we have a valid agent?
		if agentUserId := client.findAuthCode(authCode); agentUserId != "" {
			if agent, err := m.db.findAgent(agentUserId); err == nil {
				return agent.Token, nil
			}
			return Token{}, err
		} else {
			return Token{}, errors.New("unknown agent")
		}
	}

	return Token{}, errors.New("unknown client")
}

func (m *SimpleAuthenticationProvider) GetAgentUserIDForToken(token string) (agentUserId string, refresh bool, err error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	for _, a := range m.db.Agents {
		if a.Token.AccessToken == token {
			return a.ID, false, nil
		} else if a.Token.RefreshToken == token {
			return a.ID, true, nil
		}
	}
	return "", false, errors.New("invalid token")
}
