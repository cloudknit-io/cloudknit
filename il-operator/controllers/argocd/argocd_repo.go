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
	"errors"
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
	Repo          string
	Name          string
	SshPrivateKey string
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
	resp1, err1 := getRepository(log, repoOpts.ServerUrl, repoName, bearer)
	if err1 != nil {
		log.Error(err1, "Error while calling get repository")
		return err1
	}
	defer resp1.Body.Close()

	if resp1.StatusCode == 200 {
		log.Info("Repository already registered", "repository", repoOpts.RepoUrl)
		return nil
	}

	createRepoBody := CreateRepoBody{Repo: repoOpts.RepoUrl, Name: repoName, SshPrivateKey: repoOpts.SshPrivateKey}
	resp2, err2 := postCreateRepository(log, repoOpts.ServerUrl, createRepoBody, bearer)
	if err2 != nil {
		log.Error(err2, "Error while calling post create repository")
		return err2
	}
	defer resp2.Body.Close()

	if resp2.StatusCode != 200 {
		common.LogBody(log, resp2.Body)
		err2 = errors.New("status code does not equal 200")
		log.Error(err2, "Add new repo request failed", "status code", resp2.StatusCode, "response", resp2.Body)
		return err2
	}

	log.Info("Successfully registered repository", "repo", repoOpts.RepoUrl)

	return nil
}

func getRepository(log logr.Logger, host string, repoName string, bearerToken string) (*http.Response, error) {
	getRepoUrl := host + "/api/v1/repositories/" + repoName
	req, err := http.NewRequest("GET", getRepoUrl, nil)
	if err != nil {
		log.Error(err, "Failed to create POST request")
		return nil, err
	}

	req.Header.Add("Authorization", bearerToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err, "Failed to send GET request to argocd server", "url", getRepoUrl)
		return nil, err
	}

	return resp, nil
}

func postCreateRepository(log logr.Logger, host string, body CreateRepoBody, bearerToken string) (*http.Response, error) {
	jsonBody, err := common.ToJson(log, body)

	addRepoUrl := host + "/api/v1/repositories"
	req, err := http.NewRequest("POST", addRepoUrl, bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Error(err, "Failed to create POST request")
		return nil, err
	}

	req.Header.Add("Authorization", bearerToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err, "Failed to send POST request to argocd server", "url", addRepoUrl)
		return nil, err
	}

	return resp, nil
}
