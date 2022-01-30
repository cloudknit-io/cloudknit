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
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	url2 "net/url"

	appv1 "github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/common"
	"github.com/go-logr/logr"
)

func NewHTTPClient(ctx context.Context, l logr.Logger, serverURL string) API {
	return &HTTPAPI{ctx: ctx, log: l, serverURL: serverURL}
}

func (api *HTTPAPI) GetAuthToken() (*GetTokenResponse, error) {
	creds, err := getArgocdCredentialsFromEnv()
	if err != nil {
		return nil, fmt.Errorf("error getting argocd credentials: %w", err)
	}

	body := GetTokenBody{Username: creds.Username, Password: creds.Password}
	jsonBody, err := common.ToJSON(body)
	if err != nil {
		return nil, errors.Wrap(err, "error marshalling body to JSON")
	}

	url := fmt.Sprintf("%s/api/v1/session", api.serverURL)
	req, err := http.NewRequestWithContext(api.ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, errors.Wrap(err, "error creating POST request")
	}

	client := common.GetHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "error sending POST request to %s", url)
	}
	defer common.CloseBody(resp.Body)

	respBody, err := common.ReadBody(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "error reading response body")
	}

	t := new(GetTokenResponse)
	if err := common.FromJSON(t, respBody); err != nil {
		return nil, errors.Wrap(err, "error unmarshalling response body")
	}

	return t, nil
}

func isRepoRegistered(repos *RepositoryList, repoURL string) bool {
	for _, r := range repos.Items {
		if r.Repo == repoURL {
			return true
		}
	}
	return false
}

func (api *HTTPAPI) ListRepositories(bearerToken string) (*RepositoryList, *http.Response, error) {
	url := fmt.Sprintf("%s/api/v1/repositories", api.serverURL)
	req, err := http.NewRequestWithContext(api.ctx, "GET", url, nil)
	if err != nil {
		return nil, nil, errors.Wrap(err, "error creating GET request")
	}

	req.Header.Add("Authorization", bearerToken)

	client := common.GetHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "error sending GET request to %s", url)
	}

	if resp.StatusCode != 200 {
		common.CloseBody(resp.Body)
		return nil, nil, errors.Errorf("list repositories returned a non-OK response: %d", resp.StatusCode)
	}

	repos := new(RepositoryList)
	err = json.NewDecoder(resp.Body).Decode(repos)
	if err != nil {
		common.CloseBody(resp.Body)
		return nil, nil, errors.Wrap(err, "error unmarshalling response body")
	}

	return repos, resp, nil
}

func (api *HTTPAPI) CreateRepository(body *CreateRepoBody, bearerToken string) (*http.Response, error) {
	jsonBody, err := common.ToJSON(body)
	if err != nil {
		return nil, errors.Wrap(err, "error marshalling body to JSON")
	}

	url := fmt.Sprintf("%s/api/v1/repositories", api.serverURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, errors.Wrap(err, "error creating POST request")
	}

	req.Header.Add("Authorization", bearerToken)

	client := common.GetHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "error sending POST request to %s", url)
	}

	if resp.StatusCode != 200 {
		common.LogBody(api.log, resp.Body)
		common.CloseBody(resp.Body)
		return nil, errors.Errorf("create repository returned non-OK status code: %d", resp.StatusCode)
	}

	return resp, nil
}

func (api *HTTPAPI) DoesApplicationExist(name string, bearerToken string) (bool, error) {
	url := fmt.Sprintf("%s/api/v1/applications/%s", api.serverURL, name)
	req, err := http.NewRequestWithContext(api.ctx, "GET", url, nil)
	if err != nil {
		return false, errors.Wrap(err, "error creating GET request")
	}

	req.Header.Add("Authorization", bearerToken)

	client := common.GetHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return false, errors.Wrapf(err, "error sending GET request to %s", url)
	}
	defer common.CloseBody(resp.Body)

	if resp.StatusCode == 404 {
		common.CloseBody(resp.Body)
		return false, nil
	}

	if resp.StatusCode != 200 {
		common.CloseBody(resp.Body)
		return false, errors.Errorf("get application returned a non-OK response: %d", resp.StatusCode)
	}

	application := new(appv1.Application)
	if err := json.NewDecoder(resp.Body).Decode(application); err != nil {
		common.CloseBody(resp.Body)
		return false, errors.Wrap(err, "error unmarshalling response body")
	}

	if application.ObjectMeta.Name == name {
		return true, nil
	}

	return false, nil
}

func (api *HTTPAPI) CreateApplication(application *appv1.Application, bearerToken string) (*http.Response, error) {
	jsonBody, err := common.ToJSON(application)
	if err != nil {
		return nil, errors.Wrap(err, "error marshalling body to JSON")
	}

	url := fmt.Sprintf("%s/api/v1/applications", api.serverURL)
	req, err := http.NewRequestWithContext(api.ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, errors.Wrap(err, "error creating POST request")
	}

	req.Header.Add("Authorization", bearerToken)

	client := common.GetHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "error sending POST request to %s", url)
	}

	if resp.StatusCode != 200 {
		common.LogBody(api.log, resp.Body)
		common.CloseBody(resp.Body)
		return nil, errors.Errorf("create application returned non-OK status code: %d", resp.StatusCode)
	}

	return resp, nil
}

func (api *HTTPAPI) DeleteApplication(name string, bearerToken string) error {
	url := fmt.Sprintf("%s/api/v1/applications/%s", api.serverURL, name)
	req, err := http.NewRequestWithContext(api.ctx, "DELETE", url, nil)
	if err != nil {
		return errors.Wrap(err, "error creating DELETE request")
	}

	req.Header.Add("Authorization", bearerToken)

	client := common.GetHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrapf(err, "error sending DELETE request to %s", url)
	}
	defer common.CloseBody(resp.Body)

	if resp.StatusCode != 200 {
		common.LogBody(api.log, resp.Body)
		return errors.Errorf("delete application returned non-OK status code: %d", resp.StatusCode)
	}

	return nil
}

func (api *HTTPAPI) CreateProject(project *CreateProjectBody, bearerToken string) (*http.Response, error) {
	jsonBody, err := common.ToJSON(project)
	if err != nil {
		return nil, errors.Wrap(err, "error marshalling body to JSON")
	}

	url := fmt.Sprintf("%s/api/v1/projects", api.serverURL)
	req, err := http.NewRequestWithContext(api.ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, errors.Wrap(err, "error creating POST request")
	}

	req.Header.Add("Authorization", bearerToken)

	client := common.GetHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "error sending POST request to %s", url)
	}

	if resp.StatusCode != 200 {
		common.LogBody(api.log, resp.Body)
		common.CloseBody(resp.Body)
		return nil, errors.Errorf("create project returned non-OK status code: %d", resp.StatusCode)
	}

	return resp, nil
}

func (api *HTTPAPI) DoesProjectExist(name string, bearerToken string) (exists bool, response *http.Response, err error) {
	url := fmt.Sprintf("%s/api/v1/projects/%s", api.serverURL, name)
	req, err := http.NewRequestWithContext(api.ctx, "GET", url, nil)
	if err != nil {
		return false, nil, errors.Wrap(err, "error creating GET request")
	}

	req.Header.Add("Authorization", bearerToken)

	client := common.GetHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return false, nil, errors.Wrapf(err, "error sending GET request to %s", url)
	}

	if resp.StatusCode == 404 {
		common.LogBody(api.log, resp.Body)
		return false, resp, nil
	}

	return true, resp, nil
}

func (api *HTTPAPI) UpdateCluster(
	clusterURL string,
	body *UpdateClusterBody,
	updatedFields []string,
	bearerToken string,
) (*http.Response, error) {
	jsonBody, err := common.ToJSON(body)
	if err != nil {
		return nil, errors.Wrap(err, "error marshaling body to JSON")
	}

	url := fmt.Sprintf("%s/api/v1/clusters/%s?updatedFields=%s", api.serverURL, url2.QueryEscape(clusterURL), common.ToStringArray(updatedFields))
	req, err := http.NewRequestWithContext(api.ctx, "PUT", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, errors.Wrap(err, "error creating PUT request")
	}

	req.Header.Add("Authorization", bearerToken)

	client := common.GetHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "error sending PUT request to %s", url)
	}

	if resp.StatusCode == 404 {
		common.LogBody(api.log, resp.Body)
		return nil, errors.Errorf("update cluster returned non-OK response: %d", resp.StatusCode)
	}

	return resp, nil
}
