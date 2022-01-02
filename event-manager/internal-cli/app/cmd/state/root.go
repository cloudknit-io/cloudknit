package state

import (
	"github.com/compuzest/zlifecycle-internal-cli/app/env"
	"github.com/spf13/cobra"
)

// RootCmd represents the validate command
var RootCmd = &cobra.Command{
	Use:     "state {command}",
	Example: "state environment -h",
	Short:   "state command offers subcomands for managing environment or component state",
	Long:    `state command offers subcomands for managing environment or component state on remote backend using zLifecycle State Manager`,
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&env.StateManagerURL, "url", "u", "", "zLifecycle State Manager URL")

	RootCmd.AddCommand(environmentCmd)
	RootCmd.AddCommand(environmentComponentCmd)
}
