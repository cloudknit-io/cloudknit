package argocd

import (
	"errors"
	"os"
)

func GetArgocdServerAddr() string {
	addr, exists := os.LookupEnv("ARGOCD_API_URL")
	if exists {
		return addr
	} else {
		return "http://argocd-server.argocd.svc.cluster.local"
	}
}

func getArgocdCredentialsFromEnv() (*ArgocdCredentials, error) {
	username := os.Getenv("ARGOCD_USERNAME")
	password := os.Getenv("ARGOCD_PASSWORD")
	if username == "" || password == "" {
		return nil, errors.New("missing 'ARGOCD_USERNAME' or 'ARGOCD_PASSWORD' env variables")
	}

	creds := ArgocdCredentials{Username: username, Password: password}

	return &creds, nil
}
