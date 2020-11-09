package controllers

import (
       	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	appv1 "github.com/argoproj/argo-cd/pkg/apis/application/v1alpha1"
)

func GenerateYaml() *appv1.Application {
    return &appv1.Application{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Application",
			APIVersion: "argoproj.io/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-pod",
		},
		Spec: appv1.ApplicationSpec {
                        Project: "default",
                        SyncPolicy: &appv1.SyncPolicy {
                            Automated: &appv1.SyncPolicyAutomated {
                                Prune: true,
                                SelfHeal: true,
                            },
                        },
                        Destination: appv1.ApplicationDestination {
                                Server: "",
                                Namespace: "",
                        },
                        Source: appv1.ApplicationSource {
                                RepoURL: "",
                                Path: "",
                                TargetRevision: "",
                                Helm: &appv1.ApplicationSourceHelm {
                                        Values: "",
                                },
                        },
		},
	}
}
