package argocd

import "errors"
import (
	env "github.com/compuzest/zlifecycle-il-operator/controllers/util/env"
)

func getArgocdCredentialsFromEnv() (*Credentials, error) {
	username := env.Config.ArgocdUsername
	password := env.Config.ArgocdPassword
	if username == "" || password == "" {
		return nil, errors.New("missing 'ARGOCD_USERNAME' or 'ARGOCD_PASSWORD' env variables")
	}

	creds := Credentials{Username: username, Password: password}

	return &creds, nil
}
