package env

import (
	"os"
)

type config struct {
	ZlifecycleOwner     string
	ZlifecycleMasterRepoSshSecret string
	ZlifecycleOperatorNamespace   string
	CompanyName         string
	ILRepoName          string
	ILRepoURL           string
	ILRepoSourceOwner   string

	GithubSvcAccntName  string
	GithubSvcAccntEmail string
	GitHubAuthToken     string
	GitHubWebhookSecret string
	RepoBranch          string

	HelmChartsRepo string
	K8sAPIURL      string

	ArgocdServerUrl string
	ArgocdHookUrl   string
	ArgocdUsername  string
	ArgocdPassword  string
}

// Various config vars used throughout the operator
var Config = config{
	ZlifecycleOwner:               getZlifecycleOwner(),
	ZlifecycleMasterRepoSshSecret: getZlifecyleMasterRepoSshSecret(),
	ZlifecycleOperatorNamespace:   os.Getenv("ZLIFECYCLE_OPERATOR_NAMESPACE"),

	CompanyName:         os.Getenv("companyName"),
	ILRepoName:          os.Getenv("ilRepoName"),
	ILRepoURL:           os.Getenv("ilRepo"),
	ILRepoSourceOwner:   os.Getenv("ilRepoSourceOwner"),

	GithubSvcAccntName:  "zLifecycle",
	GithubSvcAccntEmail: "zLifecycle@compuzest.com",
	GitHubAuthToken:     os.Getenv("GITHUB_AUTH_TOKEN"),
	GitHubWebhookSecret: os.Getenv("GITHUB_WEBHOOK_SECRET"),
	RepoBranch:          "main",

	HelmChartsRepo:  os.Getenv("helmChartsRepo"),
	K8sAPIURL:       os.Getenv("K8s_API_URL"),

	ArgocdServerUrl:  getArgocdServerAddr(),
	ArgocdHookUrl:    os.Getenv("ARGOCD_WEBHOOK_URL"),
	ArgocdUsername:   os.Getenv("ARGOCD_USERNAME"),
	ArgocdPassword:   os.Getenv("ARGOCD_PASSWORD"),
}

func getZlifecyleMasterRepoSshSecret() string {
	val, exists := os.LookupEnv("ZLIFECYCLE_MASTER_SSH")
	if exists {
		return val
	} else {
		return "zlifecycle-master-ssh"
	}
}

func getZlifecycleOwner() string {
	val, exists := os.LookupEnv("GITHUB_ZLIFECYCLE_OWNER")
	if exists {
		return val
	} else {
		return "CompuZest"
	}
}

func getArgocdServerAddr() string {
	val, exists := os.LookupEnv("ARGOCD_URL")
	if exists {
		return val
	} else {
		return "http://argocd-server.argocd.svc.cluster.local"
	}
}

