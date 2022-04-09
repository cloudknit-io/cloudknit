package pull

import (
	"context"
	"fmt"

	"github.com/compuzest/zlifecycle-internal-cli/app/api/statemanager"
	"github.com/compuzest/zlifecycle-internal-cli/app/common"
	"github.com/compuzest/zlifecycle-internal-cli/app/env"
	"github.com/compuzest/zlifecycle-internal-cli/app/log"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// stateCmd represents the validate command
var EnvironmentStatePullCmd = &cobra.Command{
	Use:     "pull [flags]",
	Example: "pull --company dev --team checkout --environment test",
	Args:    cobra.NoArgs,
	Short:   "environment state pull command pulls the environment state and prints to stdout",
	Long:    `environment state pull command pulls the environment state from remote backend using zLifecycle State Manager and prints it to stdout`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		c := statemanager.NewHTTPStateManager(ctx, log.NewLogger().WithContext(ctx))
		req := statemanager.FetchZLStateRequest{
			Company:     env.Company,
			Team:        env.Team,
			Environment: env.Environment,
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
	EnvironmentStatePullCmd.Flags().StringVarP(&env.Company, "company", "c", "", "Company name")
	if err := EnvironmentStatePullCmd.MarkFlagRequired("company"); err != nil {
		common.Failure(1101)
	}
	EnvironmentStatePullCmd.Flags().StringVarP(&env.Team, "team", "t", "", "Team name")
	if err := EnvironmentStatePullCmd.MarkFlagRequired("team"); err != nil {
		common.Failure(1102)
	}
	EnvironmentStatePullCmd.Flags().StringVarP(&env.Environment, "environment", "e", "", "Environment name")
	if err := EnvironmentStatePullCmd.MarkFlagRequired("environment"); err != nil {
		common.Failure(1103)
	}
}
