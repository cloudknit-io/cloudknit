package pull

import (
	"context"
	"os"

	"github.com/compuzest/zlifecycle-internal-cli/app/api/statemanager"
	"github.com/compuzest/zlifecycle-internal-cli/app/common"
	"github.com/compuzest/zlifecycle-internal-cli/app/env"
	"github.com/compuzest/zlifecycle-internal-cli/app/log"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func NewEnvironmentComponentStatePullCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "pull [flags]",
		Example: "pull --company dev --team checkout --environment test --component networking",
		Args:    cobra.NoArgs,
		Short:   "pull command pulls the environment component state and prints to stdout",
		Long:    `pull command pulls the environment component state from remote backend using zLifecycle State Manager and prints it to stdout`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			logger := log.NewLogger().WithContext(ctx)
			c := statemanager.NewHTTPStateManager(ctx, logger)
			req := statemanager.FetchZLStateComponentRequest{
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

			logger.Info(string(json))

			os.Exit(getExitCode(componentState.Component.Status))
			return nil
		},
	}

	cmd.Flags().StringVarP(&env.Company, "company", "c", env.Company, "Company name")
	if err := cmd.MarkFlagRequired("company"); err != nil {
		common.Failure(1201)
	}
	cmd.Flags().StringVarP(&env.Team, "team", "t", env.Team, "Team name")
	if err := cmd.MarkFlagRequired("team"); err != nil {
		common.Failure(1202)
	}
	cmd.Flags().StringVarP(&env.Environment, "environment", "e", env.Environment, "Environment name")
	if err := cmd.MarkFlagRequired("environment"); err != nil {
		common.Failure(1203)
	}
	cmd.Flags().StringVarP(&env.Component, "component", "m", env.Component, "Environment Component name")
	if err := cmd.MarkFlagRequired("environment"); err != nil {
		common.Failure(1204)
	}

	return cmd
}

func getExitCode(status string) int {
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
