package controllers

import (
	appv1 "github.com/argoproj/argo-cd/pkg/apis/application/v1alpha1"
	stablev1alpha1 "github.com/compuzest/environment-operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GenerateYaml(terraformConfig stablev1alpha1.TerraformConfig) *appv1.Application {
	return &appv1.Application{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Application",
			APIVersion: "argoproj.io/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      terraformConfig.Name,
			Namespace: "argo",
			Annotations: map[string]string{
				"argocd.argoproj.io/sync-wave": "2",
			},
			Finalizers: []string{
				"resources-finalizer.argocd.argoproj.io",
			},
		},
		Spec: appv1.ApplicationSpec{
			Project: "default",
			SyncPolicy: &appv1.SyncPolicy{
				Automated: &appv1.SyncPolicyAutomated{
					Prune:    true,
					SelfHeal: true,
				},
			},
			Destination: appv1.ApplicationDestination{
				Server:    "https://192.168.1.155:60792",
				Namespace: "default",
			},
			Source: appv1.ApplicationSource{
				RepoURL:        terraformConfig.Module.Source,
				Path:           terraformConfig.Module.Path,
				TargetRevision: "HEAD",
				Helm: &appv1.ApplicationSourceHelm{
					Values: "",
				},
			},
		},
	}
}
