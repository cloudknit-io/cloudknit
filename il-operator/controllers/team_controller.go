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
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/file"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/repo"
	"github.com/go-logr/logr"
	"go.uber.org/atomic"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"

	stablev1 "github.com/compuzest/zlifecycle-il-operator/api/v1"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/env"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/github"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/il"

	"github.com/compuzest/zlifecycle-il-operator/controllers/argocd"
)

// TeamReconciler reconciles a Team object
type TeamReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=stable.compuzest.com,resources=teams,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=stable.compuzest.com,resources=teams/status,verbs=get;update;patch

var teamReconcileInitialRun = atomic.NewBool(true)

// Reconcile method called everytime there is a change in Team Custom Resource
func (r *TeamReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	delayTeamReconcileOnInitialRun(r.Log, 10)
	start := time.Now()

	team := &stablev1.Team{}
	if err := r.Get(ctx, req.NamespacedName, team); err != nil {
		if errors.IsNotFound(err) {
			r.Log.Info(
				"team missing from cache, ending reconcile",
				"name", req.Name,
				"namespace", req.Namespace,
			)
			return ctrl.Result{}, nil
		}
		r.Log.Error(
			err,
			"error occurred while getting Team",
			"name", req.Name,
			"namespace", req.Namespace,
		)

		return ctrl.Result{}, err
	}

	argocdApi   := argocd.NewHttpClient(r.Log, env.Config.ArgocdServerUrl)
	if _, err := argocd.TryCreateProject(r.Log, team.Spec.TeamName, env.Config.GitHubOrg); err != nil {
		return ctrl.Result{}, err
	}

	teamYAML := fmt.Sprintf("%s-team.yaml", team.Spec.TeamName)

	fileUtil := &file.UtilFileService{}
	if err := fileUtil.CreateEmptyDirectory(il.EnvironmentDirectory(team.Spec.TeamName)); err != nil {
		return ctrl.Result{}, err
	}

	teamRepo := team.Spec.ConfigRepo.Source
	operatorSshSecret := env.Config.ZlifecycleMasterRepoSshSecret
	operatorNamespace := env.Config.ZlifecycleOperatorNamespace

	if err := repo.TryRegisterRepo(r.Client, r.Log, ctx, argocdApi, teamRepo, operatorNamespace, operatorSshSecret); err != nil {
		return ctrl.Result{}, err
	}

	if err := generateAndSaveTeamApp(team, teamYAML); err != nil {
		return ctrl.Result{}, err
	}

	if err := generateAndSaveConfigWatchers(team, teamYAML); err != nil {
		return ctrl.Result{}, err
	}

	dirty, err := github.CommitAndPushFiles(
		env.Config.ILRepoSourceOwner,
		env.Config.ILRepoName,
		[]string{il.Config.TeamDirectory, il.Config.ConfigWatcherDirectory},
		env.Config.RepoBranch,
		fmt.Sprintf("Reconciling team %s", team.Spec.TeamName),
		env.Config.GithubSvcAccntName,
		env.Config.GithubSvcAccntEmail,
	)
	if err != nil {
		return ctrl.Result{}, err
	}
	if dirty {
		r.Log.Info(
			"Committed new changes to IL repo",
			"team", team.Spec.TeamName,
		)
	} else {
		r.Log.Info(
			"No git changes to commit, no-op reconciliation.",
			"team", team.Spec.TeamName,
		)
	}

	githubApi := github.NewHttpRepositoryClient(env.Config.GitHubAuthToken, ctx)
	_, err = github.CreateRepoWebhook(r.Log, githubApi, teamRepo, env.Config.ArgocdHookUrl, env.Config.GitHubWebhookSecret)
	if err != nil {
		return ctrl.Result{}, err
	}

	duration := time.Since(start)
	r.Log.Info(
		"Reconcile finished",
		"duration", duration,
		"team", team.Spec.TeamName,
	)

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the Company Controller with Manager
func (r *TeamReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&stablev1.Team{}).
		Complete(r)
}

func delayTeamReconcileOnInitialRun(log logr.Logger, seconds int64) {
	if teamReconcileInitialRun.Load() {
		log.Info(
			"Delaying Team reconcile on initial run to wait for Company operator",
			"duration", fmt.Sprintf("%ds", seconds),
		)
		time.Sleep(time.Duration(seconds) * time.Second)
		teamReconcileInitialRun.Store(false)
	}
}

func generateAndSaveTeamApp(team *stablev1.Team, teamYAML string) error {
	teamApp := argocd.GenerateTeamApp(*team)
	fileUtil := &file.UtilFileService{}

	if err := fileUtil.SaveYamlFile(*teamApp, il.Config.TeamDirectory, teamYAML); err != nil {
		return err
	}

	return nil
}

func generateAndSaveConfigWatchers(team *stablev1.Team, teamYAML string) error {
	teamConfigWatcherApp := argocd.GenerateTeamConfigWatcherApp(*team)
	fileUtil := &file.UtilFileService{}

	if err := fileUtil.SaveYamlFile(*teamConfigWatcherApp, il.Config.ConfigWatcherDirectory, teamYAML); err != nil {
		return err
	}

	return nil
}
