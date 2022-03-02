package secrets

import (
	"fmt"

	"github.com/pkg/errors"
)

type Meta struct {
	Company              string
	Team                 string
	Environment          string
	EnvironmentComponent string
}

func CreateKey(name string, scope string, meta *Meta) (string, error) {
	if meta.Company == "" {
		return "", errors.New("missing company name in secret meta")
	}
	switch scope {
	case "org":
		return GenerateOrgSecretKey(meta.Company, name), nil
	case "team":
		if meta.Team == "" {
			return "", errors.New("missing team name in secret meta")
		}
		return GenerateTeamSecretKey(meta.Company, meta.Team, name), nil
	case "environment":
		if meta.Team == "" || meta.Environment == "" {
			return "", errors.New("missing team/environment name in secret meta")
		}
		return GenerateEnvironmentSecretKey(meta.Company, meta.Team, meta.Environment, name), nil
	case "component":
		if meta.Team == "" || meta.Environment == "" || meta.EnvironmentComponent == "" {
			return "", errors.New("missing team/environment/component name in secret meta")
		}
		return GenerateEnvironmentComponentSecretKey(meta.Company, meta.Team, meta.Environment, meta.EnvironmentComponent, name), nil
	default:
		return "", errors.Errorf("invalid scope: %s", scope)
	}
}

func GenerateOrgSecretKey(company, key string) string {
	return fmt.Sprintf("/%s/%s", company, key)
}

func GenerateTeamSecretKey(company, team, key string) string {
	return fmt.Sprintf("/%s/%s/%s", company, team, key)
}

func GenerateEnvironmentSecretKey(company, team, environment, key string) string {
	return fmt.Sprintf("/%s/%s/%s/%s", company, team, environment, key)
}

func GenerateEnvironmentComponentSecretKey(company, team, environment, component, key string) string {
	return fmt.Sprintf("/%s/%s/%s/%s/%s", company, team, environment, component, key)
}
