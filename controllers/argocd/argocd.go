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

	appv1 "github.com/argoproj/argo-cd/pkg/apis/application/v1alpha1"
	stablev1alpha1 "github.com/compuzest/zlifecycle-il-operator/api/v1alpha1"
	env "github.com/compuzest/zlifecycle-il-operator/controllers/util/env"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/il"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GenerateCompanyApp(company stablev1alpha1.Company) *appv1.Application {
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
func GenerateTeamApp(team stablev1alpha1.Team) *appv1.Application {
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
				Path:           "./" + il.EnvironmentDirectory(team.Spec.TeamName),
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
				RepoURL:        env.Config.ILRepoURL,
				Path:           "./" + il.EnvironmentComponentDirectory(environment.Spec.TeamName, environment.Spec.EnvName),
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

func getHelmValues(environment stablev1alpha1.Environment, terraformConfig stablev1alpha1.TerraformConfig) string {

	helmValues := fmt.Sprintf(`
        team_name: "%s"
        env_name: %s
        config_name: %s
        module:
            source: %s`, environment.Spec.TeamName,
		environment.Spec.EnvName,
		terraformConfig.ConfigName,
		terraformConfigModel.buildModuleSource(terraformConfig.Module.Source))

	if terraformConfig.Module.Path != "" {
		helmValues += fmt.Sprintf(`
            path: %s`, terraformConfigModel.buildModulePath(terraformConfig.Module.Path))
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

func GenerateTeamConfigWatcherApp(team stablev1alpha1.Team) *appv1.Application {

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
