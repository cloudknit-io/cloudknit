package env

import (
	"os"
	"path/filepath"
)

var Version = "0.0.14" //nolint

var (
	Company             string
	Team                string
	Environment         string
	Component           string
	Status              string
	Verbose             bool
	GitHubAppID         string
	GitHubAppIDInternal = "172698"
	GitHubAppIDPublic   = "172696"
	GitHubAppSSHPath    string
	GitAuth             string
	GitToken            string
	GitCloneDir         = "."
	GitBaseURL          = "github.com"
	GitConfigDir        = "."
	GitSSHPath          = filepath.Join(os.Getenv("HOME"), ".ssh", "id_rsa")
	StateManagerURL     = getOr(
		"STATE_MANAGER_URL",
		"http://zlifecycle-state-manager.zlifecycle-il-operator-system.svc.cluster.local:8080",
	)
	ArgoCDURL = getOr(
		"ARGOCD_URL",
		"http://argocd-server.argocd.svc.cluster.local:80",
	)
)

func getOr(key string, defaultValue string) string {
	val, exists := os.LookupEnv(key)
	if exists && val != "" {
		return val
	}
	return defaultValue
}
