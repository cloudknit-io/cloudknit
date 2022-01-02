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

// environmentComponentStatePullCmd represents the validate command
var environmentComponentStatePullCmd = &cobra.Command{
	Use:     "pull [flags]",
	Example: "pull --company dev --team checkout --environment test --component networking",
	Args:    cobra.NoArgs,
	Short:   "pull command pulls the environment component state and prints to stdout",
	Long:    `pull command pulls the environment component state from remote backend using zLifecycle State Manager and prints it to stdout`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		c := zlstate.NewHTTPStateManager(ctx, log.NewLogger().WithContext(ctx))
		req := zlstate.FetchZLStateComponentRequest{
			Company:     company,
			Team:        team,
			Environment: environment,
			Component:   component,
		}
		componentState, err := c.GetComponent(&req)
		if err != nil {
			return errors.Wrap(err, "error fetching environment component zLstate")
		}
		// print output
		json, err := common.ToJSON(componentState)
		if err != nil {
			return errors.Wrap(err, "error marshaling environment component zLstate")
		}

		fmt.Println(string(json))

		return nil
	},
}

func init() {
	environmentComponentStatePullCmd.Flags().StringVarP(&company, "company", "c", "", "Company name")
	if err := environmentComponentStatePullCmd.MarkFlagRequired("company"); err != nil {
		common.Failure(1201)
	}
	environmentComponentStatePullCmd.Flags().StringVarP(&team, "team", "t", "", "Team name")
	if err := environmentComponentStatePullCmd.MarkFlagRequired("team"); err != nil {
		common.Failure(1202)
	}
	environmentComponentStatePullCmd.Flags().StringVarP(&environment, "environment", "e", "", "Environment name")
	if err := environmentComponentStatePullCmd.MarkFlagRequired("environment"); err != nil {
		common.Failure(1203)
	}
	environmentComponentStatePullCmd.Flags().StringVarP(&component, "component", "m", "", "Environment Component name")
	if err := environmentComponentStatePullCmd.MarkFlagRequired("environment"); err != nil {
		common.Failure(1204)
	}
}
