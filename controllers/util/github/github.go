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
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/common"
	"github.com/google/go-github/v42/github"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func GetAppInstallationID(
	log *logrus.Entry,
	api AppAPI,
	owner string,
	repo string,
) (*int64, error) {
	log.WithFields(logrus.Fields{
		"owner": owner,
		"repo":  repo,
	}).Info("Finding GitHub App installation ID")
	installation, resp, err := api.FindRepositoryInstallation(owner, repo)
	if err != nil {
		return nil, errors.Wrapf(err, "error finding repository installation ID for repo %s/%s", owner, repo)
	}
	defer common.CloseBody(resp.Body)

	if resp.StatusCode != 200 {
		return nil, errors.Errorf("find repository installation returned non-OK status: %d", resp.StatusCode)
	}

	log.WithFields(logrus.Fields{
		"owner":          owner,
		"repo":           repo,
		"installationId": installation.ID,
		"appId":          installation.AppID,
	}).Info("Found installation for repository")

	return installation.ID, nil
}

// TryCreateRepository tries to create a private repository in a organization (enter blank string for owner if it is a user repo)
// It returns true if repository is created, false if repository already exists, or any kind of error.
func TryCreateRepository(log *logrus.Entry, api RepositoryAPI, owner string, repo string) (bool, error) {
	log.WithFields(logrus.Fields{
		"owner": owner,
		"repo":  repo,
	}).Info("Checking does repo exist on GitHub")
	r, resp1, err := api.GetRepository(owner, repo)
	if err != nil {
		return false, errors.Wrapf(err, "error getting repository %s/%s", owner, repo)
	}
	defer common.CloseBody(resp1.Body)

	if r != nil {
		log.WithFields(logrus.Fields{
			"owner": owner,
			"repo":  repo,
		}).Info("GitHub repository already exists")
		return false, nil
	}

	log.WithFields(logrus.Fields{
		"owner": owner,
		"repo":  repo,
	}).Info("Creating new private repository in GitHub")

	nr, resp2, err := api.CreateRepository(owner, repo)
	if err != nil {
		return false, errors.Wrapf(err, "error creating private repository %s/%s", owner, repo)
	}

	defer common.CloseBody(resp2.Body)

	return nr != nil, nil
}

// CreateRepoWebhook tries to create a repository webhook
// It returns true if webhook already exists, false if new webhook was created, or any kind of error.
func CreateRepoWebhook(log *logrus.Entry, api RepositoryAPI, repoURL string, payloadURL string, webhookSecret string) (bool, error) {
	owner, repo, err := common.ParseRepositoryInfo(repoURL)
	if err != nil {
		return false, errors.Wrap(err, "error parsing owner and repo from repo url")
	}

	hooks, resp1, err := api.ListHooks(owner, repo, nil)
	if err != nil {
		return false, errors.Wrapf(err, "error listing webhooks for repo %s/%s", owner, repo)
	}
	defer common.CloseBody(resp1.Body)

	exists, err := checkIsHookRegistered(
		hooks,
		payloadURL,
	)
	if exists {
		log.WithFields(logrus.Fields{
			"url":        repoURL,
			"owner":      owner,
			"repo":       repo,
			"payloadUrl": payloadURL,
		}).Info("Hook already exists")
		return true, nil
	}
	if err != nil {
		return false, errors.Wrapf(err, "error checking is hook already registered for repo %s/%s", owner, repo)
	}

	h := newHook(payloadURL, webhookSecret)
	hook, resp2, err := api.CreateHook(owner, repo, &h)
	if err != nil {
		return false, errors.Wrapf(err, "error creating repository hook for repo %s/%s", owner, repo)
	}
	defer common.CloseBody(resp2.Body)

	log.WithFields(logrus.Fields{
		"url":        repoURL,
		"owner":      owner,
		"repo":       repo,
		"hookId":     *hook.ID,
		"payloadUrl": *hook.URL,
	}).Info("Successfully created repository webhook")

	return false, nil
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
