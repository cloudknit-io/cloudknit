package aws

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"strings"

	"github.com/compuzest/zlifecycle-internal-cli/app/api/aws"
	secrets "github.com/compuzest/zlifecycle-internal-cli/app/api/ssm"
	"github.com/compuzest/zlifecycle-internal-cli/app/lib/secret"
	"github.com/pkg/errors"

	"github.com/compuzest/zlifecycle-internal-cli/app/common"
	"github.com/compuzest/zlifecycle-internal-cli/app/env"
	"github.com/compuzest/zlifecycle-internal-cli/app/log"
	"github.com/spf13/cobra"
)

func NewConfigureCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "configure [flags]",
		Example: "configure --auth-mode profile --profile compuzest-shared --generated-profile customer-state --company zbank --team checkout",
		Args:    cobra.NoArgs,
		Short:   "configure AWS credentials",
		Long:    "configure AWS credentials by modifying the .aws folder content",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			logger := log.NewLogger().WithContext(ctx).WithFields(logrus.Fields{"awsAuthMode": env.AWSAuthMode, "awsProfile": env.AWSProfile})

			logger.Infof("Checking does AWS config file %s contains the profile [%s]", env.AWSGeneratedProfile, env.AWSConfigFile)
			data, err := os.ReadFile(env.AWSConfigFile)
			if err != nil && !errors.Is(err, os.ErrNotExist) {
				return errors.Errorf("error reading aws config file at %s", env.AWSConfigFile)
			}
			profileHeader := fmt.Sprintf("[%s]", env.AWSGeneratedProfile)
			if strings.Contains(string(data), profileHeader) {
				logger.Infof("AWS config file already has the profile [%s] configured", env.AWSGeneratedProfile)
				return nil
			}
			logger.Infof("AWS config file %s does not have the profile [%s] entry", env.AWSConfigFile, env.AWSGeneratedProfile)

			auth, err := getAuth()
			if err != nil {
				return errors.Wrap(err, "error getting aws auth")
			}

			ssmClient, err := secrets.NewSSM(ctx, auth)
			if err != nil {
				return errors.Wrap(err, "error creating ssm client")
			}

			identifier := secret.Identifier{
				Company:     env.Company,
				Team:        env.Team,
				Environment: env.Environment,
			}

			credentials, err := aws.GetStateAWSCredentials(ssmClient, &identifier, logger)
			if err != nil {
				return errors.Wrap(err, "error getting state aws credentials")
			}

			entry := aws.GenerateAWSCredentialsEntry(env.AWSGeneratedProfile, credentials.AccessKeyID, credentials.SecretAccessKey)

			logger.Infof("Adding an entry in AWS config file %s for profile [%s]", env.AWSConfigFile, env.AWSGeneratedProfile)

			f, err := os.OpenFile(env.AWSConfigFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
			if err != nil {
				return errors.Wrap(err, "error opening aws config file")
			}
			if _, err := f.WriteString("\n" + entry + "\n"); err != nil {
				return errors.Wrap(err, "error writing entry to aws config file")
			}
			if err := f.Close(); err != nil {
				return errors.Wrap(err, "error closing aws config file")
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&env.AWSAuthMode, "auth-mode", "a", env.AWSAuthMode, "AWS auth mode")
	if err := cmd.MarkFlagRequired("auth-mode"); err != nil {
		fmt.Println(err) //nolint
		common.Failure(3101)
	}
	cmd.Flags().StringVarP(&env.AWSProfile, "profile", "p", env.AWSProfile, "AWS profile")
	cmd.Flags().StringVarP(&env.AWSGeneratedProfile, "generated-profile", "g", env.AWSGeneratedProfile, "Generated AWS profile")
	if err := cmd.MarkFlagRequired("generated-profile"); err != nil {
		fmt.Println(err) //nolint
		common.Failure(3102)
	}
	cmd.Flags().StringVarP(&env.Company, "company", "c", env.Company, "Company name")
	if err := cmd.MarkFlagRequired("company"); err != nil {
		fmt.Println(err) //nolint
		common.Failure(3103)
	}
	cmd.Flags().StringVarP(&env.Team, "team", "t", env.Team, "Team name")
	cmd.Flags().StringVarP(&env.Environment, "environment", "e", env.Environment, "Environment name")
	cmd.Flags().StringVarP(&env.AWSAccessKeyID, "access-key-id", "k", env.AWSAccessKeyID, "AWS Access Key ID")
	cmd.Flags().StringVarP(&env.AWSSecretAccessKey, "secret-access-key", "s", env.AWSSecretAccessKey, "AWS Secret Access Key")

	return cmd
}

func getAuth() (*aws.Auth, error) {
	auth := aws.Auth{Region: env.AWSRegion}
	switch env.AWSAuthMode {
	case aws.AuthModeProfile:
		if env.AWSProfile == "" {
			return nil, errors.New("profile not provided")
		}
		auth.Mode = aws.AuthModeProfile
		auth.Profile = env.AWSProfile
	case aws.AuthModeStatic:
		if env.AWSAccessKeyID == "" {
			return nil, errors.New("aws access key id not provided")
		}
		if env.AWSSecretAccessKey == "" {
			return nil, errors.New("aws secret access key not provided")
		}
		auth.Mode = aws.AuthModeStatic
		auth.AccessKeyID = env.AWSAccessKeyID
		auth.SecretAccessKey = env.AWSSecretAccessKey
	default:
		auth.Mode = aws.AuthModeDefault
	}
	return &auth, nil
}
