package git

import (
	"context"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/compuzest/zlifecycle-internal-cli/app/api/git"
	"github.com/compuzest/zlifecycle-internal-cli/app/env"
	"github.com/compuzest/zlifecycle-internal-cli/app/log"
	"github.com/compuzest/zlifecycle-internal-cli/app/util"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// cloneCmd represents git clone command.
var cloneCmd = &cobra.Command{
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

func getGitAuth(ctx context.Context, repo string, logger *logrus.Entry) (*git.GoGitOptions, error) {
	opts := git.GoGitOptions{}
	if env.GitAuth == authModeSSH {
		pk, err := os.ReadFile(env.GitSSHPath)
		if err != nil {
			return nil, errors.Wrapf(err, "error reading private key at %s", env.GitAuth)
		}
		opts.Mode = git.AuthModeSSH
		opts.PrivateKey = pk
		return &opts, nil
	}

	gitOrg, _, err := util.ParseRepositoryInfo(repo)
	if err != nil {
		return nil, errors.Wrap(err, "error extracting git organization from git repository url")
	}
	token, err := getGitToken(ctx, gitOrg, logger)
	if err != nil {
		return nil, errors.Wrap(err, "error getting git token")
	}

	opts.Mode = git.AuthModeToken
	opts.Token = token

	return &opts, nil
}

func init() {
	cloneCmd.Flags().StringVarP(&env.GitCloneDir, "dir", "d", env.GitCloneDir, "Directory in which to clone the repo")
}
