package git

import (
	"context"
	"os"
	"path/filepath"

	"github.com/compuzest/zlifecycle-internal-cli/app/api/git"
	"github.com/compuzest/zlifecycle-internal-cli/app/env"
	"github.com/compuzest/zlifecycle-internal-cli/app/log"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const (
	configFile = ".gitconfig"
)

// NewLoginCmd created the cobra command for git clone operation.
func NewLoginCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "login {gitOrg} [flags]",
		Example: "zl git login zlifecycle-il --git-auth githubApp --app-id 172698 --ssh /path/to/githubapp/private_key.pem",
		Args:    cobra.ExactArgs(1),
		Short:   "login to git",
		Long:    "login to git by creating a .gitconfig file in home directory and replacing github https urls with token",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			logger := log.NewLogger().WithContext(ctx)
			if len(args) != 1 {
				return errors.Errorf("invalid number of args (must be 1 - repository URL): %d", len(args))
			}
			gitOrg := args[0]

			logger.Infof("Logging in to git provider %s using auth mode %s", env.GitBaseURL, env.GitAuth)

			token, err := getGitToken(ctx, gitOrg, logger)
			if err != nil {
				return errors.Wrap(err, "error getting git token")
			}

			gitconfig := git.Config(token, env.GitBaseURL)

			path := filepath.Join(env.GitConfigDir, configFile)
			if err := os.WriteFile(path, []byte(gitconfig), 0o444); err != nil {
				return errors.Wrapf(err, "error creating %s file", path)
			}

			logger.Infof("Successfully wrote .gitconfig file at %s", path)
			logger.Infof(".gitfile content:\n%s", gitconfig)
			return nil
		},
	}

	cmd.Flags().StringVarP(&env.GitBaseURL, "git-provider", "b", env.GitBaseURL, "Base git https url (ex. https://github.com)")
	cmd.Flags().StringVarP(&env.GitConfigDir, "git-config-path", "c", env.GitConfigDir, "Path where to create .gitconfig")

	return cmd
}
