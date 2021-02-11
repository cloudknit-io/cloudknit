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

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	stablev1alpha1 "github.com/compuzest/zlifecycle-il-operator/api/v1alpha1"
	env "github.com/compuzest/zlifecycle-il-operator/controllers/util/env"
	github "github.com/compuzest/zlifecycle-il-operator/controllers/util/github"
	il "github.com/compuzest/zlifecycle-il-operator/controllers/util/il"

	argocd "github.com/compuzest/zlifecycle-il-operator/controllers/argocd"
	fileutil "github.com/compuzest/zlifecycle-il-operator/controllers/util/file"
)

// TeamReconciler reconciles a Team object
type TeamReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=stable.compuzest.com,resources=teams,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=stable.compuzest.com,resources=teams/status,verbs=get;update;patch

func (r *TeamReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	//log := r.Log.WithValues("team", req.NamespacedName)

	team := &stablev1alpha1.Team{}

	r.Get(ctx, req.NamespacedName, team)

	teamApp := argocd.GenerateTeamApp(*team)
	companyName := env.Config.CompanyName
	ilRepoName := env.Config.ILRepoName
	teamDirectory := il.Config.TeamDirectory
	configWatcherDirectory := il.Config.ConfigWatcherDirectory
	fileutil.SaveYamlFile(*teamApp, teamDirectory, team.Spec.TeamName+"-team.yaml")

	// generate config watchers
	teamConfigWatcherApp := argocd.GenerateTeamConfigWatcherApp(companyName, team.Spec.ConfigRepo.Source)
	envConfigWatcherApp := argocd.GenerateEnvironmentConfigWatcherApp(*team)
	fileutil.SaveYamlFile(*envConfigWatcherApp, configWatcherDirectory, team.Spec.TeamName+"-team.yaml")
	fileutil.SaveYamlFile(*teamConfigWatcherApp, configWatcherDirectory, companyName+".yaml")

	github.CommitAndPushFiles(
		"CompuZest",
		ilRepoName,
		[]string{teamDirectory, configWatcherDirectory},
		env.Config.RepoBranch,
		fmt.Sprintf("Reconciling team %s", team.Spec.TeamName),
		env.Config.GithubSvcAccntName,
		env.Config.GithubSvcAccntEmail)

	return ctrl.Result{}, nil
}

func (r *TeamReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&stablev1alpha1.Team{}).
		Complete(r)
}
