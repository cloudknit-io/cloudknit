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
	authModeGitHubApp = "githubApp"
	authModeToken     = "token"
)

func getGitToken(ctx context.Context, repo string, logger *logrus.Entry) (token string, err error) {
	switch env.GitAuth {
	case authModeGitHubApp:
		token, err = getGitHubInstallationToken(ctx, repo, logger)
		if err != nil {
			return "", errors.Wrap(err, "error generating GitHub App installation token")
		}
		return token, nil
	case authModeToken:
		if env.GitToken == "" {
			return "", errors.New("git token is not passed for token auth mode via --token|-t flag")
		}
		return env.GitToken, nil
	default:
		msg := "Invalid git auth mode: %s"
		logger.Errorf(msg, env.GitAuth)
		return "", errors.Errorf(msg, env.GitAuth)
	}
}

func getGitHubInstallationToken(ctx context.Context, gitOrg string, logger *logrus.Entry) (token string, err error) {
	if env.GitHubAppID == "" {
		return "", errors.New("GitHub App ID not provided with --app-id|-p flag")
	}
	if env.GitHubAppSSHPath == "" {
		return "", errors.New("path to GitHub App private key file is not provided")
	}
	dat, err := os.ReadFile(env.GitHubAppSSHPath)
	if err != nil {
		return "", errors.Wrap(err, "error reading GitHub App private key file")
	}
	c, err := github.NewClientBuilder().WithGitHubApp(ctx, dat, env.GitHubAppID).Build()
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
