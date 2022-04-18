package git

import (
	"context"
	"os"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/compuzest/zlifecycle-internal-cli/app/api/git"
	"github.com/compuzest/zlifecycle-internal-cli/app/env"
	"github.com/compuzest/zlifecycle-internal-cli/app/log"
	"github.com/compuzest/zlifecycle-internal-cli/app/util"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func NewCloneCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "clone {repository} [flags]",
		Example: "zl git clone https://github.com/some/repo --git-auth githubApp --app-id 172698 --ssh /path/to/githubapp/private_key.pem",
		Args:    cobra.ExactArgs(1),
		Short:   "clone git repo",
		Long:    "clone git repo using configurable git auth modes",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			logger := log.NewLogger().WithContext(ctx)
			if len(args) != 1 {
				return errors.Errorf("invalid number of args (must be 1 - repository URL): %d", len(args))
			}
			repo := args[0]
			repo = strings.TrimPrefix(repo, "git::")

			logger.Infof("Cloning git repo %s using auth mode %s", repo, env.GitAuth)

			auth, err := getGitAuth(ctx, repo, logger)
			if err != nil {
				return errors.Wrap(err, "error getting git auth")
			}

			c, err := git.NewGoGit(ctx, auth)
			if err != nil {
				return errors.Wrap(err, "error instantiating git client")
			}

			if err := c.Clone(repo, env.GitCloneDir); err != nil {
				return errors.Wrap(err, "error cloning repository")
			}

			logger.Infof("Successfully cloned repo %s in directory %s", repo, env.GitCloneDir)

			return nil
		},
	}

	cmd.Flags().StringVarP(&env.GitCloneDir, "dir", "d", env.GitCloneDir, "Directory in which to clone the repo")

	return cmd
}

func getGitAuth(ctx context.Context, repo string, logger *logrus.Entry) (*git.GoGitOptions, error) {
	if env.GitAuth == authModeSSH {
		logger.Infof("Using SSH auth mode for git")
		return getGitSSHAuth(env.GitSSHPath)
	}
	logger.Infof("Using token auth mode for git")
	return getGitTokenAuth(ctx, repo, logger)
}

func getGitSSHAuth(sshPath string) (*git.GoGitOptions, error) {
	pk, err := os.ReadFile(sshPath)
	if err != nil {
		return nil, errors.Wrapf(err, "error reading private key at %s", env.GitAuth)
	}

	return &git.GoGitOptions{Mode: git.AuthModeSSH, PrivateKey: pk}, nil
}

func getGitTokenAuth(ctx context.Context, repoURL string, logger *logrus.Entry) (*git.GoGitOptions, error) {
	gitOrg, _, err := util.ReverseParseGitURL(repoURL)
	if err != nil {
		return nil, errors.Wrapf(err, "error parsing git url [%s]", repoURL)
	}
	if gitOrg == "" {
		return nil, errors.Errorf("error parsing repo url [%s]: error parsing git organization", repoURL)
	}
	token, err := getGitToken(ctx, gitOrg, logger)
	if err != nil {
		return nil, errors.Wrap(err, "error getting git token")
	}

	return &git.GoGitOptions{Mode: git.AuthModeToken, Token: token}, nil
}
