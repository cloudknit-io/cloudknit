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
	"context"
	"fmt"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/common"
	github2 "github.com/google/go-github/v32/github"
	"time"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	kClient "sigs.k8s.io/controller-runtime/pkg/client"

	stablev1alpha1 "github.com/compuzest/zlifecycle-il-operator/api/v1alpha1"
	argocd "github.com/compuzest/zlifecycle-il-operator/controllers/argocd"
	argoWorkflow "github.com/compuzest/zlifecycle-il-operator/controllers/argoworkflow"
	k8s "github.com/compuzest/zlifecycle-il-operator/controllers/kubernetes"
	env "github.com/compuzest/zlifecycle-il-operator/controllers/util/env"
	file "github.com/compuzest/zlifecycle-il-operator/controllers/util/file"
	github "github.com/compuzest/zlifecycle-il-operator/controllers/util/github"
	il "github.com/compuzest/zlifecycle-il-operator/controllers/util/il"
)

// EnvironmentReconciler reconciles a Environment object
type EnvironmentReconciler struct {
	kClient.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=stable.compuzest.com,resources=environments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=stable.compuzest.com,resources=environments/status,verbs=get;update;patch

// Reconcile method called everytime there is a change in Environment Custom Resource
func (r *EnvironmentReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	//log := r.Log.WithValues("environment", req.NamespacedName)

	environment := &stablev1alpha1.Environment{}

	if err := r.Get(ctx, req.NamespacedName, environment); err != nil {
		return ctrl.Result{}, nil
	}

	finalizer := env.Config.GithubFinalizer

	if environment.DeletionTimestamp.IsZero() {
		if !common.ContainsString(environment.GetFinalizers(), finalizer) {
			environment.SetFinalizers(append(environment.GetFinalizers(), finalizer))
			if err := r.Update(ctx, environment); err != nil {
				return ctrl.Result{}, err
			}
		}
	} else {
		if common.ContainsString(environment.GetFinalizers(), finalizer) {
			if err := r.deleteExternalResources(ctx, environment); err != nil {
				return ctrl.Result{}, err
			}

			environment.SetFinalizers(common.RemoveString(environment.GetFinalizers(), finalizer))
			if err := r.Update(ctx, environment); err != nil {
				return ctrl.Result{}, err
			}
		}

		return ctrl.Result{}, nil
	}

	envDirectory := il.EnvironmentDirectory(environment.Spec.TeamName)
	envComponentDirectory := il.EnvironmentComponentDirectory(environment.Spec.TeamName, environment.Spec.EnvName)

	if err := generateAndSaveEnvironmentApp(environment, envDirectory); err != nil {
		return ctrl.Result{}, err
	}

	if err := generateAndSaveEnvironmentComponents(environment, envComponentDirectory); err != nil {
		return ctrl.Result{}, err
	}

	if err := generateAndSaveWorkflowOfWorkflows(environment, envComponentDirectory); err != nil {
		return ctrl.Result{}, err
	}

	if err := generateAndSavePresyncJob(environment, envComponentDirectory); err != nil {
		return ctrl.Result{}, err
	}

	// Avoid race condition on initial Reconcile, collides with Team controller commit
	time.Sleep(10 * time.Second)

	if err := github.CommitAndPushFiles(
		env.Config.ILRepoSourceOwner,
		env.Config.ILRepoName,
		[]string{envDirectory, envComponentDirectory},
		env.Config.RepoBranch,
		fmt.Sprintf("Reconciling environment %s", environment.Spec.EnvName),
		env.Config.GithubSvcAccntName,
		env.Config.GithubSvcAccntEmail); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *EnvironmentReconciler) deleteExternalResources(ctx context.Context, e *stablev1alpha1.Environment) error {
	owner         := env.Config.ZlifecycleOwner
	ilRepo        := env.Config.ILRepoName
	api           := github.NewHttpGitClient(env.Config.GitHubAuthToken, ctx)
	branch        := env.Config.RepoBranch
	now           := time.Now()
	paths         := extractPathsToRemove(*e)
	team          := fmt.Sprintf("%s-team-environment", e.Spec.TeamName)
	commitAuthor  := &github2.CommitAuthor{Date: &now, Name: &env.Config.GithubSvcAccntName, Email: &env.Config.GithubSvcAccntEmail}
	commitMessage := fmt.Sprintf("Cleaning il objects for %s team in %s environment", e.Spec.TeamName, e.Spec.EnvName)
	if err := github.DeletePatternsFromRootTree(r.Log, api, owner, ilRepo, branch, team, paths, commitAuthor, commitMessage); err != nil {
		return err
	}
	return nil
}

func extractPathsToRemove(e stablev1alpha1.Environment) []string {
	envPath    := fmt.Sprintf("%s-environment-component", e.Spec.EnvName)
	envAppPath := fmt.Sprintf("%s-environment.yaml", e.Spec.EnvName)
	return []string{
		envPath,
		envAppPath,
	}
}

// SetupWithManager sets up the Environment Controller with Manager
func (r *EnvironmentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&stablev1alpha1.Environment{}).
		Complete(r)
}

func generateAndSaveEnvironmentComponents(environment *stablev1alpha1.Environment, envComponentDirectory string) error {
	for _, environmentComponent := range environment.Spec.EnvironmentComponent {
		if environmentComponent.Variables != nil {
			fileName := fmt.Sprintf("%s.tfvars", environmentComponent.Name)

			if err := file.SaveVarsToFile(environmentComponent.Variables, envComponentDirectory, fileName); err != nil {
				return err
			}

			environmentComponent.VariablesFile = &stablev1alpha1.VariablesFile{
				Source: env.Config.ILRepoURL,
				Path:   envComponentDirectory + "/" + fileName,
			}
		}

		application := argocd.GenerateEnvironmentComponentApps(*environment, *environmentComponent)

		if err := file.SaveYamlFile(*application, envComponentDirectory, environmentComponent.Name+".yaml"); err != nil {
			return err
		}
	}

	return nil
}

func generateAndSaveWorkflowOfWorkflows(environment *stablev1alpha1.Environment, envComponentDirectory string) error {
	workflow := argoWorkflow.GenerateWorkflowOfWorkflows(*environment)
	if err := file.SaveYamlFile(*workflow, envComponentDirectory, "/wofw.yaml"); err != nil {
		return err
	}

	return nil
}

func generateAndSaveEnvironmentApp(environment *stablev1alpha1.Environment, envDirectory string) error {
	envApp := argocd.GenerateEnvironmentApp(*environment)
	envYAML := fmt.Sprintf("%s-environment.yaml", environment.Spec.EnvName)

	if err := file.SaveYamlFile(*envApp, envDirectory, envYAML); err != nil {
		return err
	}

	return nil
}

func generateAndSavePresyncJob(environment *stablev1alpha1.Environment, envComponentDirectory string) error {
	presyncJob := k8s.GeneratePreSyncJob(*environment)
	if err := file.SaveYamlFile(*presyncJob, envComponentDirectory, "/presync-job.yaml"); err != nil {
		return err
	}

	return nil
}
