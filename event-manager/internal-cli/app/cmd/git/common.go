package git

import (
	"context"
	"os"

	"github.com/compuzest/zlifecycle-internal-cli/app/api/github"
	"github.com/compuzest/zlifecycle-internal-cli/app/env"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	authModeGitHubApp         = "github-app"
	authModeGitHubAppInternal = "github-app-internal"
	authModeGitHubAppPublic   = "github-app-public"
	authModeToken             = "token"
	errMsgInstallationToken   = "error generating GitHub App installation token"
)

func getGitToken(ctx context.Context, repo string, logger *logrus.Entry) (token string, err error) {

	switch env.GitAuth {
	case authModeGitHubApp:
		token, err = getGitHubInstallationToken(ctx, authModeGitHubApp, repo, logger)
		if err != nil {
			return "", errors.Wrap(err, errMsgInstallationToken)
		}
		return token, nil
	case authModeGitHubAppPublic:
		token, err = getGitHubInstallationToken(ctx, authModeGitHubAppPublic, repo, logger)
		if err != nil {
			return "", errors.Wrap(err, errMsgInstallationToken)
		}
		return token, nil
	case authModeGitHubAppInternal:
		token, err = getGitHubInstallationToken(ctx, authModeGitHubAppInternal, repo, logger)
		if err != nil {
			return "", errors.Wrap(err, errMsgInstallationToken)
		}
		return token, nil
	case authModeToken:
		if env.GitToken == "" {
			return "", errors.New("missing git token for token auth mode")
		}
		return env.GitToken, nil
	default:
		msg := "Invalid git auth mode: %s"
		logger.Errorf(msg, env.GitAuth)
		return "", errors.Errorf(msg, env.GitAuth)
	}
}

func getGitHubInstallationToken(ctx context.Context, authMode, gitOrg string, logger *logrus.Entry) (token string, err error) {
	appID, err := getGitHubAppID(authMode)
	if err != nil {
		return "", errors.Wrap(err, "error getting github app id")
	}
	if env.GitHubAppSSHPath == "" {
		return "", errors.New("path to GitHub App private key file is not provided")
	}
	dat, err := os.ReadFile(env.GitHubAppSSHPath)
	if err != nil {
		return "", errors.Wrap(err, "error reading GitHub App private key file")
	}
	c, err := github.NewClientBuilder().WithGitHubApp(ctx, dat, appID).Build()
	if err != nil {
		return "", errors.Wrap(err, "error building GitHub App client")
	}

	token, err = github.GenerateInstallationToken(logger, c, gitOrg)
	if err != nil {
		return "", errors.Wrap(err, "error generating GitHub App installation token")
	}

	logger.Infof("Generated git token %s", token)

	return token, nil
}

func getGitHubAppID(authMode string) (token string, err error) {
	switch authMode {
	case authModeGitHubApp:
		if env.GitHubAppID == "" {
			return "", errors.New("app ID not provided for github app auth mode")
		}
		token = env.GitHubAppID
	case authModeGitHubAppInternal:
		token = env.GitHubAppIDInternal
	case authModeGitHubAppPublic:
		token = env.GitHubAppIDPublic
	default:
		err = errors.Errorf("auth mode does not support github: %s", authMode)
	}
	return
}
