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

package github

import (
	"errors"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/common"
	"github.com/go-logr/logr"
	"github.com/google/go-github/v32/github"
	"strings"
)

// TryCreateRepository tries to create a private repository in a organization (enter blank string for owner if it is a user repo)
// It returns true if repository is created, false if repository already exists, or any kind of error.
func TryCreateRepository(log logr.Logger, api RepositoryAPI, owner string, repo string) (bool, error) {
	log.Info("Checking does repo exist on GitHub", "owner", owner, "repo", repo)
	r, resp1, err := api.GetRepository(owner, repo)
	if err != nil {
		log.Error(err, "Error while fetching repository from GitHub API",
			"owner", owner,
			"repo", repo,
		)
		return false, err
	}
	defer common.CloseBody(resp1.Body)

	log.Info("Call to GitHub API get repo succeeded", "code", resp1.StatusCode)

	if r != nil {
		log.Info("GitHub repository already exists", "owner", owner, "repo", repo)
		return false, nil
	}

	log.Info("Creating new private repository in GitHub", "owner", owner, "repo", repo)

	nr, resp2, err := api.CreateRepository(owner, repo)
	if err != nil {
		log.Error(err, "Error while creating private repository on GitHub",
			"owner", owner,
			"repo", repo,
		)
		return false, err
	}

	log.Info("Call to GitHub API create repo succeeded", "code", resp2.StatusCode)

	defer common.CloseBody(resp2.Body)

	return nr != nil, nil
}

// CreateRepoWebhook tries to create a repository webhook
// It returns true if webhook already exists, false if new webhook was created, or any kind of error.
func CreateRepoWebhook(log logr.Logger, api RepositoryAPI, repoURL string, payloadURL string, webhookSecret string) (bool, error) {
	owner, repo, err := parseRepoURL(repoURL)
	if err != nil {
		log.Error(err, "Error while parsing owner and repo name from repo url", "url", repoURL)
		return false, err
	}
	log.Info("Parsed repo url", "url", repoURL, "owner", owner, "repo", repo)

	hooks, resp1, err := api.ListHooks(owner, repo, nil)
	if err != nil {
		log.Error(err, "Error while fetching list of webhooks",
			"url", repoURL,
			"owner", owner,
			"repo", repo,
		)
		return false, err
	}
	defer common.CloseBody(resp1.Body)

	exists, err := checkIsHookRegistered(
		hooks,
		payloadURL,
	)
	if exists {
		log.Info(
			"Hook already exists",
			"url", repoURL,
			"owner", owner,
			"repo", repo,
			"payloadUrl", payloadURL,
		)
		return true, nil
	}
	if err != nil {
		log.Error(err, "Error while checking is hook already registered", "hooks", hooks)
		return false, err
	}

	h := newHook(payloadURL, webhookSecret)
	hook, resp2, err := api.CreateHook(owner, repo, &h)
	if err != nil {
		log.Error(
			err, "Error while calling create repository hook",
			"url", repoURL,
			"owner", owner,
			"repo", repo,
		)
		return false, err
	}
	defer common.CloseBody(resp2.Body)

	log.Info(
		"Successfully created repository webhook",
		"url", repoURL,
		"owner", owner,
		"repo", repo,
		"hookId", *hook.ID,
		"payloadUrl", *hook.URL,
	)
	return false, nil
}

func parseRepoURL(url string) (Owner, Repo, error) {
	if url == "" {
		return "", "", errors.New("URL cannot be empty")
	}

	owner := url[strings.LastIndex(url, ":")+1 : strings.LastIndex(url, "/")]
	repoURI := url[strings.LastIndex(url, "/")+1:]
	repo := strings.TrimSuffix(repoURI, ".git")

	return owner, repo, nil
}

func checkIsHookRegistered(hooks []*github.Hook, payloadURL string) (bool, error) {
	for _, h := range hooks {
		cfg := new(HookCfg)
		err := common.FromJSONMap(h.Config, cfg)
		if err != nil {
			return false, err
		}
		if cfg.URL == payloadURL {
			return true, nil
		}
	}

	return false, nil
}

func newHook(payloadURL string, secret string) github.Hook {
	events := []string{"push"}
	cfg := map[string]interface{}{
		"url":          payloadURL,
		"content_type": "json",
		"secret":       secret,
	}
	return github.Hook{Events: events, Active: github.Bool(true), Config: cfg}
}
