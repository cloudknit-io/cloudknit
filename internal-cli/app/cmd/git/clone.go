package git

import (
	"context"

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

		gitOrg, _, err := util.ParseRepositoryInfo(repo)
		if err != nil {
			return errors.Wrap(err, "error extracting git organization from git repository url")
		}
		token, err := getGitToken(ctx, gitOrg, logger)
		if err != nil {
			return errors.Wrap(err, "error getting git token")
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

func init() {
	cloneCmd.Flags().StringVarP(&env.GitCloneDir, "dir", "d", "", "Directory in which to clone the repo")
}
