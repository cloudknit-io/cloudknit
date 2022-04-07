package git

import (
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
	RootCmd.AddCommand(cloneCmd)
}
