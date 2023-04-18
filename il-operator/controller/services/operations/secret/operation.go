package secret

import (
	"context"
	"strings"

	secret2 "github.com/compuzest/zlifecycle-il-operator/controller/common/secret"

	"github.com/compuzest/zlifecycle-il-operator/controller/util"

	"github.com/compuzest/zlifecycle-il-operator/controller/codegen/secret"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var (
	ErrTerraformStateConfigMissing = errors.New("terraform state config not configured in any scope")
	ErrAWSCredentialsMissing       = errors.New("aws credentials not configured in any scope")
)

func GetCustomerTerraformStateConfig(ctx context.Context, client secret2.API, meta *secret.Identifier, log *logrus.Entry) (*secret2.TerraformStateConfig, error) {
	secretsToFetch := []string{util.StateBucketSecret, util.StateLockTableSecret}
	log.Info("Checking for terraform state config in environment scope")
	scrts, err := client.GetSecrets(
		ctx,
		getEnvironmentScopeSecrets(
			meta.Company,
			meta.Team,
			meta.Environment,
			secretsToFetch...,
		)...,
	)
	if err != nil {
		return nil, errors.Wrap(err, "error getting terraform state config from team scope secrets")
	}
	if cfg, exist := checkForTerraformStateConfig(scrts); exist {
		log.Info("Terraform state config found in environment scope")
		return cfg, nil
	}

	log.Info("Checking for terraform state config in team scope")
	scrts, err = client.GetSecrets(
		ctx,
		getTeamScopeSecrets(
			meta.Company,
			meta.Team,
			secretsToFetch...,
		)...,
	)
	if err != nil {
		return nil, errors.Wrap(err, "error getting terraform state config from team scope secrets")
	}
	if cfg, exist := checkForTerraformStateConfig(scrts); exist {
		log.Info("Terraform state config found in team scope")
		return cfg, nil
	}

	log.Info("Checking for terraform state config in company scope")
	scrts, err = client.GetSecrets(ctx, getCompanyScopeSecrets(meta.Company, secretsToFetch...)...)
	if err != nil {
		return nil, errors.Wrap(err, "error getting AWS credentials from company scope secrets")
	}
	if cfg, exist := checkForTerraformStateConfig(scrts); exist {
		log.Info("Terraform state config found in company scope")
		return cfg, nil
	}

	return nil, errors.WithStack(ErrTerraformStateConfigMissing)
}

func GetAWSCredentials(ctx context.Context, client secret2.API, meta *secret.Identifier, log *logrus.Entry) (*secret2.AWSCredentials, error) {
	secretsToFetch := []string{util.AWSAccessKeyID, util.AWSSecretAccessKey, util.AWSSessionToken}
	log.Info("Checking for AWS credentials in environment scope")
	scrts, err := client.GetSecrets(
		ctx,
		getEnvironmentScopeSecrets(
			meta.Company,
			meta.Team,
			meta.Environment,
			secretsToFetch...,
		)...,
	)
	if err != nil {
		return nil, errors.Wrap(err, "error getting AWS credentials from team scope secrets")
	}
	if credentials, exist := checkForAWSCredentials(scrts); exist {
		log.Info("AWS credentials found in environment scope")
		return credentials, nil
	}

	log.Info("Checking for AWS credentials in team scope")
	scrts, err = client.GetSecrets(
		ctx,
		getTeamScopeSecrets(
			meta.Company,
			meta.Team,
			secretsToFetch...,
		)...,
	)
	if err != nil {
		return nil, errors.Wrap(err, "error getting AWS credentials from team scope secrets")
	}
	if credentials, exist := checkForAWSCredentials(scrts); exist {
		log.Info("AWS credentials found in team scope")
		return credentials, nil
	}

	log.Info("Checking for AWS credentials in company scope")
	scrts, err = client.GetSecrets(ctx, getCompanyScopeSecrets(meta.Company, secretsToFetch...)...)
	if err != nil {
		return nil, errors.Wrap(err, "error getting AWS credentials from company scope secrets")
	}
	if credentials, exist := checkForAWSCredentials(scrts); exist {
		log.Info("AWS credentials found in company scope")
		return credentials, nil
	}

	return nil, errors.WithStack(ErrAWSCredentialsMissing)
}

func getCompanyScopeSecrets(company string, keys ...string) []string {
	secretKeys := make([]string, 0, len(keys))
	for _, k := range keys {
		secretKeys = append(secretKeys, secret.GenerateOrgSecretKey(company, k))
	}
	return secretKeys
}

func getTeamScopeSecrets(company, team string, keys ...string) []string {
	secretKeys := make([]string, 0, len(keys))
	for _, k := range keys {
		secretKeys = append(secretKeys, secret.GenerateTeamSecretKey(company, team, k))
	}
	return secretKeys
}

func getEnvironmentScopeSecrets(company, team, environment string, keys ...string) []string {
	secretKeys := make([]string, 0, len(keys))
	for _, k := range keys {
		secretKeys = append(secretKeys, secret.GenerateEnvironmentSecretKey(company, team, environment, k))
	}
	return secretKeys
}

func checkForTerraformStateConfig(scrts []*secret2.Secret) (cfg *secret2.TerraformStateConfig, exist bool) {
	var bucket, lockTable string

	for _, scrt := range scrts {
		if strings.HasSuffix(scrt.Key, util.StateBucketSecret) && scrt.Exists {
			bucket = *scrt.Value
		}
		if strings.HasSuffix(scrt.Key, util.StateLockTableSecret) && scrt.Exists {
			lockTable = *scrt.Value
		}
	}

	if bucket != "" && lockTable != "" {
		return &secret2.TerraformStateConfig{
			Bucket:    bucket,
			LockTable: lockTable,
		}, true
	}
	return nil, false
}

func checkForAWSCredentials(scrts []*secret2.Secret) (credentials *secret2.AWSCredentials, exist bool) {
	var accessKeyID, secretAccessKey, sessionToken string

	for _, scrt := range scrts {
		if strings.HasSuffix(scrt.Key, util.AWSAccessKeyID) && scrt.Exists {
			accessKeyID = *scrt.Value
		}
		if strings.HasSuffix(scrt.Key, util.AWSSecretAccessKey) && scrt.Exists {
			secretAccessKey = *scrt.Value
		}
		if strings.HasSuffix(scrt.Key, util.AWSSessionToken) && scrt.Exists {
			sessionToken = *scrt.Value
		}
	}

	if accessKeyID != "" && secretAccessKey != "" {
		return &secret2.AWSCredentials{
			AccessKeyID:     accessKeyID,
			SecretAccessKey: secretAccessKey,
			SessionToken:    sessionToken,
		}, true
	}
	return nil, false
}
