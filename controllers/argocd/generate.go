/* Copyright (C) 2020 CompuZest, Inc. - All Rights Reserved
 *
 * Unauthorized copying of this file, via any medium, is strictly prohibited
 * Proprietary and confidential
 *
 * NOTICE: All information contained herein is, and remains the property of
 * CompuZest, Inc. The intellectual and technical concepts contained herein are
 * proprietary to CompuZest, Inc. and are protected by trade secret or copyright
 * law. Dissemination of this information or reproduction of this material is
 * strictly forbidden unless prior written permission is obtained from CompuZest, Inc.
 */

package controllers

import (
	"fmt"
	appv1 "github.com/argoproj/argo-cd/pkg/apis/application/v1alpha1"
	stablev1alpha1 "github.com/compuzest/zlifecycle-il-operator/api/v1alpha1"
	config "github.com/compuzest/zlifecycle-il-operator/controllers/util/config"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GenerateEnvironmentApp(environment stablev1alpha1.Environment) *appv1.Application {

	return &appv1.Application{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "argoproj.io/v1alpha1",
			Kind:       "Application",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      environment.Spec.TeamName + "-" + environment.Spec.EnvName,
			Namespace: "argocd",
			Labels: map[string]string{
				"zlifecycle.com/model": "environment",
			},
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
				RepoURL:        config.ILRepoURL,
				Path:           environment.Spec.TeamName + "/" + environment.Spec.EnvName,
				TargetRevision: "HEAD",
				Directory: &appv1.ApplicationSourceDirectory{
					Recurse: true,
				},
			},
		},
		Status: appv1.ApplicationStatus{
			Sync: appv1.SyncStatus{
				ComparedTo: appv1.ComparedTo{
					Source: appv1.ApplicationSource{
						RepoURL: config.ILRepoURL,
					},
				},
				Status: "Synced",
			},
		},
	}
}

func GenerateTerraformConfigApps(environment stablev1alpha1.Environment, terraformConfig stablev1alpha1.TerraformConfig) *appv1.Application {

	helmValues := getHelmValues(environment, terraformConfig)

	return &appv1.Application{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Application",
			APIVersion: "argoproj.io/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      environment.Spec.TeamName + "-" + environment.Spec.EnvName + "-" + terraformConfig.ConfigName,
			Namespace: "argocd",
			Labels: map[string]string{
				"zlifecycle.com/model": "environment-component",
			},
			Finalizers: []string{
				"resources-finalizer.argocd.argoproj.io",
			},
		},
		Spec: appv1.ApplicationSpec{
			Project: "default",
			SyncPolicy: &appv1.SyncPolicy{
				Automated: &appv1.SyncPolicyAutomated{
					Prune: true,
				},
			},
			Destination: appv1.ApplicationDestination{
				Server:    config.K8sAPIURL,
				Namespace: "default",
			},
			Source: appv1.ApplicationSource{
				RepoURL:        config.HelmChartsRepo,
				Path:           "charts/terraform-config",
				TargetRevision: "HEAD",
				Helm: &appv1.ApplicationSourceHelm{
					Values: helmValues,
				},
			},
		},
		Status: appv1.ApplicationStatus{
			Sync: appv1.SyncStatus{
				ComparedTo: appv1.ComparedTo{
					Source: appv1.ApplicationSource{
						RepoURL: config.HelmChartsRepo,
					},
				},
				Status: "Synced",
			},
		},
	}
}

func getHelmValues(environment stablev1alpha1.Environment, terraformConfig stablev1alpha1.TerraformConfig) string {

	helmValues := fmt.Sprintf(`
        team_name: "%s"
        env_name: %s
        config_name: %s
        module:
            source: %s`, environment.Spec.TeamName,
		environment.Spec.EnvName,
		terraformConfig.ConfigName,
		terraformConfig.Module.Source)

	if terraformConfig.Module.Path != "" {
		helmValues += fmt.Sprintf(`
            path: %s`, terraformConfig.Module.Path)
	}

	if terraformConfig.CronSchedule != "" {
		helmValues += fmt.Sprintf(`
        cron_schedule: %s`, terraformConfig.CronSchedule)
	}

	helmValues += fmt.Sprintf(`
        variables_file:
            source: %s
            path: %s`, terraformConfig.VariablesFile.Source, terraformConfig.VariablesFile.Path)
	return helmValues
}
