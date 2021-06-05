package argocd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	appv1 "github.com/argoproj/argo-cd/pkg/apis/application/v1alpha1"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/common"
	"github.com/go-logr/logr"
)

func NewHttpClient(l logr.Logger, serverUrl string) Api {
	return HttpApi{Log: l, ServerUrl: serverUrl}
}

func (api HttpApi) GetAuthToken() (*GetTokenResponse, error) {
	creds, err := getArgocdCredentialsFromEnv()
	if err != nil {
		api.Log.Error(err, "Error getting argocd credentials")
		return nil, err
	}

	body := GetTokenBody{Username: creds.Username, Password: creds.Password}
	jsonBody, err := common.ToJson(body)

	getTokenUrl := api.ServerUrl + "/api/v1/session"
	resp, err := http.Post(getTokenUrl, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		api.Log.Error(err, "Failed to send POST request to argocd server", "url", getTokenUrl)
		return nil, err
	}

	respBody, err := common.ReadBody(resp.Body)
	if err != nil {
		return nil, err
	}

	t := new(GetTokenResponse)
	err = common.FromJson(t, respBody)

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

func (api HttpApi) ListRepositories(bearerToken string) (*RepositoryList, *http.Response, error) {
	getRepoUrl := api.ServerUrl + "/api/v1/repositories"
	req, err := http.NewRequest("GET", getRepoUrl, nil)
	if err != nil {
		api.Log.Error(err, "Failed to create GET request")
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

func (api HttpApi) CreateRepository(body CreateRepoBody, bearerToken string) (*http.Response, error) {
	jsonBody, err := common.ToJson(body)

	addRepoUrl := api.ServerUrl + "/api/v1/repositories"
	req, err := http.NewRequest("POST", addRepoUrl, bytes.NewBuffer(jsonBody))
	if err != nil {
		api.Log.Error(err, "Failed to create POST request")
		return nil, err
	}

	req.Header.Add("Authorization", bearerToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		api.Log.Error(
			err, "Failed to send POST request to /repositories",
			"serverUrl", api.ServerUrl,
			"repoUrl", addRepoUrl,
		)
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

func (api HttpApi) DoesApplicationExist(name string, bearerToken string) (bool, error) {
	getAppURL := api.ServerUrl + "/api/v1/applications/" + name
	req, err := http.NewRequest("GET", getAppURL, nil)
	if err != nil {
		api.Log.Error(err, "Failed to create GET request")
		return false, err
	}

	req.Header.Add("Authorization", bearerToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		api.Log.Error(err, "Failed to send GET request to argocd server", "url", getAppURL)
		return false, err
	}

	if resp.StatusCode == 404 {
		return false, nil
	}

	if resp.StatusCode != 200 {
		err := fmt.Errorf("get application returned a non-OK response: %d", resp.StatusCode)
		resp.Body.Close()
		return false, err
	}

	application := new(appv1.Application)
	err = json.NewDecoder(resp.Body).Decode(application)
	if err != nil {
		resp.Body.Close()
		return false, err
	}

	if application.ObjectMeta.Name == name {
		return true, nil
	}

	return false, nil
}

func (api HttpApi) CreateApplication(application *appv1.Application, bearerToken string) (*http.Response, error) {
	jsonBody, err := common.ToJson(application)

	addAppURL := api.ServerUrl + "/api/v1/applications"
	req, err := http.NewRequest("POST", addAppURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		api.Log.Error(err, "Failed to create POST request")
		return nil, err
	}

	req.Header.Add("Authorization", bearerToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		api.Log.Error(
			err, "Failed to send POST request to /applications",
			"serverUrl", api.ServerUrl,
			"repoUrl", addAppURL,
		)
		return nil, err
	}

	if resp.StatusCode != 200 {
		common.LogBody(api.Log, resp.Body)
		err = fmt.Errorf("create application returned non-OK status code: %d", resp.StatusCode)
		resp.Body.Close()
		return nil, err
	}

	return resp, nil
}
