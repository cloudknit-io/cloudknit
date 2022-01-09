package pull

import (
	"context"
	"fmt"
	"github.com/compuzest/zlifecycle-internal-cli/app/common"
	"github.com/compuzest/zlifecycle-internal-cli/app/env"
	"github.com/compuzest/zlifecycle-internal-cli/app/log"
	"github.com/compuzest/zlifecycle-internal-cli/app/zlstate"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"os"
)

// EnvironmentComponentStatePullCmd represents the validate command
var EnvironmentComponentStatePullCmd = &cobra.Command{
	Use:     "pull [flags]",
	Example: "pull --company dev --team checkout --environment test --component networking",
	Args:    cobra.NoArgs,
	Short:   "pull command pulls the environment component state and prints to stdout",
	Long:    `pull command pulls the environment component state from remote backend using zLifecycle State Manager and prints it to stdout`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		c := zlstate.NewHTTPStateManager(ctx, log.NewLogger().WithContext(ctx))
		req := zlstate.FetchZLStateComponentRequest{
			Company:     env.Company,
			Team:        env.Team,
			Environment: env.Environment,
			Component:   env.Component,
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

		os.Exit(getStatusCode(componentState.Component.Status))
		return nil
	},
}

func getStatusCode(status string) int {
	switch status {
	case "not_provisioned":
		return 1
	case "waiting_for_approval":
		return 2
	case "provisioning":
		return 3
	case "provisioned":
		return 4
	case "destroying":
		return 5
	case "destroyed":
		return 6
	default:
		return 99
	}
}

func init() {
	EnvironmentComponentStatePullCmd.Flags().StringVarP(&env.Company, "company", "c", "", "Company name")
	if err := EnvironmentComponentStatePullCmd.MarkFlagRequired("company"); err != nil {
		common.Failure(1201)
	}
	EnvironmentComponentStatePullCmd.Flags().StringVarP(&env.Team, "team", "t", "", "Team name")
	if err := EnvironmentComponentStatePullCmd.MarkFlagRequired("team"); err != nil {
		common.Failure(1202)
	}
	EnvironmentComponentStatePullCmd.Flags().StringVarP(&env.Environment, "environment", "e", "", "Environment name")
	if err := EnvironmentComponentStatePullCmd.MarkFlagRequired("environment"); err != nil {
		common.Failure(1203)
	}
	EnvironmentComponentStatePullCmd.Flags().StringVarP(&env.Component, "component", "m", "", "Environment Component name")
	if err := EnvironmentComponentStatePullCmd.MarkFlagRequired("environment"); err != nil {
		common.Failure(1204)
	}
}
