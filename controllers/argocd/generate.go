package controllers

import (
	"fmt"
	appv1 "github.com/argoproj/argo-cd/pkg/apis/application/v1alpha1"
	stablev1alpha1 "github.com/compuzest/environment-operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GenerateEnvironmentApp(environment stablev1alpha1.Environment) *appv1.Application {
	return &appv1.Application{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "argoproj.io/v1alpha1",
			Kind:       "Application",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      environment.Spec.CustomerId + "-" + environment.Spec.Name,
			Namespace: "argo",
		},
		Spec: appv1.ApplicationSpec{
			Project: "default",
			SyncPolicy: &appv1.SyncPolicy{
				Automated: &appv1.SyncPolicyAutomated{
					Prune: true,
				},
			},
			Destination: appv1.ApplicationDestination{
				Server:    "https://kubernetes.default.svc",
				Namespace: "default",
			},
			Source: appv1.ApplicationSource{
				RepoURL:        "git@github.com:CompuZest/terraform-environment.git",
				Path:           environment.Spec.CustomerId + "/" + environment.Spec.Name,
				TargetRevision: "HEAD",
				Directory: &appv1.ApplicationSourceDirectory{
					Recurse: true,
				},
			},
		},
	}
}

func GenerateTerraformConfigApps(environment stablev1alpha1.Environment, terraformConfig stablev1alpha1.TerraformConfig) *appv1.Application {
	variables := ""

	for _, variable := range terraformConfig.Variables {
		variables += "\n- name:" + variable.Name + "\n  value:" + variable.Value
	}

	helmValues := fmt.Sprintf(`
        customer_id: "%s"
        env_name: %s
        name: %s
        module:
            source: %s
            path: %s
        variables: 
        %s
        `, environment.Spec.CustomerId,
		environment.Name,
		terraformConfig.Name,
		terraformConfig.Module.Source,
		terraformConfig.Module.Path,
		variables)
	return &appv1.Application{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Application",
			APIVersion: "argoproj.io/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      environment.Spec.CustomerId + "-" + environment.Spec.Name + "-" + terraformConfig.Name,
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
				Server:    "https://192.168.1.155:51231",
				Namespace: "default",
			},
			Source: appv1.ApplicationSource{
				RepoURL:        "git@github.com:CompuZest/helm-charts.git",
				Path:           "charts/terraform-config",
				TargetRevision: "HEAD",
				Helm: &appv1.ApplicationSourceHelm{
					Values: helmValues,
				},
			},
		},
	}
}
