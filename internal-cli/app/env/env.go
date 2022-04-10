package env

import "os"

var Version = "0.0.11" //nolint

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
