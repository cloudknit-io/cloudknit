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

package argocd

import (
	"fmt"
	"strings"

	appv1 "github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	stablev1 "github.com/compuzest/zlifecycle-il-operator/api/v1"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/env"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/il"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GenerateCompanyApp(company *stablev1.Company) *appv1.Application {
	return &appv1.Application{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "argoproj.io/v1alpha1",
			Kind:       "Application",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      company.Spec.CompanyName,
			Namespace: "argocd",
			Labels: map[string]string{
				"zlifecycle.com/model": "company",
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
				RepoURL:        env.Config.ILRepoURL,
				Path:           "./" + il.Config.TeamDirectory,
				TargetRevision: "HEAD",
			},
		},
		Status: appv1.ApplicationStatus{
			Sync: appv1.SyncStatus{
				ComparedTo: appv1.ComparedTo{
					Source: appv1.ApplicationSource{
						RepoURL: env.Config.ILRepoURL,
					},
				},
				Status: "Synced",
			},
		},
	}
}

func GenerateTeamApp(team *stablev1.Team) *appv1.Application {
	return &appv1.Application{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "argoproj.io/v1alpha1",
			Kind:       "Application",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      team.Spec.TeamName,
			Namespace: "argocd",
			Labels: map[string]string{
				"zlifecycle.com/model": "team",
				"type":                 "project",
			},
		},
		Spec: appv1.ApplicationSpec{
			Project: team.Spec.TeamName,
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
				RepoURL:        env.Config.ILRepoURL,
				Path:           "./" + il.TeamDirectoryPath(team.Spec.TeamName),
				TargetRevision: "HEAD",
			},
		},
		Status: appv1.ApplicationStatus{
			Sync: appv1.SyncStatus{
				ComparedTo: appv1.ComparedTo{
					Source: appv1.ApplicationSource{
						RepoURL: env.Config.ILRepoURL,
					},
				},
				Status: "Synced",
			},
		},
	}
}

func GenerateEnvironmentApp(environment *stablev1.Environment) *appv1.Application {
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
				"type":                 "environment",
				"env_name":             environment.Spec.EnvName,
				"project_id":           environment.Spec.TeamName,
			},
		},
		Spec: appv1.ApplicationSpec{
			Project: environment.Spec.TeamName,
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
				RepoURL:        env.Config.ILRepoURL,
				Path:           "./" + il.EnvironmentDirectoryPath(environment.Spec.TeamName, environment.Spec.EnvName),
				TargetRevision: "HEAD",
			},
		},
		Status: appv1.ApplicationStatus{
			Sync: appv1.SyncStatus{
				ComparedTo: appv1.ComparedTo{
					Source: appv1.ApplicationSource{
						RepoURL: env.Config.ILRepoURL,
					},
				},
				Status: "Synced",
			},
		},
	}
}

func GenerateEnvironmentComponentApps(environment *stablev1.Environment, environmentComponent *stablev1.EnvironmentComponent) *appv1.Application {
	helmValues := getHelmValues(environment, environmentComponent)
	labels := map[string]string{
		"zlifecycle.com/model": "environment-component",
		"component_type":       environmentComponent.Type,
		"type":                 "config",
		"component_name":       environmentComponent.Name,
		"project_id":           environment.Spec.TeamName,
		"environment_id":       environment.Spec.TeamName + "-" + environment.Spec.EnvName,
		"depends_on":           strings.Join(environmentComponent.DependsOn, ".."),
	}
	for _, tag := range environmentComponent.Tags {
		labels[tag.Name] = tag.Value
	}
	return &appv1.Application{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Application",
			APIVersion: "argoproj.io/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      environment.Spec.TeamName + "-" + environment.Spec.EnvName + "-" + environmentComponent.Name,
			Namespace: "argocd",
			Labels:    labels,
			Finalizers: []string{
				"resources-finalizer.argocd.argoproj.io",
			},
		},
		Spec: appv1.ApplicationSpec{
			Project: environment.Spec.TeamName,
			SyncPolicy: &appv1.SyncPolicy{
				Automated: &appv1.SyncPolicyAutomated{
					Prune: true,
				},
			},
			Destination: appv1.ApplicationDestination{
				Server:    env.Config.K8sAPIURL,
				Namespace: "default",
			},
			Source: appv1.ApplicationSource{
				RepoURL:        env.Config.HelmChartsRepo,
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
						RepoURL: env.Config.HelmChartsRepo,
					},
				},
				Status: "Synced",
			},
		},
	}
}

func getHelmValues(environment *stablev1.Environment, environmentComponent *stablev1.EnvironmentComponent) string {
	helmValues := fmt.Sprintf(`
        team_name: "%s"
        env_name: %s
        config_name: %s
        module:
            source: %s
            path: %s`, environment.Spec.TeamName,
		environment.Spec.EnvName,
		environmentComponent.Name,
		il.EnvComponentModuleSource(environmentComponent.Module.Source, environmentComponent.Module.Name),
		il.EnvComponentModulePath(environmentComponent.Module.Path))
	if environmentComponent.CronSchedule != "" {
		helmValues += fmt.Sprintf(`
        cron_schedule: "%s"`, environmentComponent.CronSchedule)
	}
	if environmentComponent.VariablesFile != nil {
		helmValues += fmt.Sprintf(`
        variables_file:
            source: %s
            path: %s`, environmentComponent.VariablesFile.Source, environmentComponent.VariablesFile.Path)
	}
	return helmValues
}

func GenerateTeamConfigWatcherApp(team *stablev1.Team) *appv1.Application {
	return &appv1.Application{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "argoproj.io/v1alpha1",
			Kind:       "Application",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      team.Spec.TeamName + "-team-watcher",
			Namespace: "argocd",
			Labels: map[string]string{
				"zlifecycle.com/model":                   "config-watcher",
				"zlifecycle.com/watched-custom-resource": "team",
			},
		},
		Spec: appv1.ApplicationSpec{
			Project: team.Spec.TeamName,
			SyncPolicy: &appv1.SyncPolicy{
				Automated: &appv1.SyncPolicyAutomated{
					Prune: true,
				},
				Retry: &appv1.RetryStrategy{Limit: 1},
			},
			Destination: appv1.ApplicationDestination{
				Server:    "https://kubernetes.default.svc",
				Namespace: "default",
			},
			Source: appv1.ApplicationSource{
				RepoURL:        team.Spec.ConfigRepo.Source,
				Path:           team.Spec.ConfigRepo.Path,
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
						RepoURL: team.Spec.ConfigRepo.Source,
					},
				},
				Status: "Synced",
			},
		},
	}
}

func GenerateCompanyConfigWatcherApp(customerName string, companyConfigRepo string) *appv1.Application {
	return &appv1.Application{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "argoproj.io/v1alpha1",
			Kind:       "Application",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      customerName + "-watcher",
			Namespace: "argocd",
			Labels: map[string]string{
				"zlifecycle.com/model":                   "config-watcher",
				"zlifecycle.com/watched-custom-resource": "company",
			},
		},
		Spec: appv1.ApplicationSpec{
			Project: "default",
			SyncPolicy: &appv1.SyncPolicy{
				Automated: &appv1.SyncPolicyAutomated{
					Prune: true,
				},
				Retry: &appv1.RetryStrategy{Limit: 1},
			},
			Destination: appv1.ApplicationDestination{
				Server:    "https://kubernetes.default.svc",
				Namespace: "default",
			},
			Source: appv1.ApplicationSource{
				RepoURL:        companyConfigRepo,
				Path:           ".",
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
						RepoURL: companyConfigRepo,
					},
				},
				Status: "Synced",
			},
		},
	}
}

func GenerateCompanyBootstrapApp() *appv1.Application {
	return &appv1.Application{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "argoproj.io/v1alpha1",
			Kind:       "Application",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "company-bootstrap",
			Namespace: "argocd",
			Labels: map[string]string{
				"zlifecycle.com/model":                   "bootstrap",
				"zlifecycle.com/watched-custom-resource": "company",
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
				RepoURL:        env.Config.ILRepoURL,
				Path:           "company",
				TargetRevision: "HEAD",
			},
		},
	}
}

func GenerateConfigWatcherBootstrapApp() *appv1.Application {
	return &appv1.Application{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "argoproj.io/v1alpha1",
			Kind:       "Application",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "config-watcher-bootstrap",
			Namespace: "argocd",
			Labels: map[string]string{
				"zlifecycle.com/model":                   "bootstrap",
				"zlifecycle.com/watched-custom-resource": "config-watcher",
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
				RepoURL:        env.Config.ILRepoURL,
				Path:           "config-watcher",
				TargetRevision: "HEAD",
			},
		},
	}
}
