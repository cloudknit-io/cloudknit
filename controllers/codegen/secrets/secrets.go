package secrets

import (
	"fmt"
	v1 "github.com/compuzest/zlifecycle-il-operator/api/v1"
	"github.com/compuzest/zlifecycle-il-operator/controllers/env"

	"github.com/pkg/errors"
)

type Identifier struct {
	Company              string
	Team                 string
	Environment          string
	EnvironmentComponent string
}

func NewIdentifierFromEnvironment(e *v1.Environment) *Identifier {
	return &Identifier{
		Company:     env.Config.CompanyName,
		Team:        e.Spec.TeamName,
		Environment: e.Spec.EnvName,
	}
}

func (i *Identifier) GenerateKey(name, scope string) (string, error) {
	if i.Company == "" {
		return "", errors.New("missing company name in secret meta")
	}
	switch scope {
	case "org":
		return GenerateOrgSecretKey(i.Company, name), nil
	case "team":
		if i.Team == "" {
			return "", errors.New("missing team name in secret meta")
		}
		return GenerateTeamSecretKey(i.Company, i.Team, name), nil
	case "environment":
		if i.Team == "" || i.Environment == "" {
			return "", errors.New("missing team/environment name in secret meta")
		}
		return GenerateEnvironmentSecretKey(i.Company, i.Team, i.Environment, name), nil
	case "component":
		if i.Team == "" || i.Environment == "" || i.EnvironmentComponent == "" {
			return "", errors.New("missing team/environment/component name in secret meta")
		}
		return GenerateEnvironmentComponentSecretKey(i.Company, i.Team, i.Environment, i.EnvironmentComponent, name), nil
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
