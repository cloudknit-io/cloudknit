package argocd

import (
	"bytes"
	"encoding/json"
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
		return nil, fmt.Errorf("error getting argocd credentials: %v", err)
	}

	body := GetTokenBody{Username: creds.Username, Password: creds.Password}
	jsonBody, err := common.ToJson(body)

	getTokenUrl := api.ServerUrl + "/api/v1/session"
	resp, err := http.Post(getTokenUrl, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to send POST request to argocd server: %v", err)
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
		return nil, nil, fmt.Errorf("failed to create GET request: %v", err)
	}

	req.Header.Add("Authorization", bearerToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to send GET request to argocd server: %v", err)
	}

	if resp.StatusCode != 200 {
		resp.Body.Close()
		return nil, nil, fmt.Errorf("list repositories returned a non-OK response: %d", resp.StatusCode)
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
		return nil, fmt.Errorf("failed to create POST request: %v", err)
	}

	req.Header.Add("Authorization", bearerToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send POST request to /repositories: %v", err)
	}

	if resp.StatusCode != 200 {
		common.LogBody(api.Log, resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("create repository returned non-OK status code: %d", resp.StatusCode)
	}

	return resp, nil
}

func (api HttpApi) DoesApplicationExist(name string, bearerToken string) (bool, error) {
	getAppURL := fmt.Sprintf("%s/api/v1/applications/%s", api.ServerUrl, name)
	req, err := http.NewRequest("GET", getAppURL, nil)
	if err != nil {
		return false, fmt.Errorf("failed to create GET request: %v", err)
	}

	req.Header.Add("Authorization", bearerToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("failed to send GET request to applications/%s: %v", name, err)
	}

	if resp.StatusCode == 404 {
		return false, nil
	}

	if resp.StatusCode != 200 {
		resp.Body.Close()
		return false, fmt.Errorf("get application returned a non-OK response: %d", resp.StatusCode)
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
		return nil, fmt.Errorf("failed to create POST request: %v", err)
	}

	req.Header.Add("Authorization", bearerToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send POST request to /applications: %v", err)
	}

	if resp.StatusCode != 200 {
		common.LogBody(api.Log, resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("create application returned non-OK status code: %d", resp.StatusCode)
	}

	return resp, nil
}

func (api HttpApi) DeleteApplication(name string, bearerToken string) error {
	deleteAppURL := fmt.Sprintf("%s/api/v1/applications/%s", api.ServerUrl, name)
	req, err := http.NewRequest("DELETE", deleteAppURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create DELETE request: %v", err)
	}

	req.Header.Add("Authorization", bearerToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send DELETE request to /applications/%s: %v", name, err)
	}

	if resp.StatusCode != 200 {
		common.LogBody(api.Log, resp.Body)
		resp.Body.Close()
		return fmt.Errorf("delete application returned non-OK status code: %d", resp.StatusCode)
	}

	return nil
}
