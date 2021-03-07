package argocd

import (
	"errors"
	"os"
)

func GetArgocdServerAddr() string {
	addr, exists := os.LookupEnv("argocd_url")
	if exists {
		return addr
	} else {
		return "http://argocd-server.argocd.svc.cluster.local"
	}
}

func getArgocdCredentialsFromEnv() (*ArgocdCredentials, error) {
	username := os.Getenv("argocd_username")
	password := os.Getenv("argocd_password")
	if username == "" || password == "" {
		return nil, errors.New("missing 'argocd_username' or 'argocd_password' env variables")
	}

	creds := ArgocdCredentials{Username: username, Password: password}

	return &creds, nil
}
