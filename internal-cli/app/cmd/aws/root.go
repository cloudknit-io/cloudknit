package aws

import (
	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "aws {command}",
		Example: "aws -h",
		Short:   "aws command offers subcommands for AWS actions",
	}

	cmd.AddCommand(NewConfigureCmd())

	return cmd
}
