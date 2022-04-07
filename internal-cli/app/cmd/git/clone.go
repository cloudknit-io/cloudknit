package git

import (
	"context"
	"github.com/compuzest/zlifecycle-internal-cli/app/api/git"
	"github.com/compuzest/zlifecycle-internal-cli/app/api/github"
	"github.com/compuzest/zlifecycle-internal-cli/app/common"
	"github.com/compuzest/zlifecycle-internal-cli/app/env"
	"github.com/compuzest/zlifecycle-internal-cli/app/log"
	"github.com/compuzest/zlifecycle-internal-cli/app/util"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

const (
	authModeGitHubApp = "githubApp"
	authModeToken     = "token"
)

// cloneCmd represents git clone command
var cloneCmd = &cobra.Command{
	Use:     "clone {repository} [flags]",
	Example: "zl git clone https://github.com/some/repo -h",
	Args:    cobra.ExactArgs(1),
	Short:   "clone git repo",
	Long:    "clone git repo by selecting appropriate auth mode",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		log := log.NewLogger().WithContext(ctx)
		if len(args) != 1 {
			return errors.Errorf("invalid number of args (must be 1 - repository URL): %d", len(args))
		}
		repo := args[0]
		var token string
		if env.GitAuth == authModeGitHubApp {
			_token, err := getGitHubInstallationToken(ctx, repo, log)
			if err != nil {
				return errors.Wrap(err, "error generating GitHub App installation token")
			}
			token = _token
		} else {
			if env.GitToken == "" {
				return errors.New("git token is not passed for token auth mode via --token|-t flag")
			}
			token = env.GitToken
		}

		c, err := git.NewGoGit(ctx, &git.GoGitOptions{Mode: git.AuthModeToken, Token: token})
		if err != nil {
			return errors.Wrap(err, "error instantiating git client")
		}

		if err := c.Clone(repo, env.GitCloneDir); err != nil {
			return errors.Wrap(err, "error cloning repository")
		}

		return nil
	},
}

func getGitHubInstallationToken(ctx context.Context, repo string, log *logrus.Entry) (token string, err error) {
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
	c, err := github.NewClientBuilder().WithGitHubApp(ctx, dat, env.GitHubAppInstallationID).Build()
	if err != nil {
		return "", errors.Wrap(err, "error building GitHub App client")
	}

	owner, _, err := util.ParseRepositoryInfo(repo)
	if err != nil {
		return "", errors.Wrap(err, "error parsing repository info")
	}

	token, err = github.GenerateInstallationToken(log, c, owner)
	if err != nil {
		return "", errors.Wrap(err, "error generating GitHub App installation token")
	}

	return token, nil
}

func init() {
	cloneCmd.Flags().StringVarP(&env.GitAuth, "auth", "a", "", "Git auth method")
	if err := cloneCmd.MarkFlagRequired("auth"); err != nil {
		common.Failure(3201)
	}

	cloneCmd.Flags().StringVarP(&env.GitHubAppID, "app-id", "p", "", "GitHub App organization ID")

	cloneCmd.Flags().StringVarP(&env.GitHubAppInstallationID, "installation-id", "i", "", "GitHub App installation ID")

	cloneCmd.Flags().StringVarP(&env.GitHubAppSSHPath, "ssh", "s", "", "GitHub App private key filepath")

	cloneCmd.Flags().StringVarP(&env.GitToken, "token", "t", "", "Git token")

	cloneCmd.Flags().StringVarP(&env.GitCloneDir, "dir", "d", "", "Directory in which to clone the repo")
}
