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
	"github.com/go-logr/logr"
	"strings"

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
				"type":                 "project",
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
				"type":                 "environment",
				"project_id":           environment.Spec.TeamName,
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

func GenerateEnvironmentComponentApps(environment stablev1alpha1.Environment, environmentComponent stablev1alpha1.EnvironmentComponent) *appv1.Application {

	helmValues := getHelmValues(environment, environmentComponent)

	return &appv1.Application{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Application",
			APIVersion: "argoproj.io/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      environment.Spec.TeamName + "-" + environment.Spec.EnvName + "-" + environmentComponent.Name,
			Namespace: "argocd",
			Labels: map[string]string{
				"zlifecycle.com/model": "environment-component",
				"component_type":       environmentComponent.Type,
				"type":                 "config",
				"project_id":           environment.Spec.TeamName,
				"environment_id":       environment.Spec.TeamName + "-" + environment.Spec.EnvName,
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

func getHelmValues(environment stablev1alpha1.Environment, environmentComponent stablev1alpha1.EnvironmentComponent) string {

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

	helmValues += fmt.Sprintf(`
        variables_file:
            source: %s
            path: %s`, environmentComponent.VariablesFile.Source, environmentComponent.VariablesFile.Path)
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

func RegisterRepo(log logr.Logger, api Api, repoOpts RepoOpts) (bool, error) {
	repoUri := repoOpts.RepoUrl[strings.LastIndex(repoOpts.RepoUrl, "/")+1:]
	repoName := strings.TrimSuffix(repoUri, ".git")

	tokenResponse, err := api.GetAuthToken()
	if err != nil {
		log.Error(err, "Error while calling get auth token")
		return false, err
	}

	bearer := "Bearer " + tokenResponse.Token
	repositories, resp1, err1 := api.ListRepositories(bearer)
	if err1 != nil {
		return false, err1
	}
	defer resp1.Body.Close()

	if isRepoRegistered(*repositories, repoOpts.RepoUrl) {
		log.Info("Repository already registered",
			"repos", repositories.Items,
			"repoName", repoName,
			"repoUrl", repoOpts.RepoUrl,
		)
		return false, nil
	}

	log.Info("Repository is not registered, registering now...", "repos", repositories, "repoName", repoName)

	createRepoBody := CreateRepoBody{Repo: repoOpts.RepoUrl, Name: repoName, SshPrivateKey: repoOpts.SshPrivateKey}
	resp2, err2 := api.CreateRepository(createRepoBody, bearer)
	if err2 != nil {
		log.Error(err2, "Error while calling post create repository")
		return false, err2
	}
	defer resp2.Body.Close()

	log.Info("Successfully registered repository", "repo", repoOpts.RepoUrl)

	return true, nil
}
