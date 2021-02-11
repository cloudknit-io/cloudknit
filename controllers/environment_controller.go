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

func (r *EnvironmentReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	//log := r.Log.WithValues("environment", req.NamespacedName)

	environment := &stablev1alpha1.Environment{}

	r.Get(ctx, req.NamespacedName, environment)

	envDirectory := il.Config.TeamDirectory + "/" + environment.Spec.TeamName + "-team-environment"
	envComponentDirectory := envDirectory + "/" + environment.Spec.EnvName + "-environment-component"

	environ := argocd.GenerateEnvironmentApp(*environment)
	if err := file.SaveYamlFile(*environ, envDirectory, environment.Spec.EnvName+"-environment.yaml"); err != nil {
		return ctrl.Result{}, err
	}

	for _, terraformConfig := range environment.Spec.TerraformConfigs {
		if terraformConfig.Variables != nil {
			fileName := terraformConfig.ConfigName + ".tfvars"
			if err := file.SaveVarsToFile(terraformConfig.Variables, envComponentDirectory, fileName); err != nil {
				return ctrl.Result{}, err
			}

			terraformConfig.VariablesFile = &stablev1alpha1.VariablesFile{
				Source: env.Config.ILRepoURL,
				Path:   envComponentDirectory + "/" + fileName,
			}
		}

		application := argocd.GenerateTerraformConfigApps(*environment, *terraformConfig)

		if err := file.SaveYamlFile(*application, envComponentDirectory, terraformConfig.ConfigName+".yaml"); err != nil {
			return ctrl.Result{}, err
		}
	}

	workflow := argoWorkflow.GenerateWorkflowOfWorkflows(*environment)
	if err := file.SaveYamlFile(*workflow, envComponentDirectory, "/wofw.yaml"); err != nil {
		return ctrl.Result{}, err
	}

	presyncJob := k8s.GeneratePreSyncJob(*environment)
	if err := file.SaveYamlFile(*presyncJob, envComponentDirectory, "/presync-job.yaml"); err != nil {
		return ctrl.Result{}, err
	}

	// Avoid race condition on initial Reconcile, collides with Team controller commit
	time.Sleep(10 * time.Second)

	github.CommitAndPushFiles(
		"CompuZest",
		env.Config.ILRepoName,
		[]string{envDirectory, envComponentDirectory},
		env.Config.RepoBranch,
		fmt.Sprintf("Reconciling environment %s", environment.Spec.EnvName),
		env.Config.GithubSvcAccntName,
		env.Config.GithubSvcAccntEmail)
	return ctrl.Result{}, nil
}

func (r *EnvironmentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&stablev1alpha1.Environment{}).
		Complete(r)
}
