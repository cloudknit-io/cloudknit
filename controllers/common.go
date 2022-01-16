package controllers

import (
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/env"
	"strings"
)

var (
	helmChartsRepo         = env.Config.HelmChartsRepo
	operatorSSHSecret      = env.Config.ZlifecycleMasterRepoSSHSecret
	operatorNamespace      = env.Config.ZlifecycleOperatorNamespace
	zlILRepoURL            = env.Config.ZLILRepoURL
	zlILRepoName           = env.Config.ZLILRepoName
	ilRepoOwner            = env.Config.ZlifecycleILRepoOwner
	githubSvcAccntName     = env.Config.GithubSvcAccntName
	githubSvcAccntEmail    = env.Config.GithubSvcAccntEmail
	gitHubWebhookSecret    = env.Config.GitHubWebhookSecret
	argocdHookURL          = env.Config.ArgocdHookURL
	argocdServerURL        = env.Config.ArgocdServerURL
	argoWorkflowsNamespace = env.Config.ArgoWorkflowsNamespace
)

func checkIsNamespaceWatched(namespace string) bool {
	watchedNamespace := env.Config.KubernetesOperatorWatchedNamespace
	return namespace == watchedNamespace
}

func checkIsResourceWatched(resource string) bool {
	watchedResources := strings.Split(env.Config.KubernetesOperatorWatchedResources, ",")

	for _, r := range watchedResources {
		if strings.EqualFold(strings.TrimSpace(r), resource) {
			return true
		}
	}

	return false
}
