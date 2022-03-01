package controllers

import (
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/compuzest/zlifecycle-il-operator/controllers/env"
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

func checkReconcileMode(log *logrus.Entry) (end bool) {
	end = false
	if env.Config.ReconcileMode == "noop" {
		log.Info("Reconcile mode configured as noop, ending reconcile...")
		end = true
	}
	return
}
