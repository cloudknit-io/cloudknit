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
	"github.com/compuzest/zlifecycle-internal-cli/app/util"
	"strconv"

	"github.com/google/go-github/v42/github"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func GetAppInstallationID(
	log *logrus.Entry,
	client API,
	org string,
) (installationID *int64, appID *int64, err error) {
	log.WithFields(logrus.Fields{
		"org": org,
	}).Infof("Finding GitHub App installation ID for organization %s", org)
	installation, resp, err := client.FindOrganizationInstallation(org)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "error finding repository installation ID for org %s", org)
	}
	defer util.CloseBody(resp.Body)

	if resp.StatusCode != 200 {
		return nil, nil, errors.Errorf("find repository installation returned non-OK status: %d", resp.StatusCode)
	}

	log.WithFields(logrus.Fields{
		"org":            org,
		"installationId": strconv.FormatInt(*installation.ID, 10),
		"appId":          strconv.FormatInt(*installation.AppID, 10),
	}).Infof("Found installation ID for organization %s", org)

	return installation.ID, installation.AppID, nil
}

func GenerateInstallationToken(log *logrus.Entry, client API, org string) (token string, err error) {
	installationID, _, err := GetAppInstallationID(log, client, org)
	if err != nil {
		return "", errors.Wrapf(err, "error getting installation id for org [%s]", org)
	}
	log.WithFields(logrus.Fields{
		"org": org,
	}).Infof("Creating installation token for organization %s", org)
	installationToken, resp, err := client.CreateInstallationToken(*installationID)
	if err != nil {
		return "", errors.Wrapf(
			err,
			"error creating installation token for org [%s] and installation ID [%s]",
			strconv.FormatInt(*installationID, 10), org,
		)
	}
	defer util.CloseBody(resp.Body)

	if resp.StatusCode != 201 {
		return "", errors.Errorf("create installation token returned non-Created status: %d", resp.StatusCode)
	}

	return installationToken.GetToken(), nil
}

// CreateRepository tries to create a private repository in a organization (enter blank string for owner if it is a user repo)
// It returns true if repository is created, false if repository already exists, or any kind of error.
func CreateRepository(log *logrus.Entry, api RepositoryAPI, owner string, repo string) (exists bool, err error) {
	log.WithFields(logrus.Fields{
		"owner": owner,
		"repo":  repo,
	}).Infof("Checking does repository %s/%s exist on GitHub", owner, repo)
	r, resp1, err := api.GetRepository(owner, repo)
	if err != nil {
		return false, errors.Wrapf(err, "error getting repository %s/%s", owner, repo)
	}
	defer util.CloseBody(resp1.Body)

	if r != nil {
		log.WithFields(logrus.Fields{
			"owner": owner,
			"repo":  repo,
		}).Infof("GitHub repository %s/%s already exists", owner, repo)
		return false, nil
	}

	log.WithFields(logrus.Fields{
		"owner": owner,
		"repo":  repo,
	}).Infof("Creating new private repository %s/%s in GitHub", owner, repo)

	nr, resp2, err := api.CreateRepository(owner, repo)
	if err != nil {
		return false, errors.Wrapf(err, "error creating private repository %s/%s", owner, repo)
	}

	defer util.CloseBody(resp2.Body)

	return nr != nil, nil
}

// CreateRepoWebhook tries to create a repository webhook
// It returns true if webhook already exists, false if new webhook was created, or any kind of error.
func CreateRepoWebhook(log *logrus.Entry, api API, repoURL string, payloadURL string, webhookSecret string) (bool, error) {
	owner, repo, err := util.ParseRepositoryInfo(repoURL)
	if err != nil {
		return false, errors.Wrap(err, "error parsing owner and repo from repo url")
	}

	hooks, resp1, err := api.ListHooks(owner, repo, nil)
	if err != nil {
		return false, errors.Wrapf(err, "error listing webhooks for repo %s/%s", owner, repo)
	}
	defer util.CloseBody(resp1.Body)

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
		}).Infof("Webhook already exists for repository %s", repoURL)
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
	defer util.CloseBody(resp2.Body)

	log.WithFields(logrus.Fields{
		"url":        repoURL,
		"owner":      owner,
		"repo":       repo,
		"hookId":     *hook.ID,
		"payloadUrl": *hook.URL,
	}).Infof("Successfully created webhook for repository %s", repoURL)

	return false, nil
}

func checkIsHookRegistered(hooks []*github.Hook, payloadURL string) (bool, error) {
	for _, h := range hooks {
		cfg := new(HookCfg)
		err := util.FromJSONMap(h.Config, cfg)
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
