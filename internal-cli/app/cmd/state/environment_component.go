package state

import (
	"github.com/compuzest/zlifecycle-internal-cli/app/cmd/state/patch"
	"github.com/compuzest/zlifecycle-internal-cli/app/cmd/state/pull"
	"github.com/compuzest/zlifecycle-internal-cli/app/common"
	"github.com/spf13/cobra"
)

func NewEnvironmentComponentCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "component {command}",
		Example: "zl state environment pull -h",
		Args:    cobra.ExactArgs(1),
		Short:   "component command exposes subcommands for managing environment component state",
		Long: "component command exposes subcommands for managing environment component state" +
			" from remote backend using zLifecycle State Manager and prints it to stdout",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				if err := cmd.Help(); err != nil {
					common.Failure(1)
				}
				common.Success()
			}
		},
	}

	cmd.AddCommand(pull.NewEnvironmentComponentStatePullCmd())
	cmd.AddCommand(patch.NewEnvironmentComponentStatusPatchCmd())

	return cmd
}
