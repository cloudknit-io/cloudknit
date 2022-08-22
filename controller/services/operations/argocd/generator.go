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

	"github.com/compuzest/zlifecycle-il-operator/controller/util"

	"github.com/compuzest/zlifecycle-il-operator/controller/codegen/il"

	appv1 "github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	stablev1 "github.com/compuzest/zlifecycle-il-operator/api/v1"
	"github.com/compuzest/zlifecycle-il-operator/controller/env"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const localCluster = "https://kubernetes.default.svc"

func CompanyProjectName() string {
	return fmt.Sprintf("project-%s", env.Config.CompanyName)
}

func newTypeMeta() metav1.TypeMeta {
	return metav1.TypeMeta{
		APIVersion: "argoproj.io/v1alpha1",
		Kind:       "Application",
	}
}

func newApplicationDestination(server, namespace string) appv1.ApplicationDestination {
	return appv1.ApplicationDestination{
		Server:    server,
		Namespace: namespace,
	}
}

func newApplicationSource(repoURL, path string, recurse bool) appv1.ApplicationSource {
	repoURL = util.RewriteGitHubURLToHTTPS(repoURL, false)
	as := appv1.ApplicationSource{
		RepoURL:        repoURL,
		Path:           path,
		TargetRevision: "HEAD",
	}
	if recurse {
		as.Directory = &appv1.ApplicationSourceDirectory{Recurse: true}
	}
	return as
}

func newApplicationStatus(repoURL string) appv1.ApplicationStatus {
	return appv1.ApplicationStatus{
		Sync: appv1.SyncStatus{
			ComparedTo: appv1.ComparedTo{
				Source: appv1.ApplicationSource{
					RepoURL: util.RewriteGitHubURLToHTTPS(repoURL, false),
				},
			},
			Status: "Synced",
		},
	}
}

func GenerateCompanyApp(company *stablev1.Company) *appv1.Application {
	return &appv1.Application{
		TypeMeta: newTypeMeta(),
		ObjectMeta: metav1.ObjectMeta{
			Name:      company.Spec.CompanyName,
			Namespace: env.ArgocdNamespace(),
			Labels: map[string]string{
				"zlifecycle.com/model": "company",
			},
		},
		Spec: appv1.ApplicationSpec{
			Project: CompanyProjectName(),
			SyncPolicy: &appv1.SyncPolicy{
				Automated: &appv1.SyncPolicyAutomated{
					Prune: true,
				},
			},
			Destination: newApplicationDestination(localCluster, "default"),
			Source:      newApplicationSource(env.Config.ILZLifecycleRepositoryURL, "./"+il.Config.TeamDirectory, false),
		},
		Status: newApplicationStatus(env.Config.ILZLifecycleRepositoryURL),
	}
}

func GenerateTeamApp(team *stablev1.Team) *appv1.Application {
	return &appv1.Application{
		TypeMeta: newTypeMeta(),
		ObjectMeta: metav1.ObjectMeta{
			Name:      team.Spec.TeamName,
			Namespace: env.ArgocdNamespace(),
			Labels: map[string]string{
				"zlifecycle.com/model": "team",
				"type":                 "project",
			},
		},
		Spec: appv1.ApplicationSpec{
			Project: CompanyProjectName(),
			SyncPolicy: &appv1.SyncPolicy{
				Automated: &appv1.SyncPolicyAutomated{
					Prune: true,
				},
			},
			Destination: newApplicationDestination(localCluster, "default"),
			Source:      newApplicationSource(env.Config.ILZLifecycleRepositoryURL, "./"+il.EnvironmentDirectoryPath(team.Spec.TeamName), false),
		},
		Status: newApplicationStatus(env.Config.ILZLifecycleRepositoryURL),
	}
}

func GenerateEnvironmentApp(environment *stablev1.Environment) *appv1.Application {
	return &appv1.Application{
		TypeMeta: newTypeMeta(),
		ObjectMeta: metav1.ObjectMeta{
			Name:      environment.Spec.TeamName + "-" + environment.Spec.EnvName,
			Namespace: env.ArgocdNamespace(),
			Labels: map[string]string{
				"zlifecycle.com/model": "environment",
				"type":                 "environment",
				"env_name":             environment.Spec.EnvName,
				"project_id":           environment.Spec.TeamName,
			},
		},
		Spec: appv1.ApplicationSpec{
			Project: CompanyProjectName(),
			SyncPolicy: &appv1.SyncPolicy{
				Automated: &appv1.SyncPolicyAutomated{
					Prune:    true,
					SelfHeal: true,
				},
			},
			Destination: newApplicationDestination(localCluster, "default"),
			Source: newApplicationSource(
				env.Config.ILZLifecycleRepositoryURL,
				"./"+il.EnvironmentComponentsDirectoryPath(environment.Spec.TeamName, environment.Spec.EnvName),
				false,
			),
		},
		Status: newApplicationStatus(env.Config.ILZLifecycleRepositoryURL),
	}
}

func GenerateEnvironmentComponentApps(e *stablev1.Environment, ec *stablev1.EnvironmentComponent) *appv1.Application {
	helmValues := getHelmValues(e, ec)

	labels := map[string]string{
		"zlifecycle.com/model": "environment-component",
		"component_type":       ec.Type,
		"type":                 "config",
		"component_name":       ec.Name,
		"project_id":           e.Spec.TeamName,
		"environment_id":       fmt.Sprintf("%s-%s", e.Spec.TeamName, e.Spec.EnvName),
	}

	for i, dependsOn := range ec.DependsOn {
		labels[fmt.Sprintf("depends_on_%d", i)] = dependsOn
	}

	for _, tag := range ec.Tags {
		labels[tag.Name] = tag.Value
	}

	source := appv1.ApplicationSource{
		RepoURL:        util.RewriteGitHubURLToHTTPS(env.Config.GitHelmChartsRepository, false),
		Path:           "charts/terraform-config",
		TargetRevision: "HEAD",
		Helm: &appv1.ApplicationSourceHelm{
			Values: helmValues,
		},
	}
	if ec.Type == "argocd" {
		source = appv1.ApplicationSource{
			RepoURL:        util.RewriteGitHubURLToHTTPS(env.Config.ILZLifecycleRepositoryURL, false),
			Path:           il.EnvironmentComponentArgocdAppsDirectoryPath(e.Spec.TeamName, e.Spec.EnvName, ec.Name),
			TargetRevision: "HEAD",
		}
	}

	return &appv1.Application{
		TypeMeta: newTypeMeta(),
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-%s-%s", e.Spec.TeamName, e.Spec.EnvName, ec.Name),
			Namespace: env.ArgocdNamespace(),
			Labels:    labels,
			Finalizers: []string{
				"resources-finalizer.argocd.argoproj.io",
			},
		},
		Spec: appv1.ApplicationSpec{
			Project: CompanyProjectName(),
			SyncPolicy: &appv1.SyncPolicy{
				Automated: &appv1.SyncPolicyAutomated{
					Prune: true,
				},
			},
			Destination: newApplicationDestination(localCluster, "default"),
			Source:      source,
		},
		Status: newApplicationStatus(env.Config.GitHelmChartsRepository),
	}
}

func getHelmValues(environment *stablev1.Environment, environmentComponent *stablev1.EnvironmentComponent) string {
	helmValues := fmt.Sprintf(`
        namespace: "%s"
        team_name: "%s"
        env_name: %s
        config_name: %s
        module:
            source: %s
            path: %s`, env.ArgocdNamespace(),
		environment.Spec.TeamName,
		environment.Spec.EnvName,
		environmentComponent.Name,
		il.EnvironmentComponentModuleSource(environmentComponent.Module.Source, environmentComponent.Module.Name),
		il.EnvironmentComponentModulePath(environmentComponent.Module.Path))
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
		TypeMeta: newTypeMeta(),
		ObjectMeta: metav1.ObjectMeta{
			Name:      team.Spec.TeamName + "-team-watcher",
			Namespace: env.ArgocdNamespace(),
			Labels: map[string]string{
				"zlifecycle.com/model":                   "config-watcher",
				"zlifecycle.com/watched-custom-resource": "team",
			},
		},
		Spec: appv1.ApplicationSpec{
			Project: CompanyProjectName(),
			SyncPolicy: &appv1.SyncPolicy{
				Automated: &appv1.SyncPolicyAutomated{
					Prune: true,
				},
				Retry: &appv1.RetryStrategy{Limit: 1},
			},
			Destination: newApplicationDestination(localCluster, "default"),
			Source:      newApplicationSource(team.Spec.ConfigRepo.Source, team.Spec.ConfigRepo.Path, true),
		},
		Status: newApplicationStatus(team.Spec.ConfigRepo.Source),
	}
}

func GenerateCompanyConfigWatcherApp(customerName string, companyConfigRepo string, companyConfigRepoPath string) *appv1.Application {
	return &appv1.Application{
		TypeMeta: newTypeMeta(),
		ObjectMeta: metav1.ObjectMeta{
			Name:      customerName + "-watcher",
			Namespace: env.ArgocdNamespace(),
			Labels: map[string]string{
				"zlifecycle.com/model":                   "config-watcher",
				"zlifecycle.com/watched-custom-resource": "company",
			},
		},
		Spec: appv1.ApplicationSpec{
			Project: CompanyProjectName(),
			SyncPolicy: &appv1.SyncPolicy{
				Automated: &appv1.SyncPolicyAutomated{
					Prune: true,
				},
				Retry: &appv1.RetryStrategy{Limit: 1},
			},
			Destination: newApplicationDestination(localCluster, "default"),
			Source:      newApplicationSource(companyConfigRepo, companyConfigRepoPath, true),
		},
		Status: newApplicationStatus(companyConfigRepo),
	}
}

func GenerateCompanyBootstrapApp() *appv1.Application {
	return &appv1.Application{
		TypeMeta: newTypeMeta(),
		ObjectMeta: metav1.ObjectMeta{
			Name:      "company-bootstrap",
			Namespace: env.ArgocdNamespace(),
			Labels: map[string]string{
				"zlifecycle.com/model":                   "bootstrap",
				"zlifecycle.com/watched-custom-resource": "company",
			},
		},
		Spec: appv1.ApplicationSpec{
			Project: CompanyProjectName(),
			SyncPolicy: &appv1.SyncPolicy{
				Automated: &appv1.SyncPolicyAutomated{
					Prune: true,
				},
			},
			Destination: newApplicationDestination(localCluster, "default"),
			Source:      newApplicationSource(env.Config.ILZLifecycleRepositoryURL, "company", false),
		},
	}
}

func GenerateConfigWatcherBootstrapApp() *appv1.Application {
	return &appv1.Application{
		TypeMeta: newTypeMeta(),
		ObjectMeta: metav1.ObjectMeta{
			Name:      "config-watcher-bootstrap",
			Namespace: env.ArgocdNamespace(),
			Labels: map[string]string{
				"zlifecycle.com/model":                   "bootstrap",
				"zlifecycle.com/watched-custom-resource": "config-watcher",
			},
		},
		Spec: appv1.ApplicationSpec{
			Project: CompanyProjectName(),
			SyncPolicy: &appv1.SyncPolicy{
				Automated: &appv1.SyncPolicyAutomated{
					Prune: true,
				},
			},
			Destination: newApplicationDestination(localCluster, "default"),
			Source:      newApplicationSource(env.Config.ILZLifecycleRepositoryURL, "config-watcher", false),
		},
	}
}

func AddLabelsToCustomerApp(app *appv1.Application, e *stablev1.Environment, ec *stablev1.EnvironmentComponent, filename string) {
	app.Labels = map[string]string{
		"zlifecycle.com/model": "environment-component",
		"component_type":       "argocd",
		"component_name":       ec.Name,
		"project_id":           e.Spec.TeamName,
		"environment_id":       fmt.Sprintf("%s-%s", e.Spec.TeamName, e.Spec.EnvName),
		"source_file_name":     filename,
	}
	for i, dep := range ec.DependsOn {
		label := fmt.Sprintf("depends_on_%d", i)
		app.Labels[label] = dep
	}
	for _, tag := range ec.Tags {
		app.Labels[tag.Name] = tag.Value
	}
}
