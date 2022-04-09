package git

import (
	"context"
	"github.com/compuzest/zlifecycle-internal-cli/app/api/git"
	"github.com/compuzest/zlifecycle-internal-cli/app/env"
	"github.com/compuzest/zlifecycle-internal-cli/app/log"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

const (
	configFile = ".gitconfig"
)

// loginCmd represents git clone command
var loginCmd = &cobra.Command{
	Use:     "login {gitOrg} [flags]",
	Example: "zl git login zlifecycle-il",
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
		token, err := getGitToken(ctx, gitOrg, logger)
		if err != nil {
			return errors.Wrap(err, "error getting git token")
		}

		gitconfig := git.Config(token, env.GitBaseURL)

		path := filepath.Join(env.GitConfigDir, configFile)
		if err := os.WriteFile(path, []byte(gitconfig), 0444); err != nil {
			return errors.Wrapf(err, "error creating %s file", path)
		}

		return nil
	},
}

func init() {
	loginCmd.Flags().StringVarP(&env.GitBaseURL, "--base-url", "b", "", "Base git https url (ex. https://github.com)")
	loginCmd.Flags().StringVarP(&env.GitConfigDir, "--config-path", "c", "", "Path where to create .gitconfig")

}
