package oauth

import (
	"errors"
	"math/rand"
	"sync"
	"time"
)

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type AuthenticationProvider interface {
	GenerateAuthCodeForClient(clientId string) (authCode string, err error)
	ValidateAuthCodeAndGenerateToken(clientId, clientSecret, authCode string) (token Token, err error)
}

type mapBasedAuthProviderClient struct {
	ID           string
	Secret       string
	AuthCode     string
	AccessToken  string
	RefreshToken string
}

type MapBasedAuthProvider struct {
	db   map[string]mapBasedAuthProviderClient
	lock sync.Mutex
}

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateRandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func (m *MapBasedAuthProvider) RegisterClient(clientId string, clientSecret string) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if m.db == nil {
		m.db = make(map[string]mapBasedAuthProviderClient)
	}
	m.db[clientId] = mapBasedAuthProviderClient{
		ID:     clientId,
		Secret: clientSecret,
	}
}

func (m *MapBasedAuthProvider) GenerateAuthCodeForClient(clientId string) (string, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if m.db != nil {
		if client, ok := m.db[clientId]; ok {
			client.AuthCode = GenerateRandomString(36)
			m.db[clientId] = client
			return client.AuthCode, nil
		}
	}

	return "", errors.New("unknown client")
}

func (m *MapBasedAuthProvider) ValidateAuthCodeAndGenerateToken(clientId, clientSecret, authCode string) (Token, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if m.db != nil {
		if client, ok := m.db[clientId]; ok {
			if client.Secret != clientSecret {
				return Token{}, errors.New("invalid secret")
			}
			if client.AuthCode != authCode {
				return Token{}, errors.New("invalid auth code")
			}
			client.AccessToken = GenerateRandomString()
			client.RefreshToken = GenerateRandomString()

			m.db[clientId] = client
			return Token{
				AccessToken:  client.AccessToken,
				RefreshToken: client.RefreshToken,
			}, nil
		}
	}

	return Token{}, errors.New("unknown client")
}
