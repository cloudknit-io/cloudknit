package patch

import (
	"context"

	"github.com/compuzest/zlifecycle-internal-cli/app/api/statemanager"
	"github.com/compuzest/zlifecycle-internal-cli/app/common"
	"github.com/compuzest/zlifecycle-internal-cli/app/env"
	"github.com/compuzest/zlifecycle-internal-cli/app/log"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func NewEnvironmentComponentStatusPatchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "patch [flags]",
		Example: "patch --company dev --team checkout --environment test --component networking --status provisioned",
		Args:    cobra.NoArgs,
		Short:   "patch command modifies the environment component state status and prints to stdout",
		Long:    `patch command modifies the environment component state status on remote backend using zLifecycle State Manager and prints it to stdout`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			logger := log.NewLogger().WithContext(ctx)
			c := statemanager.NewHTTPStateManager(ctx, logger)
			req := statemanager.UpdateZLStateComponentStatusRequest{
				Company:     env.Company,
				Team:        env.Team,
				Environment: env.Environment,
				Component:   env.Component,
				Status:      env.Status,
			}
			resp, err := c.PatchEnvironmentComponentStatus(&req)
			if err != nil {
				return errors.Wrap(err, "error patching environment component zLstate status")
			}
			// print output
			json, err := common.ToJSON(resp)
			if err != nil {
				return errors.Wrap(err, "error marshaling patch environment component status response")
			}

			logger.Info(string(json))

			return nil
		},
	}

	cmd.Flags().StringVarP(&env.Company, "company", "c", "", "Company name")
	if err := cmd.MarkFlagRequired("company"); err != nil {
		common.Failure(2201)
	}
	cmd.Flags().StringVarP(&env.Team, "team", "t", "", "Team name")
	if err := cmd.MarkFlagRequired("team"); err != nil {
		common.Failure(2202)
	}
	cmd.Flags().StringVarP(&env.Environment, "environment", "e", "", "Environment name")
	if err := cmd.MarkFlagRequired("environment"); err != nil {
		common.Failure(2203)
	}
	cmd.Flags().StringVarP(&env.Component, "component", "m", "", "Environment Component name")
	if err := cmd.MarkFlagRequired("environment"); err != nil {
		common.Failure(2204)
	}
	cmd.Flags().StringVarP(&env.Status, "status", "s", "", "Environment Component status")
	if err := cmd.MarkFlagRequired("status"); err != nil {
		common.Failure(2204)
	}

	return cmd
}
