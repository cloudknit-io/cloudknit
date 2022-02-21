package secrets

import (
	"fmt"

	"github.com/pkg/errors"
)

type SecretMeta struct {
	Company              string
	Team                 string
	Environment          string
	EnvironmentComponent string
}

func CreateKey(name string, scope string, meta SecretMeta) (string, error) {
	if meta.Company == "" {
		return "", errors.New("missing company name in secret meta")
	}
	switch scope {
	case "org":
		return fmt.Sprintf("/%s/%s", meta.Company, name), nil
	case "team":
		if meta.Team == "" {
			return "", errors.New("missing team name in secret meta")
		}
		return fmt.Sprintf("/%s/%s/%s", meta.Company, meta.Team, name), nil
	case "environment":
		if meta.Team == "" || meta.Environment == "" {
			return "", errors.New("missing team/environment name in secret meta")
		}
		return fmt.Sprintf("/%s/%s/%s/%s", meta.Company, meta.Team, meta.Environment, name), nil
	case "component":
		if meta.Team == "" || meta.Environment == "" || meta.EnvironmentComponent == "" {
			return "", errors.New("missing team/environment/component name in secret meta")
		}
		return fmt.Sprintf("/%s/%s/%s/%s/%s", meta.Company, meta.Team, meta.Environment, meta.EnvironmentComponent, name), nil
	default:
		return "", errors.Errorf("invalid scope: %s", scope)
	}
}
