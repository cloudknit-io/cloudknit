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
	"strings"

	"github.com/go-logr/logr"

	env "github.com/compuzest/zlifecycle-il-operator/controllers/util/env"
)

func RegisterRepo(log logr.Logger, api Api, repoOpts RepoOpts) (bool, error) {
	repoUri := repoOpts.RepoUrl[strings.LastIndex(repoOpts.RepoUrl, "/")+1:]
	repoName := strings.TrimSuffix(repoUri, ".git")
	tokenResponse, err := api.GetAuthToken()
	if err != nil {
		log.Error(err, "Error while calling ArgoCD API get auth token")
		return false, err
	}
	bearer := "Bearer " + tokenResponse.Token
	repositories, resp1, err1 := api.ListRepositories(bearer)
	if err1 != nil {
		return false, err1
	}
	defer resp1.Body.Close()
	log.Info("List of repositories registered on ArgoCD", "repos", repositories)
	if isRepoRegistered(*repositories, repoOpts.RepoUrl) {
		log.Info("Repository already registered on ArgoCD",
			"repos", repositories.Items,
			"repoName", repoName,
			"repoUrl", repoOpts.RepoUrl,
		)
		return false, nil
	}
	log.Info(
		"Repository is not registered on ArgoCD, registering now...",
		"repoName", repoName,
		"repoUrl", repoOpts.RepoUrl,
	)
	createRepoBody := CreateRepoBody{Repo: repoOpts.RepoUrl, Name: repoName, SshPrivateKey: repoOpts.SshPrivateKey}
	resp2, err2 := api.CreateRepository(createRepoBody, bearer)
	if err2 != nil {
		log.Error(err2, "Error while calling ArgoCD API post create repository")
		return false, err2
	}
	defer resp2.Body.Close()
	log.Info("Successfully registered repository on ArgoCD", "repo", repoOpts.RepoUrl)
	return true, nil
}

func TryCreateBootstrapApps(log logr.Logger) error {
	argocdAPI := NewHttpClient(log, env.Config.ArgocdServerUrl)

	tokenResponse, err := argocdAPI.GetAuthToken()
	if err != nil {
		log.Error(err, "Error while calling ArgoCD API get auth token")
		return err
	}
	bearer := "Bearer " + tokenResponse.Token

	exists, err := argocdAPI.DoesApplicationExist("company-bootstrap", bearer)
	if err != nil {
		log.Error(err, "Error while calling ArgoCD API to check if Application Exists")
		return err
	}
	if exists {
		log.Info("Application already registered on ArgoCD",
			"application", "company-bootstrap",
		)
	} else {
		companyResp, companyErr := argocdAPI.CreateApplication(GenerateCompanyBootstrapApp(), bearer)
		if companyErr != nil {
			log.Error(companyErr, "Error while creating Company Bootstrap Application")
			return companyErr
		}
		defer companyResp.Body.Close()
		log.Info("Successfully registered application on ArgoCD",
			"application", "company-bootstrap",
		)
	}

	exists2, err2 := argocdAPI.DoesApplicationExist("config-watcher-bootstrap", bearer)
	if err2 != nil {
		log.Error(err2, "Error while calling ArgoCD API to check if Application Exists")
		return err2
	}
	if exists2 {
		log.Info("Application already registered on ArgoCD",
			"application", "config-watcher-bootstrap",
		)
	} else {
		companyResp2, companyErr2 := argocdAPI.CreateApplication(GenerateConfigWatcherBootstrapApp(), bearer)
		if companyErr2 != nil {
			log.Error(companyErr2, "Error while creating Config Watcher Bootstrap Application")
			return companyErr2
		}
		defer companyResp2.Body.Close()
		log.Info("Successfully registered application on ArgoCD",
			"application", "config-watcher-bootstrap",
		)
	}

	return nil
}
