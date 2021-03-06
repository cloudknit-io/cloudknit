package argocd

import (
	"bytes"
	"github.com/go-logr/logr"
	"net/http"
)

type ArgocdCredentials struct {
	Username string
	Password string
}

type GetTokenBody struct {
	Username string
	Password string
}

type GetTokenResponse struct {
	Token string
}

func GetAuthToken(log logr.Logger, argocdUrl string) (*GetTokenResponse, error) {
	creds, err := getArgocdCredentialsFromEnv()
	if err != nil {
		log.Error(err, "Error getting argocd credentials")
		return nil, err
	}

	body := GetTokenBody{Username: creds.Username, Password: creds.Password}
	jsonBody, err := toJson(log, body)

	getTokenUrl := argocdUrl + "/api/v1/session"
	resp, err := http.Post(getTokenUrl, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Error(err, "Failed to send POST request to argocd server", "url", getTokenUrl)
		return nil, err
	}

	respBody, err := readBody(log, resp.Body)
	if err != nil {
		return nil, err
	}

	t := new(GetTokenResponse)
	err = fromJson(log, t, respBody)

	return t, nil
}
