package state

import (
	"context"
	"fmt"
	"github.com/compuzest/zlifecycle-internal-cli/app/common"
	"github.com/compuzest/zlifecycle-internal-cli/app/log"
	"github.com/compuzest/zlifecycle-internal-cli/app/zlstate"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// stateCmd represents the validate command
var environmentStatePullCmd = &cobra.Command{
	Use:     "pull [flags]",
	Example: "pull --company dev --team checkout --environment test",
	Args:    cobra.NoArgs,
	Short:   "environment state pull command pulls the environment state and prints to stdout",
	Long:    `environment state pull command pulls the environment state from remote backend using zLifecycle State Manager and prints it to stdout`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		c := zlstate.NewHTTPStateManager(ctx, log.NewLogger().WithContext(ctx))
		req := zlstate.FetchZLStateRequest{
			Company:     company,
			Team:        team,
			Environment: environment,
		}
		zlstate, err := c.GetState(&req)
		if err != nil {
			return errors.Wrap(err, "error fetching environment zLstate")
		}
		// print output
		json, err := common.ToJSON(zlstate)
		if err != nil {
			return errors.Wrap(err, "error marshaling environment zLstate")
		}

		fmt.Println(string(json))

		return nil
	},
}

func init() {
	environmentStatePullCmd.Flags().StringVarP(&company, "company", "c", "", "Company name")
	if err := environmentStatePullCmd.MarkFlagRequired("company"); err != nil {
		common.Failure(1101)
	}
	environmentStatePullCmd.Flags().StringVarP(&team, "team", "t", "", "Team name")
	if err := environmentStatePullCmd.MarkFlagRequired("team"); err != nil {
		common.Failure(1102)
	}
	environmentStatePullCmd.Flags().StringVarP(&environment, "environment", "e", "", "Environment name")
	if err := environmentStatePullCmd.MarkFlagRequired("environment"); err != nil {
		common.Failure(1103)
	}
}
