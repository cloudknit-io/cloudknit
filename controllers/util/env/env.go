package env

import (
	"os"
)

type config struct {
	App string

	ZlifecycleILRepoOwner         string
	ZlifecycleMasterRepoSSHSecret string
	ZlifecycleOperatorNamespace   string
	ZlifecycleOperatorRepo        string
	CompanyName                   string
	ZLILRepoName                  string
	ZLILRepoURL                   string
	TFILRepoName                  string
	TFILRepoURL                   string
	ILRepoSourceOwner             string

	DefaultTerraformVersion string
	DefaultRegion           string
	DefaultSharedRegion     string
	DefaultSharedProfile    string
	DefaultSharedAlias      string

	GithubSvcAccntName  string
	GithubSvcAccntEmail string
	GitHubAuthToken     string
	GitHubWebhookSecret string
	GitHubOrg           string
	RepoBranch          string

	DisableWebhooks   string
	KubernetesCertDir string

	Mode string

	NewRelicAPIKey string
	EnableNewRelic string

	DisableEnvironmentFinalizer string
	EnvironmentFinalizer        string

	HelmChartsRepo string
	K8sAPIURL      string

	ArgocdServerURL string
	ArgocdHookURL   string
	ArgocdUsername  string
	ArgocdPassword  string

	ArgoWorkflowsServerURL string
	ArgoWorkflowsNamespace string

	APIURL string

	TelemetryEnvironment string

	SlackWebhookURL string
}

// Config exposes vars used throughout the operator.
var Config = config{
	App: "zlifecycle-il-operator",

	ZlifecycleILRepoOwner:         getOr("GITHUB_ZLIFECYCLE_OWNER", "zlifecycle-il"),
	ZlifecycleMasterRepoSSHSecret: getOr("ZLIFECYCLE_MASTER_SSH", "zlifecycle-operator-ssh"),
	ZlifecycleOperatorNamespace:   getOr("ZLIFECYCLE_OPERATOR_NAMESPACE", "zlifecycle-il-operator-system"),
	ZlifecycleOperatorRepo:        "zlifecycle-il-operator",

	CompanyName:       os.Getenv("companyName"),
	ZLILRepoName:      os.Getenv("ilRepoName"),
	ZLILRepoURL:       os.Getenv("ilRepo"),
	TFILRepoName:      os.Getenv("TF_IL_REPO_NAME"),
	TFILRepoURL:       os.Getenv("TF_IL_REPO_URL"),
	ILRepoSourceOwner: os.Getenv("ilRepoSourceOwner"),

	DefaultTerraformVersion: getOr("TERRAFORM_DEFAULT_VERSION", "1.0.9"),
	DefaultRegion:           getOr("TERRAFORM_DEFAULT_REGION", "us-east-1"),
	DefaultSharedRegion:     getOr("TERRAFORM_DEFAULT_SHARED_REGION", "us-east-1"),
	DefaultSharedProfile:    getOr("TERRAFORM_DEFAULT_SHARED_PROFILE", "compuzest-shared"),
	DefaultSharedAlias:      getOr("TERRAFORM_DEFAULT_SHARED_PROFILE", "shared"),

	DisableWebhooks:   getOr("DISABLE_WEBHOOKS", "false"),
	KubernetesCertDir: os.Getenv("KUBERNETES_CERT_DIR"),

	Mode: getOr("MODE", "cloud"),

	NewRelicAPIKey: os.Getenv("NEW_RELIC_API_KEY"),
	EnableNewRelic: getOr("ENABLE_NEW_RELIC", "false"),

	GithubSvcAccntName:  "zLifecycle",
	GithubSvcAccntEmail: "zLifecycle@compuzest.com",
	GitHubAuthToken:     os.Getenv("GITHUB_AUTH_TOKEN"),
	GitHubWebhookSecret: os.Getenv("GITHUB_WEBHOOK_SECRET"),
	GitHubOrg:           os.Getenv("GITHUB_ORG"),
	RepoBranch:          "main",

	DisableEnvironmentFinalizer: getOr("DISABLE_ENVIRONMENT_FINALIZER", "false"),
	EnvironmentFinalizer:        "zlifecycle.compuzest.com/github-finalizer",

	HelmChartsRepo: os.Getenv("helmChartsRepo"),
	K8sAPIURL:      "https://kubernetes.default.svc",

	ArgocdServerURL: getOr("ARGOCD_URL", "http://argocd-server.argocd.svc.cluster.local"),
	ArgocdHookURL:   os.Getenv("ARGOCD_WEBHOOK_URL"),
	ArgocdUsername:  os.Getenv("ARGOCD_USERNAME"),
	ArgocdPassword:  os.Getenv("ARGOCD_PASSWORD"),

	ArgoWorkflowsServerURL: getOr("ARGOWORKFLOWS_URL", "http://argo-workflow-server.argocd.svc.cluster.local:2746"),
	ArgoWorkflowsNamespace: "argocd",

	APIURL: getOr("API_URL", "http://zlifecycle-api.zlifecycle-ui.svc.cluster.local"),

	TelemetryEnvironment: getOr("TELEMETRY_ENVIRONMENT", "dev"),

	SlackWebhookURL: getOr("SLACK_WEBHOOK_URL", "https://hooks.slack.com/services/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"),
}

func getOr(key string, defaultValue string) string {
	val, exists := os.LookupEnv(key)
	if exists && val != "" {
		return val
	}
	return defaultValue
}
