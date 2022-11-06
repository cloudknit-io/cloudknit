package aws

import (
	"fmt"
	"strings"

	"github.com/compuzest/zlifecycle-internal-cli/app/lib/secret"
	"github.com/compuzest/zlifecycle-internal-cli/app/util"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

//nolint
const credentialsFormat = `[%s]
aws_access_key_id = %s
aws_secret_access_key = %s
region = %s`

func GenerateAWSCredentialsEntry(profile, accessKeyID, secretAccessKey string, region string) string {
	return fmt.Sprintf(credentialsFormat, profile, accessKeyID, secretAccessKey, region)
}

func GetStateAWSCredentials(client secret.API, meta *secret.Identifier, log *logrus.Entry) (*Credentials, error) {
	secretsToFetch := []string{util.StateAWSAccessKeyID, util.StateAWSSecretAccessKey, util.StateAWSRegion}
	log.Info("Checking for AWS creds in environment scope")
	scrts, err := client.GetSecrets(
		getEnvironmentScopeSecrets(
			meta.Company,
			meta.Team,
			meta.Environment,
			secretsToFetch...,
		)...,
	)
	if err != nil {
		return nil, errors.Wrap(err, "error getting AWS creds from team scope secrets")
	}
	if creds, exist := checkForAWSCreds(scrts); exist {
		log.Info("AWS creds found in environment scope")
		return creds, nil
	}

	log.Info("Checking for AWS creds in team scope")
	scrts, err = client.GetSecrets(
		getTeamScopeSecrets(
			meta.Company,
			meta.Team,
			secretsToFetch...,
		)...,
	)
	if err != nil {
		return nil, errors.Wrap(err, "error getting AWS creds from team scope secrets")
	}
	if creds, exist := checkForAWSCreds(scrts); exist {
		log.Info("AWS creds found in team scope")
		return creds, nil
	}

	log.Info("Checking for AWS creds in company scope")
	scrts, err = client.GetSecrets(getCompanyScopeSecrets(meta.Company, secretsToFetch...)...)
	if err != nil {
		return nil, errors.Wrap(err, "error getting AWS creds from company scope secrets")
	}
	if creds, exist := checkForAWSCreds(scrts); exist {
		log.Info("AWS creds found in company scope")
		return creds, nil
	}

	return nil, errors.New("AWS creds not configured in any scope")
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

func checkForAWSCreds(scrts []*secret.Secret) (creds *Credentials, exist bool) {
	var accessKeyID, secretAccessKey, region string

	// Setting default region to us-east-1
	region = "us-east-1"

	for _, scrt := range scrts {
		if strings.HasSuffix(scrt.Key, util.StateAWSAccessKeyID) && scrt.Exists {
			accessKeyID = *scrt.Value
		}
		if strings.HasSuffix(scrt.Key, util.StateAWSSecretAccessKey) && scrt.Exists {
			secretAccessKey = *scrt.Value
		}
		if strings.HasSuffix(scrt.Key, util.StateAWSRegion) && scrt.Exists {
			region = *scrt.Value
		}
	}

	if accessKeyID != "" && secretAccessKey != "" {
		return &Credentials{
			AccessKeyID:     accessKeyID,
			SecretAccessKey: secretAccessKey,
			Region: region,
		}, true
	}
	return nil, false
}
