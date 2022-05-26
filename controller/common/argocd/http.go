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
	"net/http"
	url2 "net/url"

	"github.com/sirupsen/logrus"

	"github.com/compuzest/zlifecycle-il-operator/controller/util"

	"github.com/pkg/errors"

	appv1 "github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
)

func NewHTTPClient(ctx context.Context, l *logrus.Entry, serverURL string) API {
	return &HTTPClient{ctx: ctx, log: l, serverURL: serverURL}
}

func (c *HTTPClient) GetAuthToken() (*GetTokenResponse, error) {
	creds, err := getArgocdCredentialsFromEnv()
	if err != nil {
		return nil, fmt.Errorf("error getting argocd credentials: %w", err)
	}

	body := GetTokenBody{Username: creds.Username, Password: creds.Password}
	jsonBody, err := util.ToJSON(body)
	if err != nil {
		return nil, errors.Wrap(err, "error marshaling body to JSON")
	}

	url := fmt.Sprintf("%s/api/v1/session", c.serverURL)
	req, err := http.NewRequestWithContext(c.ctx, http.MethodPost, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, errors.Wrap(err, "error creating POST request")
	}

	client := util.GetHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "error sending POST request to %s", url)
	}
	defer util.CloseBody(resp.Body)

	respBody, err := util.ReadBody(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "error reading get auth token response body")
	}

	t := GetTokenResponse{}
	if err = util.FromJSON(&t, respBody); err != nil {
		return nil, errors.Wrap(err, "error unmarshalling get auth token response body")
	}

	return &t, nil
}

func (c *HTTPClient) ListRepositories(bearerToken string) (*RepositoryList, *http.Response, error) {
	url := fmt.Sprintf("%s/api/v1/repositories", c.serverURL)
	req, err := http.NewRequestWithContext(c.ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, nil, errors.Wrap(err, "error creating GET request")
	}

	req.Header.Add("Authorization", bearerToken)

	client := util.GetHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "error sending GET request to %s", url)
	}

	if resp.StatusCode != 200 {
		util.CloseBody(resp.Body)
		return nil, nil, errors.Errorf("list repositories returned a non-OK response: %d", resp.StatusCode)
	}

	repos := RepositoryList{}
	if err = json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		util.CloseBody(resp.Body)
		return nil, nil, errors.Wrap(err, "error unmarshalling response body")
	}

	return &repos, resp, nil
}

func (c *HTTPClient) CreateRepository(body interface{}, bearerToken string) (*http.Response, error) {
	jsonBody, err := util.ToJSON(body)
	if err != nil {
		return nil, errors.Wrap(err, "error marshaling body to JSON")
	}

	url := fmt.Sprintf("%s/api/v1/repositories", c.serverURL)
	req, err := http.NewRequestWithContext(c.ctx, http.MethodPost, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, errors.Wrap(err, "error creating POST request")
	}

	req.Header.Add("Authorization", bearerToken)

	client := util.GetHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "error sending POST request to %s", url)
	}

	if resp.StatusCode != 200 {
		util.LogBody(c.log, resp.Body)
		util.CloseBody(resp.Body)
		return nil, errors.Errorf("create repository returned non-OK status code: %d", resp.StatusCode)
	}

	return resp, nil
}

func (c *HTTPClient) DoesApplicationExist(name string, bearerToken string) (bool, error) {
	url := fmt.Sprintf("%s/api/v1/applications/%s", c.serverURL, name)
	req, err := http.NewRequestWithContext(c.ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return false, errors.Wrap(err, "error creating GET request")
	}

	req.Header.Add("Authorization", bearerToken)

	client := util.GetHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return false, errors.Wrapf(err, "error sending GET request to %s", url)
	}
	defer util.CloseBody(resp.Body)

	if resp.StatusCode == 404 {
		util.CloseBody(resp.Body)
		return false, nil
	}

	if resp.StatusCode != 200 {
		util.CloseBody(resp.Body)
		return false, errors.Errorf("get application returned a non-OK response: %d", resp.StatusCode)
	}

	application := appv1.Application{}
	if err := json.NewDecoder(resp.Body).Decode(&application); err != nil {
		util.CloseBody(resp.Body)
		return false, errors.Wrap(err, "error unmarshalling response body")
	}

	if application.ObjectMeta.Name == name {
		return true, nil
	}

	return false, nil
}

func (c *HTTPClient) CreateApplication(application *appv1.Application, bearerToken string) (*http.Response, error) {
	jsonBody, err := util.ToJSON(application)
	if err != nil {
		return nil, errors.Wrap(err, "error marshaling body to JSON")
	}

	url := fmt.Sprintf("%s/api/v1/applications", c.serverURL)
	req, err := http.NewRequestWithContext(c.ctx, http.MethodPost, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, errors.Wrap(err, "error creating POST request")
	}

	req.Header.Add("Authorization", bearerToken)

	client := util.GetHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "error sending POST request to %s", url)
	}

	if resp.StatusCode != 200 {
		util.LogBody(c.log, resp.Body)
		util.CloseBody(resp.Body)
		return nil, errors.Errorf("create application returned non-OK status code: %d", resp.StatusCode)
	}

	return resp, nil
}

func (c *HTTPClient) DeleteApplication(name string, bearerToken string) error {
	url := fmt.Sprintf("%s/api/v1/applications/%s", c.serverURL, name)
	req, err := http.NewRequestWithContext(c.ctx, http.MethodDelete, url, http.NoBody)
	if err != nil {
		return errors.Wrap(err, "error creating DELETE request")
	}

	req.Header.Add("Authorization", bearerToken)

	client := util.GetHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrapf(err, "error sending DELETE request to %s", url)
	}
	defer util.CloseBody(resp.Body)

	if resp.StatusCode != 200 {
		util.LogBody(c.log, resp.Body)
		return errors.Errorf("delete application returned non-OK status code: %d", resp.StatusCode)
	}

	return nil
}

func (c *HTTPClient) CreateProject(project *CreateProjectBody, bearerToken string) (*http.Response, error) {
	jsonBody, err := util.ToJSON(project)
	if err != nil {
		return nil, errors.Wrap(err, "error marshaling body to JSON")
	}

	url := fmt.Sprintf("%s/api/v1/projects", c.serverURL)
	req, err := http.NewRequestWithContext(c.ctx, http.MethodPost, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, errors.Wrap(err, "error creating POST request")
	}

	req.Header.Add("Authorization", bearerToken)

	client := util.GetHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "error sending POST request to %s", url)
	}

	if resp.StatusCode != 200 {
		util.LogBody(c.log, resp.Body)
		util.CloseBody(resp.Body)
		return nil, errors.Errorf("create project returned non-OK status code: %d", resp.StatusCode)
	}

	return resp, nil
}

func (c *HTTPClient) DoesProjectExist(name string, bearerToken string) (exists bool, response *http.Response, err error) {
	url := fmt.Sprintf("%s/api/v1/projects/%s", c.serverURL, name)
	req, err := http.NewRequestWithContext(c.ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return false, nil, errors.Wrap(err, "error creating GET request")
	}

	req.Header.Add("Authorization", bearerToken)

	client := util.GetHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return false, nil, errors.Wrapf(err, "error sending GET request to %s", url)
	}

	if resp.StatusCode == 404 {
		util.LogBody(c.log, resp.Body)
		return false, resp, nil
	}

	return true, resp, nil
}

func (c *HTTPClient) UpdateCluster(
	clusterURL string,
	body *UpdateClusterBody,
	updatedFields []string,
	bearerToken string,
) (*http.Response, error) {
	jsonBody, err := util.ToJSON(body)
	if err != nil {
		return nil, errors.Wrap(err, "error marshaling body to JSON")
	}

	url := fmt.Sprintf("%s/api/v1/clusters/%s?updatedFields=%s", c.serverURL, url2.QueryEscape(clusterURL), util.ToStringArray(updatedFields))
	req, err := http.NewRequestWithContext(c.ctx, http.MethodPut, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, errors.Wrap(err, "error creating PUT request")
	}

	req.Header.Add("Authorization", bearerToken)

	client := util.GetHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "error sending PUT request to %s", url)
	}

	if resp.StatusCode == 404 {
		util.LogBody(c.log, resp.Body)
		return nil, errors.Errorf("update cluster returned non-OK response: %d", resp.StatusCode)
	}

	return resp, nil
}

func (c *HTTPClient) RegisterCluster(body *RegisterClusterBody, bearerToken string) (*http.Response, error) {
	jsonBody, err := util.ToJSON(body)
	if err != nil {
		return nil, errors.Wrap(err, "error marshaling body to JSON")
	}

	url := fmt.Sprintf("%s/api/v1/clusters", c.serverURL)
	req, err := http.NewRequestWithContext(c.ctx, http.MethodPost, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, errors.Wrap(err, "error creating POST request")
	}

	req.Header.Add("Authorization", bearerToken)

	client := util.GetHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "error sending POST request to %s", url)
	}

	if resp.StatusCode != 200 {
		util.LogBody(c.log, resp.Body)
		return nil, errors.Errorf("register cluster returned non-OK response: %d", resp.StatusCode)
	}

	return resp, nil
}

func (c *HTTPClient) ListClusters(name *string, bearerToken string) (*ListClustersResponse, error) {
	url := fmt.Sprintf("%s/api/v1/clusters", c.serverURL)
	if name != nil {
		url += fmt.Sprintf("?name=%s", *name)
	}
	req, err := http.NewRequestWithContext(c.ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, errors.Wrap(err, "error creating GET request")
	}

	req.Header.Add("Authorization", bearerToken)

	client := util.GetHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "error sending GET request to %s", url)
	}
	defer util.CloseBody(resp.Body)

	if resp.StatusCode != 200 {
		util.LogBody(c.log, resp.Body)
		return nil, errors.Errorf("list clusters returned non-OK response: %d", resp.StatusCode)
	}

	respBody, err := util.ReadBody(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "error reading list clusters response body")
	}

	clusters := ListClustersResponse{}
	if err = util.FromJSON(&clusters, respBody); err != nil {
		return nil, errors.Wrap(err, "error unmarshalling list clusters response body")
	}

	return &clusters, nil
}
