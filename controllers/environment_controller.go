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

	r.Get(ctx, req.NamespacedName, environment)

	envDirectory := il.EnvironmentDirectory(environment.Spec.TeamName)
	envComponentDirectory := il.EnvironmentComponentDirectory(environment.Spec.TeamName, environment.Spec.EnvName)
	var fileUtil file.UtilFile = &file.UtilFileService{}

	if err := generateAndSaveEnvironmentApp(fileUtil, environment, envDirectory); err != nil {
		return ctrl.Result{}, err
	}

	if err := generateAndSaveEnvironmentComponents(fileUtil, environment, envComponentDirectory); err != nil {
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

// SetupWithManager sets up the Environment Controller with Manager
func (r *EnvironmentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&stablev1alpha1.Environment{}).
		Complete(r)
}

func generateAndSaveEnvironmentComponents(fileUtil file.UtilFile, environment *stablev1alpha1.Environment, envComponentDirectory string) error {
	for _, environmentComponent := range environment.Spec.EnvironmentComponent {
		if environmentComponent.Variables != nil {
			fileName := fmt.Sprintf("%s.tfvars", environmentComponent.Name)

			if err := fileUtil.SaveVarsToFile(environmentComponent.Variables, envComponentDirectory, fileName); err != nil {
				return err
			}

			environmentComponent.VariablesFile = &stablev1alpha1.VariablesFile{
				Source: env.Config.ILRepoURL,
				Path:   envComponentDirectory + "/" + fileName,
			}
		}

		application := argocd.GenerateEnvironmentComponentApps(*environment, *environmentComponent)

		var tf terraformgenerator.UtilTerraformGenerator = terraformgenerator.TerraformGenerator{}
		err := tf.GenerateTerraform(fileUtil, environmentComponent, environment, envComponentDirectory)

		if err != nil {
			return err
		}

		if err = fileUtil.SaveYamlFile(*application, envComponentDirectory, environmentComponent.Name+".yaml"); err != nil {
			return err
		}
	}

	return nil
}

func generateAndSaveWorkflowOfWorkflows(fileUtil file.UtilFile, environment *stablev1alpha1.Environment, envComponentDirectory string) error {

	// WIP, below command is for testing
	experimentalworkflow := argoWorkflow.GenerateWorkflowOfWorkflows(*environment)
	if err := fileUtil.SaveYamlFile(*experimentalworkflow, envComponentDirectory, "/experimental_wofw.yaml"); err != nil {
		return err
	}

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
