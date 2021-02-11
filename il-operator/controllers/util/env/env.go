package env

import "os"

type config struct {
	CompanyName string
	ILRepoName  string
	ILRepoURL   string
	SourceOwner string

	GithubSvcAccntName  string
	GithubSvcAccntEmail string
	GitHubAuthToken     string
	RepoBranch          string

	HelmChartsRepo string
	K8sAPIURL      string
}

// Various config vars used throughout the operator
var Config = config{
	CompanyName: os.Getenv("companyName"),
	ILRepoName:  os.Getenv("ilRepoName"),
	ILRepoURL:   os.Getenv("ilRepo"),
	SourceOwner: "CompuZest",

	GithubSvcAccntName:  "zLifecycle",
	GithubSvcAccntEmail: "zLifecycle@compuzest.com",
	GitHubAuthToken:     os.Getenv("GITHUB_AUTH_TOKEN"),
	RepoBranch:          "main",

	HelmChartsRepo: os.Getenv("helmChartsRepo"),
	K8sAPIURL:      os.Getenv("K8s_API_URL"),
}
