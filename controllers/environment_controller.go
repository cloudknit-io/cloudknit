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
	"github.com/google/go-cmp/cmp"
	github2 "github.com/google/go-github/v32/github"
	"go.uber.org/atomic"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"time"

	stablev1 "github.com/compuzest/zlifecycle-il-operator/api/v1"
	"github.com/compuzest/zlifecycle-il-operator/controllers/argocd"
	"github.com/compuzest/zlifecycle-il-operator/controllers/argoworkflow"
	"github.com/compuzest/zlifecycle-il-operator/controllers/terraformgenerator"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/env"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/file"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/github"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/il"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	kClient "sigs.k8s.io/controller-runtime/pkg/client"
)

// EnvironmentReconciler reconciles a Environment object
type EnvironmentReconciler struct {
	kClient.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=stable.compuzest.com,resources=environments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=stable.compuzest.com,resources=environments/status,verbs=get;update;patch

var environmentInitialRun = atomic.NewBool(true)

// Reconcile method called everytime there is a change in Environment Custom Resource
func (r *EnvironmentReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	delayEnvironmentReconcileOnInitialRun(r.Log, 25)
	start := time.Now()
	ctx := context.Background()

	environment := &stablev1.Environment{}

	exists, err := r.tryGetEnvironment(ctx, req, environment)
	if err != nil {
		return ctrl.Result{}, err
	}
	if !exists {
		return ctrl.Result{}, nil
	}

	finalizer := env.Config.EnvironmentFinalizer
	finalizerCompleted, err := r.handleFinalizer(ctx, environment, finalizer)
	if err != nil {
		return ctrl.Result{}, err
	}
	if finalizerCompleted {
		r.Log.Info(
			"Finalizer completed, ending reconcile...",
			"team", environment.Spec.TeamName,
			"environment", environment.Spec.EnvName,
		)
		return ctrl.Result{}, nil
	}

	envDirectory := il.EnvironmentDirectory(environment.Spec.TeamName)
	envComponentDirectory := il.EnvironmentComponentDirectory(environment.Spec.TeamName, environment.Spec.EnvName)
	fileUtil := &file.UtilFileService{}

	isDeleteEvent := !environment.DeletionTimestamp.IsZero() || environment.Spec.Teardown
	if !isDeleteEvent {
		if err := r.updateStatus(ctx, environment); err != nil {
			r.Log.Error(
				err,
				"Error occurred while updating Environment status...",
				"name", environment.Name,
			)
			return ctrl.Result{}, nil
		}

		r.Log.Info(
			"Generating Environment application...",
			"team", environment.Spec.TeamName,
			"environment", environment.Spec.EnvName,
		)
		if err := generateAndSaveEnvironmentApp(fileUtil, environment, envDirectory); err != nil {
			return ctrl.Result{}, err
		}

		githubRepoApi := github.NewHttpRepositoryClient(env.Config.GitHubAuthToken, ctx)
		if err := generateAndSaveEnvironmentComponents(
			r.Log,
			fileUtil,
			environment,
			envComponentDirectory,
			githubRepoApi,
		); err != nil {
			return ctrl.Result{}, err
		}
	} else {
		r.Log.Info(
			"Generating teardown workflow",
			"team", environment.Spec.TeamName,
			"environment", environment.Spec.EnvName,
		)
	}

	r.Log.Info(
		"Generating workflow of workflows...",
		"team", environment.Spec.TeamName,
		"environment", environment.Spec.EnvName,
		"isDeleteEvent", isDeleteEvent,
	)
	if err := generateAndSaveWorkflowOfWorkflows(fileUtil, environment, envComponentDirectory); err != nil {
		return ctrl.Result{}, err
	}

	dirty, err := github.CommitAndPushFiles(
		env.Config.ILRepoSourceOwner,
		env.Config.ILRepoName,
		[]string{envDirectory, envComponentDirectory},
		env.Config.RepoBranch,
		fmt.Sprintf("Reconciling environment %s", environment.Spec.EnvName),
		env.Config.GithubSvcAccntName,
		env.Config.GithubSvcAccntEmail,
	)
	if err != nil {
		return ctrl.Result{}, err
	}
	if dirty {
		r.Log.Info(
			"Committed new changes to IL repo",
			"team", environment.Spec.TeamName,
			"environment", environment.Spec.EnvName,
		)
	} else {
		r.Log.Info(
			"No git changes to commit, no-op reconciliation.",
			"team", environment.Spec.TeamName,
			"environment", environment.Spec.EnvName,
		)
	}

	r.Log.Info(
		"Cleaning up local git files",
		"team", environment.Spec.TeamName,
		"environment", environment.Spec.EnvName,
		"path", envDirectory,
	)
	if err := fileUtil.RemoveAll(envDirectory); err != nil {
		return ctrl.Result{}, err
	}

	duration := time.Since(start)
	r.Log.Info(
		"Reconcile finished",
		"duration", duration,
		"team", environment.Spec.TeamName,
		"environment", environment.Spec.EnvName,
	)

	return ctrl.Result{}, nil
}

func (r *EnvironmentReconciler) tryGetEnvironment(ctx context.Context, req ctrl.Request, e *stablev1.Environment) (exists bool, err error) {
	if err := r.Get(ctx, req.NamespacedName, e); err != nil {
		if errors.IsNotFound(err) {
			r.Log.Info(
				"Environment missing from cache, ending reconcile...",
				"name", req.Name,
				"namespace", req.Namespace,
			)
			return false, nil
		}
		r.Log.Error(
			err,
			"Error occurred while getting Environment...",
			"name", req.Name,
			"namespace", req.Namespace,
		)

		return false, err
	}

	return true, nil
}

// SetupWithManager sets up the Environment Controller with Manager
func (r *EnvironmentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&stablev1.Environment{}).
		Complete(r)
}

func delayEnvironmentReconcileOnInitialRun(log logr.Logger, seconds int64) {
	if environmentInitialRun.Load() {
		log.Info(
			"Delaying Environment reconcile on initial run to wait for Team operator",
			"duration", fmt.Sprintf("%ds", seconds*1000),
		)
		time.Sleep(time.Duration(seconds) * time.Second)
		environmentInitialRun.Store(false)
	}
}

func (r *EnvironmentReconciler) updateStatus(ctx context.Context, e *stablev1.Environment) error {
	isStateDirty := e.Status.TeamName != e.Spec.TeamName ||
		e.Status.EnvName != e.Spec.EnvName ||
		!cmp.Equal(e.Status.EnvironmentComponent, e.Spec.EnvironmentComponent)

	if isStateDirty {
		r.Log.Info(
			"Environment state is dirty and needs to be updated",
			"team", e.Spec.TeamName,
			"environment", e.Spec.EnvName,
		)
		e.Status.TeamName = e.Spec.TeamName
		e.Status.EnvName = e.Spec.EnvName
		e.Status.EnvironmentComponent = e.Spec.EnvironmentComponent
		if err := r.Status().Update(ctx, e); err != nil {
			return err
		}
	} else {
		r.Log.Info(
			"environment state is up-to-date",
			"team", e.Spec.TeamName,
			"environment", e.Spec.EnvName,
		)
	}

	return nil
}

func (r *EnvironmentReconciler) handleFinalizer(ctx context.Context, e *stablev1.Environment, finalizer string) (completed bool, err error) {
	// hack in order to avoid error for stale cached object
	// ref: https://github.com/operator-framework/operator-sdk/issues/3968
	freshEnvironment := &stablev1.Environment{}
	if err := r.Get(ctx, types.NamespacedName{Name: e.Name, Namespace: e.Namespace}, freshEnvironment); err != nil {
		return false, err
	}

	if freshEnvironment.DeletionTimestamp.IsZero() {
		if !common.ContainsString(freshEnvironment.GetFinalizers(), finalizer) {
			r.Log.Info(
				"Setting finalizer for environment",
				"environment", freshEnvironment.Spec.EnvName,
				"team", freshEnvironment.Spec.TeamName,
			)
			freshEnvironment.SetFinalizers(append(freshEnvironment.GetFinalizers(), finalizer))
			if err := r.Update(ctx, freshEnvironment); err != nil {
				return false, err
			}
		}
	} else {
		if common.ContainsString(freshEnvironment.GetFinalizers(), finalizer) {
			if err := r.postDeleteHook(ctx, freshEnvironment); err != nil {
				return false, err
			}

			r.Log.Info(
				"Removing finalizer",
				"team", freshEnvironment.Spec.TeamName,
				"environment", freshEnvironment.Spec.EnvName,
			)
			freshEnvironment.SetFinalizers(common.RemoveString(freshEnvironment.GetFinalizers(), finalizer))
			if err := r.Update(ctx, freshEnvironment); err != nil {
				return false, err
			}

			return true, nil
		}
		return true, nil
	}

	return false, nil
}

func (r *EnvironmentReconciler) postDeleteHook(ctx context.Context, e *stablev1.Environment) error {
	r.Log.Info(
		"Executing post delete hook for environment finalizer",
		"environment", e.Spec.EnvName,
		"team", e.Spec.TeamName,
	)

	argocdApi := argocd.NewHttpClient(r.Log, env.Config.ArgocdServerUrl)
	argoworkflowApi := argoworkflow.NewHttpClient(r.Log, env.Config.ArgoWorkflowsServerUrl)

	if err := r.cleanupIlRepo(ctx, e); err != nil {
		return err
	}
	_ = r.deleteDanglingArgocdApps(e, argocdApi)
	_ = r.deleteDanglingArgoWorkflows(e, argoworkflowApi)
	return nil
}

func (r *EnvironmentReconciler) deleteDanglingArgocdApps(e *stablev1.Environment, argocdApi argocd.Api) error {
	r.Log.Info(
		"Cleaning up dangling argocd apps",
		"team", e.Spec.TeamName,
		"environment", e.Spec.EnvName,
	)
	for _, ec := range e.Spec.EnvironmentComponent {
		appName := fmt.Sprintf("%s-%s-%s", e.Spec.TeamName, e.Spec.EnvName, ec.Name)
		r.Log.Info(
			"Deleting argocd application",
			"team", e.Spec.TeamName,
			"environment", e.Spec.EnvName,
			"component", ec.Name,
			"app", appName,
		)
		if err := argocd.DeleteApplication(r.Log, argocdApi, appName); err != nil {
			r.Log.Error(err, "error deleting argocd app")
		}
	}
	return nil
}

func (r *EnvironmentReconciler) deleteDanglingArgoWorkflows(e *stablev1.Environment, api argoworkflow.Api) error {
	prefix := fmt.Sprintf("%s-%s", e.Spec.TeamName, e.Spec.EnvName)
	namespace := env.Config.ArgoWorkflowsNamespace
	r.Log.Info(
		"Cleaning up dangling Argo Workflows",
		"team", e.Spec.TeamName,
		"environment", e.Spec.EnvName,
		"prefix", prefix,
		"workflowNamespace", namespace,
	)
	if err := argoworkflow.DeleteWorkflowsWithPrefix(r.Log, prefix, namespace, api); err != nil {
		return err
	}

	return nil
}

func (r *EnvironmentReconciler) cleanupIlRepo(ctx context.Context, e *stablev1.Environment) error {
	owner := env.Config.ZlifecycleOwner
	ilRepo := env.Config.ILRepoName
	api := github.NewHttpGitClient(env.Config.GitHubAuthToken, ctx)
	branch := env.Config.RepoBranch
	now := time.Now()
	paths := extractPathsToRemove(*e)
	team := fmt.Sprintf("%s-team-environment", e.Spec.TeamName)
	commitAuthor := &github2.CommitAuthor{Date: &now, Name: &env.Config.GithubSvcAccntName, Email: &env.Config.GithubSvcAccntEmail}
	commitMessage := fmt.Sprintf("Cleaning il objects for %s team in %s environment", e.Spec.TeamName, e.Spec.EnvName)
	r.Log.Info(
		"Cleaning up IL repo",
		"team", e.Spec.TeamName,
		"environment", e.Spec.EnvName,
		"objects", paths,
	)
	if err := github.DeletePatternsFromRootTree(r.Log, api, owner, ilRepo, branch, team, paths, commitAuthor, commitMessage); err != nil {
		return err
	}
	return nil
}

// TODO: Should we remove objects in IL repo?
func extractPathsToRemove(e stablev1.Environment) []string {
	envPath := fmt.Sprintf("%s-environment-component", e.Spec.EnvName)
	envAppPath := fmt.Sprintf("%s-environment.yaml", e.Spec.EnvName)
	return []string{
		envPath,
		envAppPath,
	}
}

func generateAndSaveEnvironmentComponents(
	log logr.Logger,
	fileUtil file.UtilFile,
	environment *stablev1.Environment,
	envComponentDirectory string,
	githubRepoApi github.RepositoryApi,
) error {
	log.Info(
		"Generating environment",
		"name", environment.Name,
		"environment", environment.Spec.EnvName,
		"team", environment.Spec.TeamName,
	)
	for _, ec := range environment.Spec.EnvironmentComponent {
		log.Info(
			"Generating environment component",
			"environment", environment.Spec.EnvName,
			"team", environment.Spec.TeamName,
			"component", ec.Name,
			"type", ec.Type,
		)
		if ec.Variables != nil {
			fileName := fmt.Sprintf("%s.tfvars", ec.Name)

			var variables []*stablev1.Variable
			for _, v := range ec.Variables {
				// TODO: This is a hack to just to make it work, needs to be revisited
				v.Value = fmt.Sprintf("\"%s\"", v.Value)
				variables = append(variables, v)
			}

			if err := fileUtil.SaveVarsToFile(variables, envComponentDirectory, fileName); err != nil {
				return err
			}
		}

		tfvars := ""
		vf := ec.VariablesFile
		if vf != nil {
			tfv, err := getVariablesFromTfvarsFile(
				log,
				githubRepoApi,
				vf.Source,
				env.Config.RepoBranch,
				vf.Path,
			)
			if err != nil {
				return err
			}
			tfvars = tfv
		}

		application := argocd.GenerateEnvironmentComponentApps(*environment, *ec)

		tf := terraformgenerator.TerraformGenerator{Log: log}
		vars := terraformgenerator.TemplateVariables{
			TeamName:             environment.Spec.TeamName,
			EnvName:              environment.Spec.EnvName,
			EnvCompName:          ec.Name,
			EnvCompModulePath:    ec.Module.Path,
			EnvCompModuleSource:  ec.Module.Source,
			EnvCompModuleName:    ec.Module.Name,
			EnvCompOutputs:       ec.Outputs,
			EnvCompDependsOn:     ec.DependsOn,
			EnvCompVariablesFile: tfvars,
			EnvCompVariables:     ec.Variables,
			EnvCompSecrets:       ec.Secrets,
		}
		err := tf.GenerateTerraform(fileUtil, vars, envComponentDirectory)

		if err != nil {
			return err
		}

		if err = fileUtil.SaveYamlFile(*application, envComponentDirectory, ec.Name+".yaml"); err != nil {
			return err
		}
	}

	return nil
}

func getVariablesFromTfvarsFile(log logr.Logger, api github.RepositoryApi, repoUrl string, ref string, path string) (string, error) {
	log.Info("Downloading tfvars file", "repoUrl", repoUrl, "ref", ref, "path", path)
	buff, err := downloadTfvarsFile(log, api, repoUrl, ref, path)
	if err != nil {
		return "", err
	}
	tfvars := string(buff)

	return tfvars, nil
}

func downloadTfvarsFile(log logr.Logger, api github.RepositoryApi, repoUrl string, ref string, path string) ([]byte, error) {
	rc, err := github.DownloadFile(log, api, repoUrl, ref, path)
	if err != nil {
		return nil, err
	}
	defer common.CloseBody(rc)
	buff, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, err
	}
	return buff, nil
}

func generateAndSaveWorkflowOfWorkflows(fileUtil file.UtilFile, environment *stablev1.Environment, envComponentDirectory string) error {

	// WIP, below command is for testing
	// experimentalworkflow := argoWorkflow.GenerateWorkflowOfWorkflows(*environment)
	// if err := fileUtil.SaveYamlFile(*experimentalworkflow, envComponentDirectory, "/experimental_wofw.yaml"); err != nil {
	// 	return err
	// }

	workflow := argoworkflow.GenerateLegacyWorkflowOfWorkflows(*environment)
	if err := fileUtil.SaveYamlFile(*workflow, envComponentDirectory, "/wofw.yaml"); err != nil {
		return err
	}

	return nil
}

func generateAndSaveEnvironmentApp(fileUtil file.UtilFile, environment *stablev1.Environment, envDirectory string) error {
	envApp := argocd.GenerateEnvironmentApp(*environment)
	envYAML := fmt.Sprintf("%s-environment.yaml", environment.Spec.EnvName)

	if err := fileUtil.SaveYamlFile(*envApp, envDirectory, envYAML); err != nil {
		return err
	}

	return nil
}
