package pull

import (
	"context"
	"github.com/compuzest/zlifecycle-internal-cli/app/api/statemanager"
	"github.com/compuzest/zlifecycle-internal-cli/app/common"
	"github.com/compuzest/zlifecycle-internal-cli/app/env"
	"github.com/compuzest/zlifecycle-internal-cli/app/log"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func NewEnvironmentStatePullCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "pull [flags]",
		Example: "pull --company dev --team checkout --environment test",
		Args:    cobra.NoArgs,
		Short:   "environment state pull command pulls the environment state and prints to stdout",
		Long:    `environment state pull command pulls the environment state from remote backend using zLifecycle State Manager and prints it to stdout`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			logger := log.Logger.WithContext(ctx)
			c := statemanager.NewHTTPStateManager(ctx, logger)
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

			logger.Info(string(json))

			return nil
		},
	}

	cmd.Flags().StringVarP(&env.Company, "company", "c", "", "Company name")
	if err := cmd.MarkFlagRequired("company"); err != nil {
		common.Failure(1101)
	}
	cmd.Flags().StringVarP(&env.Team, "team", "t", "", "Team name")
	if err := cmd.MarkFlagRequired("team"); err != nil {
		common.Failure(1102)
	}
	cmd.Flags().StringVarP(&env.Environment, "environment", "e", "", "Environment name")
	if err := cmd.MarkFlagRequired("environment"); err != nil {
		common.Failure(1103)
	}

	return cmd
}
