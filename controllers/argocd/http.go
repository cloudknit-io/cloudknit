/* Copyright (C) 2020 CompuZest, Inc. - All Rights Reserved
 *
 * Unauthorized copying of this file, via any medium, is strictly prohibited
 * Proprietary and confidential
 *
 * NOTICE: All information contained herein is, and remains the property of
 * CompuZest, Inc. The intellectual and technical concepts contained herein are
 * proprietary to CompuZest, Inc. and are protected by trade secret or copyright
 * law. Dissemination of this information or reproduction of this material is
 * strictly forbidden unless prior written permission is obtained from CompuZest, Inc.
 */

package argocd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	appv1 "github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
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
	if err != nil {
		return nil, err
	}

	getTokenUrl := fmt.Sprintf("%s/api/v1/session", api.ServerUrl)
	resp, err := http.Post(getTokenUrl, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to send POST request to /api/v1/session: %v", err)
	}

	respBody, err := common.ReadBody(resp.Body)
	if err != nil {
		return nil, err
	}

	t := new(GetTokenResponse)
	if err := common.FromJson(t, respBody); err != nil {
		return nil, err
	}

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
	getRepoUrl := fmt.Sprintf("%s/api/v1/repositories", api.ServerUrl)
	req, err := http.NewRequest("GET", getRepoUrl, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create GET request: %v", err)
	}

	req.Header.Add("Authorization", bearerToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to send GET request to /api/v1/repositories: %v", err)
	}

	if resp.StatusCode != 200 {
		common.CloseBody(resp.Body)
		return nil, nil, fmt.Errorf("list repositories returned a non-OK response: %d", resp.StatusCode)
	}

	repos := new(RepositoryList)
	err = json.NewDecoder(resp.Body).Decode(repos)
	if err != nil {
		common.CloseBody(resp.Body)
		return nil, nil, err
	}

	return repos, resp, nil
}

func (api HttpApi) CreateRepository(body CreateRepoBody, bearerToken string) (*http.Response, error) {
	jsonBody, err := common.ToJson(body)
	if err != nil {
		return nil, err
	}

	addRepoUrl := fmt.Sprintf("%s/api/v1/repositories", api.ServerUrl)
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
		common.CloseBody(resp.Body)
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
		common.CloseBody(resp.Body)
		return false, nil
	}

	if resp.StatusCode != 200 {
		common.CloseBody(resp.Body)
		return false, fmt.Errorf("get application returned a non-OK response: %d", resp.StatusCode)
	}

	application := new(appv1.Application)
	if err := json.NewDecoder(resp.Body).Decode(application); err != nil {
		common.CloseBody(resp.Body)
		return false, err
	}

	if application.ObjectMeta.Name == name {
		return true, nil
	}

	return false, nil
}

func (api HttpApi) CreateApplication(application *appv1.Application, bearerToken string) (*http.Response, error) {
	jsonBody, err := common.ToJson(application)
	if err != nil {
		return nil, err
	}

	addAppURL := fmt.Sprintf("%s/api/v1/applications", api.ServerUrl)
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
		common.CloseBody(resp.Body)
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
		common.CloseBody(resp.Body)
		return fmt.Errorf("delete application returned non-OK status code: %d", resp.StatusCode)
	}

	return nil
}
