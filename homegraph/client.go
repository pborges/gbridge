package homegraph

import (
	"context"
	"encoding/json"
	"errors"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"net/http"
	"strings"
)

const ApiEndpoint = "https://homegraph.googleapis.com/v1/"
const ApiScope = "https://www.googleapis.com/auth/homegraph"

type Response struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Status  string `json:"status"`
	} `json:"error"`
}

func Dial(credentialsJsonFile string) (Client, error) {
	data, err := ioutil.ReadFile(credentialsJsonFile)
	if err != nil {
		return Client{}, err
	}

	conf, err := google.JWTConfigFromJSON(data, ApiScope)
	if err != nil {
		return Client{}, err
	}

	return Client{httpClient: conf.Client(context.Background())}, nil
}

type Client struct {
	httpClient *http.Client
}

func (c Client) RequestResync(agentUserId string) error {
	res, err := c.httpClient.Post(
		ApiEndpoint+"devices:requestSync",
		"application/json",
		strings.NewReader(`{"agentUserId": "`+agentUserId+`", "async":false}`),
	)
	if err != nil {
		return err
	}

	var response Response
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return err
	}

	if err := res.Body.Close(); err != nil {
		return err
	}

	if response.Error.Code != 0 {
		return errors.New(response.Error.Message)
	}

	return nil
}
