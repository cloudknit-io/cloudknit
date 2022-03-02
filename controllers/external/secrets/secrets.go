package secrets

import (
	"strings"

	"github.com/compuzest/zlifecycle-il-operator/controllers/codegen/secrets"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	secretNameAccessKeyID     = "aws_access_key_id"
	secretNameSecretAccessKey = "aws_secret_access_key"
)

func GetAWSCreds(client API, meta *secrets.Meta, log *logrus.Entry) (*AWSCreds, error) {
	log.Info("Checking for AWS creds in environment scope")
	scrts, err := client.GetSecrets(
		getEnvironmentScopeSecrets(
			meta.Company,
			meta.Team,
			meta.Environment,
			secretNameAccessKeyID,
			secretNameSecretAccessKey,
		)...,
	)
	if err != nil {
		return nil, errors.Wrap(err, "error getting AWS creds from team scope secrets")
	}
	if creds, exist := checkForCreds(scrts); exist {
		log.Info("AWS creds found in environment scope")
		return creds, nil
	}

	log.Info("Checking for AWS creds in team scope")
	scrts, err = client.GetSecrets(
		getTeamScopeSecrets(
			meta.Company,
			meta.Team,
			secretNameAccessKeyID,
			secretNameSecretAccessKey,
		)...,
	)
	if err != nil {
		return nil, errors.Wrap(err, "error getting AWS creds from team scope secrets")
	}
	if creds, exist := checkForCreds(scrts); exist {
		log.Info("AWS creds found in team scope")
		return creds, nil
	}

	log.Info("Checking for AWS creds in company scope")
	scrts, err = client.GetSecrets(getCompanyScopeSecrets(meta.Company, secretNameAccessKeyID, secretNameSecretAccessKey)...)
	if err != nil {
		return nil, errors.Wrap(err, "error getting AWS creds from company scope secrets")
	}
	if creds, exist := checkForCreds(scrts); exist {
		log.Info("AWS creds found in company scope")
		return creds, nil
	}

	return nil, errors.New("AWS creds not configured in any scope")
}

func getCompanyScopeSecrets(company string, keys ...string) []string {
	secretKeys := make([]string, 0, len(keys))
	for _, k := range keys {
		secretKeys = append(secretKeys, secrets.GenerateOrgSecretKey(company, k))
	}
	return secretKeys
}

func getTeamScopeSecrets(company, team string, keys ...string) []string {
	secretKeys := make([]string, 0, len(keys))
	for _, k := range keys {
		secretKeys = append(secretKeys, secrets.GenerateTeamSecretKey(company, team, k))
	}
	return secretKeys
}

func getEnvironmentScopeSecrets(company, team, environment string, keys ...string) []string {
	secretKeys := make([]string, 0, len(keys))
	for _, k := range keys {
		secretKeys = append(secretKeys, secrets.GenerateEnvironmentSecretKey(company, team, environment, k))
	}
	return secretKeys
}

func checkForCreds(scrts []*Secret) (creds *AWSCreds, exist bool) {
	var accessKeyID, secretAccessKey string

	for _, scrt := range scrts {
		if strings.HasSuffix(scrt.Key, secretNameAccessKeyID) && scrt.Exists {
			accessKeyID = *scrt.Value
		}
		if strings.HasSuffix(scrt.Key, secretNameSecretAccessKey) && scrt.Exists {
			secretAccessKey = *scrt.Value
		}
	}

	if accessKeyID != "" && secretAccessKey != "" {
		return &AWSCreds{
			AccessKeyID:     accessKeyID,
			SecretAccessKey: secretAccessKey,
		}, true
	}
	return nil, false
}
