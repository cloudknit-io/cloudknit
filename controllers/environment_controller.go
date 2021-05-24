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
	"github.com/magiconair/properties"
	"io/ioutil"
	"time"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	kClient "sigs.k8s.io/controller-runtime/pkg/client"

	stablev1alpha1 "github.com/compuzest/zlifecycle-il-operator/api/v1alpha1"
	argocd "github.com/compuzest/zlifecycle-il-operator/controllers/argocd"
	argoWorkflow "github.com/compuzest/zlifecycle-il-operator/controllers/argoworkflow"
	terraformgenerator "github.com/compuzest/zlifecycle-il-operator/controllers/terraformgenerator"
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

	finalizer     := env.Config.GithubFinalizer
	githubGitApi  := github.NewHttpGitClient(env.Config.GitHubAuthToken, ctx)
	githubRepoApi := github.NewHttpRepositoryClient(env.Config.GitHubAuthToken, ctx)
	if environment.DeletionTimestamp.IsZero() {
		if !common.ContainsString(environment.GetFinalizers(), finalizer) {
			environment.SetFinalizers(append(environment.GetFinalizers(), finalizer))
			if err := r.Update(ctx, environment); err != nil {
				return ctrl.Result{}, err
			}
		}
	} else {
		if common.ContainsString(environment.GetFinalizers(), finalizer) {
			if err := r.deleteExternalResources(ctx, environment, githubGitApi); err != nil {
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
	var fileUtil file.UtilFile = &file.UtilFileService{}

	if err := generateAndSaveEnvironmentApp(fileUtil, environment, envDirectory); err != nil {
		return ctrl.Result{}, err
	}

	if err := generateAndSaveEnvironmentComponents(
		r.Log,
		fileUtil,
		environment,
		envComponentDirectory,
		githubRepoApi,
		); err != nil {
		return ctrl.Result{}, err
	}

	if err := generateAndSaveWorkflowOfWorkflows(fileUtil, environment, envComponentDirectory); err != nil {
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

func (r *EnvironmentReconciler) deleteExternalResources(ctx context.Context, e *stablev1alpha1.Environment, githubGitApi github.GitApi) error {
	owner         := env.Config.ZlifecycleOwner
	ilRepo        := env.Config.ILRepoName
	branch        := env.Config.RepoBranch
	now           := time.Now()
	paths         := extractPathsToRemove(*e)
	team          := fmt.Sprintf("%s-team-environment", e.Spec.TeamName)
	commitAuthor  := &github2.CommitAuthor{Date: &now, Name: &env.Config.GithubSvcAccntName, Email: &env.Config.GithubSvcAccntEmail}
	commitMessage := fmt.Sprintf("Cleaning il objects for %s team in %s environment", e.Spec.TeamName, e.Spec.EnvName)
	if err := github.DeletePatternsFromRootTree(r.Log, githubGitApi, owner, ilRepo, branch, team, paths, commitAuthor, commitMessage); err != nil {
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

func generateAndSaveEnvironmentComponents(
	log logr.Logger,
	fileUtil file.UtilFile,
	environment *stablev1alpha1.Environment,
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

			var variables []*stablev1alpha1.Variable
			for _, v := range ec.Variables {
				// TODO: This is a hack to just to make it work, needs to be revisited
				v.Value = fmt.Sprintf("\"%s\"", v.Value)
				variables = append(variables, v)
			}

			if err := fileUtil.SaveVarsToFile(variables, envComponentDirectory, fileName); err != nil {
				return err
			}
		}

		vf := ec.VariablesFile
		if vf != nil {
			tfvars, err := getVariablesFromTfvarsFile(
				log,
				githubRepoApi,
				vf.Source,
				env.Config.RepoBranch,
				vf.Path,
				)
			if err != nil {
				return err
			}
			vf.Variables = tfvars
		}

		application := argocd.GenerateEnvironmentComponentApps(*environment, *ec)

		tf  := terraformgenerator.TerraformGenerator{Log: log}
		err := tf.GenerateTerraform(fileUtil, ec, environment, envComponentDirectory)

		if err != nil {
			return err
		}

		if err = fileUtil.SaveYamlFile(*application, envComponentDirectory, ec.Name+".yaml"); err != nil {
			return err
		}
	}

	return nil
}

func getVariablesFromTfvarsFile(log logr.Logger, api github.RepositoryApi, repoUrl string, ref string, path string) ([]*stablev1alpha1.Variable, error) {
	log.Info("Downloading tfvars file", "repoUrl", repoUrl, "ref", ref, "path", path)
	buff, err := downloadTfvarsFile(log, api, repoUrl, ref, path)
	if err != nil {
		return nil, err
	}
	log.Info("Parsing variables from tfvars file")
	tfvars, err := parseTfvars(buff)
	if err != nil {
		return nil, err
	}
	return tfvars, nil
}

func downloadTfvarsFile(log logr.Logger, api github.RepositoryApi, repoUrl string, ref string, path string) ([]byte, error)  {
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

func parseTfvars(buff []byte) ([]*stablev1alpha1.Variable, error) {
	props, err := properties.Load(buff, properties.UTF8)
	if err != nil {
		return nil, err
	}

	var tfvars []*stablev1alpha1.Variable
	for name, value := range props.Map() {
		tfvars = append(tfvars, &stablev1alpha1.Variable{Name: name, Value: value})
	}
	return tfvars, nil
}

func generateAndSaveWorkflowOfWorkflows(fileUtil file.UtilFile, environment *stablev1alpha1.Environment, envComponentDirectory string) error {

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

func generateAndSaveEnvironmentApp(fileUtil file.UtilFile, environment *stablev1alpha1.Environment, envDirectory string) error {
	envApp := argocd.GenerateEnvironmentApp(*environment)
	envYAML := fmt.Sprintf("%s-environment.yaml", environment.Spec.EnvName)

	if err := fileUtil.SaveYamlFile(*envApp, envDirectory, envYAML); err != nil {
		return err
	}

	return nil
}
