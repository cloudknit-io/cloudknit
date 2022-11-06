package git

import (
	"fmt"

	"github.com/compuzest/zlifecycle-internal-cli/app/common"
	"github.com/compuzest/zlifecycle-internal-cli/app/env"
	"github.com/spf13/cobra"
)

// NewRootCmd creates the root command for git operations.
func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "git {command}",
		Example: "git clone -h",
		Short:   "git command offers subcommands for cloning zlifecycle repos",
	}

	cmd.PersistentFlags().StringVarP(&env.GitAuth, "git-auth", "g", "", "Git auth method")
	if err := cmd.MarkPersistentFlagRequired("git-auth"); err != nil {
		fmt.Println(err)
		common.Failure(3201)
	}
	cmd.PersistentFlags().StringVarP(&env.GitSSHPath, "git-ssh", "s", env.GitSSHPath, "Git private key filepath")
	cmd.PersistentFlags().StringVarP(&env.GitHubAppID, "github-app-id", "a", "", "GitHub App organization ID")
	cmd.PersistentFlags().StringVarP(&env.GitToken, "token", "t", "", "Git token")

	cmd.AddCommand(NewCloneCmd())
	cmd.AddCommand(NewLoginCmd())

	return cmd
}
