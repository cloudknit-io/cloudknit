package env

import (
	"os"
	"path/filepath"
)

var Version = "0.0.23" //nolint

const (
	TestModeIntegration = "integration"
	TestModeUnit        = "unit"
	TestModeAll         = "all"
)

var (
	TestMode = getOr("ZLI_TEST_MODE", "unit")

	TestDir             = "/tmp/zli_test"
	Company             = os.Getenv("ZLI_COMPANY")
	Team                = os.Getenv("ZLI_TEAM")
	Environment         = os.Getenv("ZLI_ENVIRONMENT")
	Component           = os.Getenv("ZLI_TEAM")
	Status              string
	Verbose             bool
	GitHubAppID         string
	GitHubAppIDInternal = getOr("GITHUB_APP_ID_INTERNAL", "172698")
	GitHubAppIDPublic   = getOr("GITHUB_APP_ID_PUBLIC", "172696")
	GitAuth             string
	GitToken            string
	GitCloneDir         = getOr("GIT_CLONE_DIR", ".")
	GitBaseURL          = getOr("GIT_BASE_URL", "github.com")
	GitConfigDir        = getOr("GIT_CONFIG_DIR", ".")
	GitSSHPath          = getOr("GIT_SSH_PATH", filepath.Join(os.Getenv("HOME"), ".ssh", "id_rsa"))
	AWSAuthMode         string
	AWSProfile          string
	AWSGeneratedProfile string
	AWSRegion           = getOr("AWS_REGION", "us-east-1")
	AWSConfigFile       = getOr("AWS_CONFIG_FILE", filepath.Join(os.Getenv("HOME"), ".aws", "credentials"))
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
