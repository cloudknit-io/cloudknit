package state

import (
	"github.com/compuzest/zlifecycle-internal-cli/app/cmd/state/pull"
	"github.com/compuzest/zlifecycle-internal-cli/app/common"
	"github.com/spf13/cobra"
)

func NewEnvironmentCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "environment {command}",
		Example: "zl state environment pull -h",
		Args:    cobra.ExactArgs(1),
		Short:   "environment command exposes subcommands for managing environment state",
		Long:    `environment command exposes subcommands for managing environment state from remote backend using zLifecycle State Manager and prints it to stdout`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				if err := cmd.Help(); err != nil {
					common.Failure(1)
				}
				common.Success()
			}
		},
	}

	cmd.AddCommand(pull.NewEnvironmentStatePullCmd())

	return cmd
}
