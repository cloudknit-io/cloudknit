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
	"go.uber.org/atomic"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/api/errors"
	"time"

	stablev1 "github.com/compuzest/zlifecycle-il-operator/api/v1"
	"github.com/compuzest/zlifecycle-il-operator/controllers/argocd"
	argoWorkflow "github.com/compuzest/zlifecycle-il-operator/controllers/argoworkflow"
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

	if err := r.Get(ctx, req.NamespacedName, environment); err != nil {
		if errors.IsNotFound(err) {
			r.Log.Info(
				"Environment missing from cache, ending reconcile...",
				"name", req.Name,
				"namespace", req.Namespace,
			)
		} else {
			r.Log.Error(
				err,
				"Error occurred while getting Environment...",
				"name", req.Name,
				"namespace", req.Namespace,
			)
		}

		return ctrl.Result{}, nil
	}

	if err := r.updateStatus(ctx, environment); err != nil {
		r.Log.Error(
			err,
			"Error occurred while updating Environment status...",
			"name", environment.Name,
		)
		return ctrl.Result{}, nil
	}

	finalizer := env.Config.GithubFinalizer
	if err := r.handleFinalizer(ctx, environment, finalizer); err != nil {
		return ctrl.Result{}, err
	}

	envDirectory := il.EnvironmentDirectory(environment.Spec.TeamName)
	envComponentDirectory := il.EnvironmentComponentDirectory(environment.Spec.TeamName, environment.Spec.EnvName)
	fileUtil := &file.UtilFileService{}
	isDeleteEvent := !environment.DeletionTimestamp.IsZero()
	if !isDeleteEvent {
		r.Log.Info(
			"Generating Environment application...",
			"team", environment.Spec.TeamName,
			"environment", environment.Spec.EnvName,
			"isDeleteEvent", isDeleteEvent,
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

		// TODO: This will be enabled when we revisit validation
		//if err := r.saveEnvironmentState(ctx, environment, envTrackerCm); err != nil {
		//	return ctrl.Result{}, err
		//}
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

// SetupWithManager sets up the Environment Controller with Manager
func (r *EnvironmentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&stablev1.Environment{}).
		Complete(r)
}

func delayEnvironmentReconcileOnInitialRun(log logr.Logger, seconds int64) {
	if environmentInitialRun.Load() == true {
		log.Info(
			"Delaying Environment reconcile on initial run to wait for Team operator",
			"duration", fmt.Sprintf("%ds", seconds * 1000),
		)
		time.Sleep(time.Duration(seconds) * time.Second)
		environmentInitialRun.Store(false)
	}
}

func (r *EnvironmentReconciler) updateStatus(ctx context.Context, e *stablev1.Environment) error {
	e.Status.TeamName = e.Spec.TeamName
	e.Status.EnvName = e.Spec.EnvName
	e.Status.EnvironmentComponent = e.Spec.EnvironmentComponent
	if err := r.Status().Update(ctx, e); err != nil {
		return err
	}
	return nil
}

func (r *EnvironmentReconciler) handleFinalizer(ctx context.Context, e *stablev1.Environment, finalizer string) error {
	if e.DeletionTimestamp.IsZero() {
		if !common.ContainsString(e.GetFinalizers(), finalizer) {
			r.Log.Info(
				"Setting finalizer for environment",
				"env", e.Spec.EnvName,
				"team", e.Spec.TeamName,
			)
			e.SetFinalizers(append(e.GetFinalizers(), finalizer))
			if err := r.Update(ctx, e); err != nil {
				return err
			}
		}
	} else {
		if common.ContainsString(e.GetFinalizers(), finalizer) {
			if err := r.postDeleteHook(e); err != nil {
				return err
			}

			e.SetFinalizers(common.RemoveString(e.GetFinalizers(), finalizer))
			if err := r.Update(ctx, e); err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *EnvironmentReconciler) postDeleteHook(e *stablev1.Environment) error {
	r.Log.Info(
		"Executing post delete hook for environment finalizer",
		"environment", e.Spec.EnvName,
		"team", e.Spec.TeamName,
	)
	return nil
}

// TODO: Should we remove objects in IL repo?
//func extractPathsToRemove(e stablev1.Environment) []string {
//	envPath    := fmt.Sprintf("%s-environment-component", e.Spec.EnvName)
//	envAppPath := fmt.Sprintf("%s-environment.yaml", e.Spec.EnvName)
//	return []string{
//		envPath,
//		envAppPath,
//	}
//}

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
		"env", environment.Spec.EnvName,
		"team", environment.Spec.TeamName,
	)
	for _, ec := range environment.Spec.EnvironmentComponent {
		log.Info("Generating environment component", "component", ec.Name, "type", ec.Type)
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
	defer rc.Close()
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

	workflow := argoWorkflow.GenerateLegacyWorkflowOfWorkflows(*environment)
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
