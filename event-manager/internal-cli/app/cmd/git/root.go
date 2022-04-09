package git

import (
	"fmt"
	"github.com/compuzest/zlifecycle-internal-cli/app/common"
	"github.com/compuzest/zlifecycle-internal-cli/app/env"
	"github.com/spf13/cobra"
)

// RootCmd represents the validate command
var RootCmd = &cobra.Command{
	Use:     "git {command}",
	Example: "git clone -h",
	Short:   "git command offers subcommands for cloning zlifecycle repos",
	Long:    "git command offers subcommands for cloning zlifecycle repos",
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&env.GitAuth, "git-auth", "g", "", "Git auth method")
	if err := RootCmd.MarkPersistentFlagRequired("git-auth"); err != nil {
		fmt.Println(err)
		common.Failure(3201)
	}

	RootCmd.PersistentFlags().StringVarP(&env.GitHubAppID, "app-id", "a", "", "GitHub App organization ID")

	RootCmd.PersistentFlags().StringVarP(&env.GitHubAppSSHPath, "ssh", "s", "", "GitHub App private key filepath")

	RootCmd.PersistentFlags().StringVarP(&env.GitToken, "token", "t", "", "Git token")

	RootCmd.AddCommand(cloneCmd)
	RootCmd.AddCommand(loginCmd)
}
