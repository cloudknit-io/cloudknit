package controllers

import (
	"strings"

	"github.com/compuzest/zlifecycle-il-operator/controllers/util/env"
)

var (
	helmChartsRepo         = env.Config.GitHelmChartsRepository
	zlILRepoURL            = env.Config.ILZLifecycleRepositoryURL
	ilRepoOwner            = env.Config.GitILRepositoryOwner
	githubSvcAccntName     = env.Config.GitServiceAccountName
	githubSvcAccntEmail    = env.Config.GitServiceAccountEmail
	gitHubWebhookSecret    = env.Config.GitHubWebhookSecret
	argocdHookURL          = env.Config.ArgocdWebhookURL
	argocdServerURL        = env.Config.ArgocdServerURL
	argoWorkflowsNamespace = env.Config.ArgoWorkflowsWorkflowNamespace
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
