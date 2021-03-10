package argocd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/common"
	"github.com/go-logr/logr"
	"net/http"
)

func NewHttpClient(l logr.Logger) ArgocdAPI {
	return ArgocdHttpAPI{Log: l}
}

func (api ArgocdHttpAPI) GetAuthToken(argocdUrl string) (*GetTokenResponse, error) {
	creds, err := getArgocdCredentialsFromEnv()
	if err != nil {
		api.Log.Error(err, "Error getting argocd credentials")
		return nil, err
	}

	body := GetTokenBody{Username: creds.Username, Password: creds.Password}
	jsonBody, err := common.ToJson(api.Log, body)

	getTokenUrl := argocdUrl + "/api/v1/session"
	resp, err := http.Post(getTokenUrl, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		api.Log.Error(err, "Failed to send POST request to argocd server", "url", getTokenUrl)
		return nil, err
	}

	respBody, err := common.ReadBody(api.Log, resp.Body)
	if err != nil {
		return nil, err
	}

	t := new(GetTokenResponse)
	err = common.FromJson(api.Log, t, respBody)

	return t, nil
}

func isRepoRegistered(repos RepositoryList, repoUrl string) bool {
	for _, r := range repos.Items {
		if r.Repo == repoUrl {
			return true
		}
	}
	return false
}

func (api ArgocdHttpAPI) ListRepositories(host string, bearerToken string) (*RepositoryList, *http.Response, error) {
	getRepoUrl := host + "/api/v1/repositories"
	req, err := http.NewRequest("GET", getRepoUrl, nil)
	if err != nil {
		api.Log.Error(err, "Failed to create POST request")
		return nil, nil, err
	}

	req.Header.Add("Authorization", bearerToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		api.Log.Error(err, "Failed to send GET request to argocd server", "url", getRepoUrl)
		return nil, nil, err
	}

	if resp.StatusCode != 200 {
		err := errors.New(
			fmt.Sprintf("list repositories returned a non-OK response: %d", resp.StatusCode),
		)
		resp.Body.Close()
		return nil, nil, err
	}

	repos := new(RepositoryList)
	err = json.NewDecoder(resp.Body).Decode(repos)
	if err != nil {
		resp.Body.Close()
		return nil, nil, err
	}

	return repos, resp, nil
}

func (api ArgocdHttpAPI) CreateRepository(serverUrl string, body CreateRepoBody, bearerToken string) (*http.Response, error) {
	jsonBody, err := common.ToJson(api.Log, body)

	addRepoUrl := serverUrl + "/api/v1/repositories"
	req, err := http.NewRequest("POST", addRepoUrl, bytes.NewBuffer(jsonBody))
	if err != nil {
		api.Log.Error(err, "Failed to create POST request")
		return nil, err
	}

	req.Header.Add("Authorization", bearerToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		api.Log.Error(err, "Failed to send POST request to /repositories", "server", serverUrl, "repoUrl", addRepoUrl)
		return nil, err
	}

	if resp.StatusCode != 200 {
		common.LogBody(api.Log, resp.Body)
		err = errors.New(fmt.Sprintf("create repository returned non-OK status code: %d", resp.StatusCode))
		resp.Body.Close()
		return nil, err
	}

	return resp, nil
}
