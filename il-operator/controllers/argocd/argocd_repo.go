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
	"errors"
	"fmt"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/common"
	"github.com/go-logr/logr"
	"net/http"
	"strings"
)

type RepoOpts struct {
	ServerUrl     string
	RepoUrl       string
	SshPrivateKey string
}

type CreateRepoBody struct {
	Repo          string `json:"repo"`
	Name          string `json:"name"`
	SshPrivateKey string `json:"sshPrivateKey"`
}

type RepositoryList struct {
	Items []Repository `json:"items"`
}

type Repository struct {
	Repo string `json:"repo"`
	Name string `json:"name"`
}

func RegisterRepo(log logr.Logger, repoOpts RepoOpts) error {
	repoUri := repoOpts.RepoUrl[strings.LastIndex(repoOpts.RepoUrl, "/")+1:]
	repoName := strings.TrimSuffix(repoUri, ".git")

	tokenResponse, err := GetAuthToken(log, repoOpts.ServerUrl)
	if err != nil {
		log.Error(err, "Error while calling get auth token")
		return err
	}

	bearer := "Bearer " + tokenResponse.Token
	repositories, resp1, err1 := listRepositories(log, repoOpts.ServerUrl, bearer)
	if err1 != nil {
		return err1
	}
	defer resp1.Body.Close()

	if isRepoRegistered(*repositories, repoOpts.RepoUrl) {
		log.Info("Repository already registered",
			"repos", repositories.Items,
			"repoName", repoName,
			"repoUrl", repoOpts.RepoUrl,
		)
		return nil
	}

	log.Info("Repository is not registered, registering now...", "repos", repositories, "repoName", repoName)

	createRepoBody := CreateRepoBody{Repo: repoOpts.RepoUrl, Name: repoName, SshPrivateKey: repoOpts.SshPrivateKey}
	resp2, err2 := postCreateRepository(log, repoOpts.ServerUrl, createRepoBody, bearer)
	if err2 != nil {
		log.Error(err2, "Error while calling post create repository")
		return err2
	}
	defer resp2.Body.Close()

	log.Info("Successfully registered repository", "repo", repoOpts.RepoUrl)

	return nil
}

func isRepoRegistered(repos RepositoryList, repoUrl string) bool {
	for _, r := range repos.Items {
		if r.Repo == repoUrl {
			return true
		}
	}
	return false
}

func listRepositories(log logr.Logger, host string, bearerToken string) (*RepositoryList, *http.Response, error) {
	getRepoUrl := host + "/api/v1/repositories"
	req, err := http.NewRequest("GET", getRepoUrl, nil)
	if err != nil {
		log.Error(err, "Failed to create POST request")
		return nil, nil, err
	}

	req.Header.Add("Authorization", bearerToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err, "Failed to send GET request to argocd server", "url", getRepoUrl)
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

func postCreateRepository(log logr.Logger, serverUrl string, body CreateRepoBody, bearerToken string) (*http.Response, error) {
	jsonBody, err := common.ToJson(log, body)

	addRepoUrl := serverUrl + "/api/v1/repositories"
	req, err := http.NewRequest("POST", addRepoUrl, bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Error(err, "Failed to create POST request")
		return nil, err
	}

	req.Header.Add("Authorization", bearerToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err, "Failed to send POST request to /repositories", "server", serverUrl, "repoUrl", addRepoUrl)
		return nil, err
	}

	if resp.StatusCode != 200 {
		common.LogBody(log, resp.Body)
		err = errors.New(fmt.Sprintf("create repository returned non-OK status code: %d", resp.StatusCode))
		resp.Body.Close()
		return nil, err
	}

	return resp, nil
}
