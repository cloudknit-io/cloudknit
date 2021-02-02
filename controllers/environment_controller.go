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

	github "github.com/compuzest/zlifecycle-il-operator/controllers/util/github"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	kClient "sigs.k8s.io/controller-runtime/pkg/client"

	stablev1alpha1 "github.com/compuzest/zlifecycle-il-operator/api/v1alpha1"
	argocd "github.com/compuzest/zlifecycle-il-operator/controllers/argocd"
	argoWorkflow "github.com/compuzest/zlifecycle-il-operator/controllers/argoworkflow"
	k8s "github.com/compuzest/zlifecycle-il-operator/controllers/kubernetes"
	config "github.com/compuzest/zlifecycle-il-operator/controllers/util/config"
	file "github.com/compuzest/zlifecycle-il-operator/controllers/util/file"
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

	teamEnvPrefix := environment.Spec.TeamName + "/" + environment.Spec.EnvName

	env := argocd.GenerateEnvironmentApp(*environment)
	file.SaveYamlFile(*env, teamEnvPrefix+".yaml")

	for _, terraformConfig := range environment.Spec.TerraformConfigs {
		if terraformConfig.Variables != nil {
			filePath := teamEnvPrefix + "/" + terraformConfig.ConfigName + ".tfvars"
			file.SaveVarsToFile(terraformConfig.Variables, filePath)
			terraformConfig.VariablesFile = &stablev1alpha1.VariablesFile{
				Source: config.ILRepoURL,
				Path:   filePath,
			}
		}

		application := argocd.GenerateTerraformConfigApps(*environment, *terraformConfig)

		file.SaveYamlFile(*application, teamEnvPrefix+"/"+terraformConfig.ConfigName+".yaml")
	}

	workflow := argoWorkflow.GenerateWorkflowOfWorkflows(*environment)
	file.SaveYamlFile(*workflow, teamEnvPrefix+"/wofw.yaml")

	presyncJob := k8s.GeneratePreSyncJob(*environment)
	file.SaveYamlFile(*presyncJob, teamEnvPrefix+"/presync-job.yaml")

	github.CommitAndPushFiles(
		config.CompanyName,
		config.ILRepoName,
		environment.Spec.TeamName+"/",
		config.RepoBranch,
		config.GithubSvcAccntName,
		config.GithubSvcAccntEmail)

	return ctrl.Result{}, nil
}

func (r *EnvironmentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&stablev1alpha1.Environment{}).
		Complete(r)
}
