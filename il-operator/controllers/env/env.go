package env

import (
	"fmt"
	"os"
)

type config struct {
	App  string
	Mode string

	Environment string

	CompanyName      string
	CompanyNamespace string

	TelemetryEnvironment string
	SlackWebhookURL      string
	EnableErrorNotifier  string

	ILZLifecycleRepositoryURL string
	ILTerraformRepositoryURL  string
	ILCompanyFolder           string
	ILTeamFolder              string
	ILConfigWatcherFolder     string

	SharedAWSCredsSecret string
	AWSRegion            string

	TerraformDefaultVersion          string
	TerraformDefaultAWSRegion        string
	TerraformDefaultSharedAWSRegion  string
	TerraformDefaultSharedAWSProfile string
	TerraformDefaultSharedAWSAlias   string

	// git
	GitHelmChartsRepository string
	GitCoreRepositoryOwner  string
	GitILRepositoryOwner    string
	GitSSHSecretName        string
	GitServiceAccountName   string
	GitServiceAccountEmail  string
	GitToken                string
	GitRepositoryBranch     string

	// github
	GitHubWebhookSecret              string
	GitHubCompanyOrganization        string
	GitHubAppIDCompany               string
	GitHubAppSecretNameCompany       string
	GitHubAppSecretNamespaceCompany  string
	GitHubCompanyAuthMethod          string
	GitHubAppIDInternal              string
	GitHubAppSecretNameInternal      string
	GitHubAppSecretNamespaceInternal string
	GitHubInternalAuthMethod         string

	// kubernetes
	KubernetesDisableWebhooks             string
	KubernetesCertDir                     string
	KubernetesServiceNamespace            string
	KubernetesDisableEnvironmentFinalizer string
	KubernetesEnvironmentFinalizerName    string
	KubernetesAPIURL                      string
	KubernetesOperatorWatchedNamespace    string
	KubernetesOperatorWatchedResources    string

	// new relic
	NewRelicAPIKey string
	EnableNewRelic string

	// argocd
	ArgocdServerURL     string
	ArgocdWebhookURL    string
	ArgocdUsername      string
	ArgocdPassword      string
	ArgocdRBACConfigMap string

	// argo workflows
	ArgoWorkflowsServerURL         string
	ArgoWorkflowsWorkflowNamespace string

	ZLifecycleStateManagerURL string
	ZLifecycleAPIURL          string

	// test
	ReconcileMode string
	SkipReconcile string
}

// Config exposes vars used throughout the operator.
var Config = config{
	App:                  "zlifecycle-il-operator",
	Mode:                 getOr("MODE", "cloud"),
	TelemetryEnvironment: getOr("TELEMETRY_ENVIRONMENT", "dev"),
	SlackWebhookURL:      os.Getenv("SLACK_WEBHOOK_URL"),
	EnableErrorNotifier:  getOr("ENABLE_ERROR_NOTIFIER", "false"),

	Environment: os.Getenv("ENVIRONMENT"),

	// company/customer config
	CompanyName:      os.Getenv("COMPANY_NAME"),
	CompanyNamespace: APINamespace(),

	// k8s
	KubernetesDisableWebhooks:             getOr("KUBERNETES_DISABLE_WEBHOOKS", "false"),
	KubernetesCertDir:                     os.Getenv("KUBERNETES_CERT_DIR"),
	KubernetesServiceNamespace:            getOr("KUBERNETES_SERVICE_NAMESPACE", "zlifecycle-system"),
	KubernetesDisableEnvironmentFinalizer: getOr("KUBERNETES_DISABLE_ENVIRONMENT_FINALIZER", "false"),
	KubernetesEnvironmentFinalizerName:    "zlifecycle.compuzest.com/github-finalizer",
	KubernetesAPIURL:                      "https://kubernetes.default.svc",
	KubernetesOperatorWatchedNamespace:    getOr("KUBERNETES_OPERATOR_WATCHED_NAMESPACE", "zlifecycle-system"),
	KubernetesOperatorWatchedResources:    getOr("KUBERNETES_OPERATOR_WATCHED_RESOURCES", "company,team,environment"),

	// il
	ILZLifecycleRepositoryURL: os.Getenv("IL_ZLIFECYCLE_REPOSITORY_URL"),
	ILTerraformRepositoryURL:  os.Getenv("IL_TERRAFORM_REPOSITORY_URL"),
	ILCompanyFolder:           getOr("IL_COMPANY_FOLDER", "company"),
	ILTeamFolder:              getOr("IL_TEAM_FOLDER", "team"),
	ILConfigWatcherFolder:     getOr("IL_CONFIG_WATCHER_FOLDER", "config-watcher"),

	// aws
	SharedAWSCredsSecret: getOr("AWS_SHARED_CREDS_SECRET", "shared-aws-creds"),
	AWSRegion:            getOr("AWS_REGION", "us-east-1"),

	// terraform config
	TerraformDefaultVersion:          getOr("TERRAFORM_DEFAULT_VERSION", "1.0.9"),
	TerraformDefaultAWSRegion:        getOr("TERRAFORM_DEFAULT_REGION", "us-east-1"),
	TerraformDefaultSharedAWSRegion:  getOr("TERRAFORM_DEFAULT_SHARED_REGION", "us-east-1"),
	TerraformDefaultSharedAWSProfile: getOr("TERRAFORM_DEFAULT_SHARED_PROFILE", "compuzest-shared"),
	TerraformDefaultSharedAWSAlias:   getOr("TERRAFORM_DEFAULT_SHARED_ALIAS", "shared"),

	// new relic
	EnableNewRelic: getOr("ENABLE_NEW_RELIC", "false"),
	NewRelicAPIKey: os.Getenv("NEW_RELIC_API_KEY"),

	// git
	GitHelmChartsRepository: os.Getenv("GIT_HELM_CHARTS_REPOSITORY"),
	GitCoreRepositoryOwner:  getOr("GIT_CORE_REPOSITORY_OWNER", "CompuZest"),
	GitILRepositoryOwner:    getOr("GIT_IL_REPOSITORY_OWNER", "zlifecycle-il"),
	GitSSHSecretName:        getOr("GIT_SSH_SECRET_NAME", "zlifecycle-operator-ssh"),
	GitServiceAccountName:   getOr("GIT_SERVICE_ACCOUNT_NAME", "zLifecycle"),
	GitServiceAccountEmail:  getOr("GIT_SERVICE_ACCOUNT_EMAIL", "zLifecycle@compuzest.com"),
	GitToken:                os.Getenv("GIT_TOKEN"),
	GitRepositoryBranch:     getOr("GIT_REPOSITORY_BRANCH", "main"),

	// github
	GitHubWebhookSecret:              os.Getenv("GITHUB_WEBHOOK_SECRET"),
	GitHubCompanyOrganization:        os.Getenv("GITHUB_COMPANY_ORGANIZATION"),
	GitHubAppIDCompany:               os.Getenv("GITHUB_APP_ID_COMPANY"),
	GitHubAppSecretNameCompany:       getOr("GITHUB_APP_SECRET_NAME_COMPANY", "public-github-app-ssh"),
	GitHubAppSecretNamespaceCompany:  SystemNamespace(),
	GitHubCompanyAuthMethod:          getOr("GITHUB_COMPANY_AUTH_METHOD", "ssh"),
	GitHubAppIDInternal:              os.Getenv("GITHUB_APP_ID_INTERNAL"),
	GitHubAppSecretNameInternal:      getOr("GITHUB_APP_SECRET_NAME_INTERNAL", "internal-github-app-ssh"),
	GitHubAppSecretNamespaceInternal: SystemNamespace(),
	GitHubInternalAuthMethod:         getOr("GITHUB_INTERNAL_AUTH_METHOD", "ssh"),

	// argocd
	ArgocdServerURL:     getOr("ARGOCD_SERVER_URL", fmt.Sprintf("http://argocd-%s-server.%s.svc.cluster.local", CompanyName(), ArgocdNamespace())),
	ArgocdWebhookURL:    os.Getenv("ARGOCD_WEBHOOK_URL"),
	ArgocdUsername:      os.Getenv("ARGOCD_USERNAME"),
	ArgocdPassword:      os.Getenv("ARGOCD_PASSWORD"),
	ArgocdRBACConfigMap: getOr("ARGOCD_RBAC_CONFIG_MAP", "argocd-rbac-cm"),

	// argo workflows
	ArgoWorkflowsServerURL: getOr("ARGOWORKFLOWS_URL", fmt.Sprintf(
		"http://argo-workflow-server.%s.svc.cluster.local:2746",
		ArgoWorkflowsNamespace(),
	)),
	ArgoWorkflowsWorkflowNamespace: getOr("ARGOWORKFLOWS_WORKFLOW_NAMESPACE", ExecutorNamespace()),

	// zlifecycle
	ZLifecycleStateManagerURL: getOr(
		"ZLIFECYCLE_STATE_MANAGER_URL",
		fmt.Sprintf("http://zlifecycle-state-manager.%s.svc.cluster.local:8080", StateManagerNamespace()),
	),
	ZLifecycleAPIURL: getOr("ZLIFECYCLE_API_URL", fmt.Sprintf(
		"http://zlifecycle-api.%s.svc.cluster.local",
		APINamespace(),
	)),

	// test
	ReconcileMode: getOr("RECONCILE_MODE", "normal"),
	SkipReconcile: getOr("SKIP_RECONCILE", ""),
}

func CompanyName() string {
	return os.Getenv("COMPANY_NAME")
}

func APINamespace() string {
	val, exists := os.LookupEnv("COMPANY_NAME")
	if exists {
		return fmt.Sprintf("%s-system", val)
	}
	return "zlifecycle-ui"
}

func StateManagerNamespace() string {
	val, exists := os.LookupEnv("COMPANY_NAME")
	if exists {
		return fmt.Sprintf("%s-system", val)
	}
	return "zlifecycle-il-operator-system"
}

func ArgocdNamespace() string {
	val, exists := os.LookupEnv("COMPANY_NAME")
	if exists {
		return fmt.Sprintf("%s-system", val)
	}
	return "argocd"
}

func ArgoWorkflowsNamespace() string {
	val, exists := os.LookupEnv("COMPANY_NAME")
	if exists {
		return fmt.Sprintf("%s-executor", val)
	}
	return "argocd"
}

func SystemNamespace() string {
	val, exists := os.LookupEnv("COMPANY_NAME")
	if exists {
		return fmt.Sprintf("%s-system", val)
	}
	return "zlifecycle-il-operator"
}

func ConfigNamespace() string {
	val, exists := os.LookupEnv("COMPANY_NAME")
	if exists {
		return fmt.Sprintf("%s-config", val)
	}
	return "zlifecycle"
}

func ExecutorNamespace() string {
	val, exists := os.LookupEnv("COMPANY_NAME")
	if exists {
		return fmt.Sprintf("%s-executor", val)
	}
	return "argocd"
}

func getOr(key string, defaultValue string) string {
	val, exists := os.LookupEnv(key)
	if exists && val != "" {
		return val
	}
	return defaultValue
}
