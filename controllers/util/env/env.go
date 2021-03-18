package env

import (
	"os"
)

type config struct {
	ZlifecycleOwner   string
	WebhookSecret     string
	CompanyName       string
	ILRepoName        string
	ILRepoURL         string
	ILRepoSourceOwner string

	GithubSvcAccntName  string
	GithubSvcAccntEmail string
	GitHubAuthToken     string
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
	ZlifecycleOwner:   getZlifecycleOwner(),
	WebhookSecret:     getWebhookSecret(),
	CompanyName:       os.Getenv("companyName"),
	ILRepoName:        os.Getenv("ilRepoName"),
	ILRepoURL:         os.Getenv("ilRepo"),
	ILRepoSourceOwner: os.Getenv("ilRepoSourceOwner"),

	GithubSvcAccntName:  "zLifecycle",
	GithubSvcAccntEmail: "zLifecycle@compuzest.com",
	GitHubAuthToken:     os.Getenv("GITHUB_AUTH_TOKEN"),
	RepoBranch:          "main",

	HelmChartsRepo:  os.Getenv("helmChartsRepo"),
	K8sAPIURL:       os.Getenv("K8s_API_URL"),

	ArgocdServerUrl:  getArgocdServerAddr(),
	ArgocdHookUrl:    os.Getenv("ARGOCD_WEBHOOK_URL"),
	ArgocdUsername:   os.Getenv("ARGOCD_USERNAME"),
	ArgocdPassword:   os.Getenv("ARGOCD_PASSWORD"),
}

func getWebhookSecret() string {
	secret, exists := os.LookupEnv("GITHUB_WEBHOOK_SECRET")
	if exists {
		return secret
	} else {
		return ""
	}
}

func getZlifecycleOwner() string {
	secret, exists := os.LookupEnv("GITHUB_ZLIFECYCLE_OWNER")
	if exists {
		return secret
	} else {
		return "CompuZest"
	}
}

func getArgocdServerAddr() string {
	addr, exists := os.LookupEnv("ARGOCD_URL")
	if exists {
		return addr
	} else {
		return "http://argocd-server.argocd.svc.cluster.local"
	}
}

