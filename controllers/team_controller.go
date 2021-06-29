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
	file "github.com/compuzest/zlifecycle-il-operator/controllers/util/file"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/repo"
	"github.com/go-logr/logr"
	"go.uber.org/atomic"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"

	stablev1alpha1 "github.com/compuzest/zlifecycle-il-operator/api/v1alpha1"
	env "github.com/compuzest/zlifecycle-il-operator/controllers/util/env"
	github "github.com/compuzest/zlifecycle-il-operator/controllers/util/github"
	il "github.com/compuzest/zlifecycle-il-operator/controllers/util/il"

	argocd "github.com/compuzest/zlifecycle-il-operator/controllers/argocd"
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
func (r *TeamReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	delayTeamReconcileOnInitialRun(r.Log, 10)
	start := time.Now()
	ctx := context.Background()

	team := &stablev1alpha1.Team{}
	fileUtil := &file.UtilFileService{}
	if err := r.Get(ctx, req.NamespacedName, team); err != nil {
		if errors.IsNotFound(err) {
			r.Log.Info(
				"Team missing from cache, ending reconcile...",
				"name", req.Name,
				"namespace", req.Namespace,
			)
		} else {
			r.Log.Error(
				err,
				"Error occurred while getting Team...",
				"name", req.Name,
				"namespace", req.Namespace,
			)
		}

		return ctrl.Result{}, nil
	}

	teamYAML := fmt.Sprintf("%s-team.yaml", team.Spec.TeamName)

	if err := fileUtil.CreateEmptyDirectory(il.EnvironmentDirectory(team.Spec.TeamName)); err != nil {
		return ctrl.Result{}, err
	}

	teamRepo := team.Spec.ConfigRepo.Source
	repoSecret := il.SSHKeyName()

	argocdApi := argocd.NewHttpClient(r.Log, env.Config.ArgocdServerUrl)
	if err := repo.TryRegisterRepo(r.Client, r.Log, ctx, argocdApi, teamRepo, req.Namespace, repoSecret); err != nil {
		return ctrl.Result{}, err
	}

	if err := generateAndSaveTeamApp(team, teamYAML); err != nil {
		return ctrl.Result{}, err
	}

	if err := generateAndSaveConfigWatchers(team, teamYAML); err != nil {
		return ctrl.Result{}, err
	}

	if err := github.CommitAndPushFiles(
		env.Config.ILRepoSourceOwner,
		env.Config.ILRepoName,
		[]string{il.Config.TeamDirectory, il.Config.ConfigWatcherDirectory},
		env.Config.RepoBranch,
		fmt.Sprintf("Reconciling team %s", team.Spec.TeamName),
		env.Config.GithubSvcAccntName,
		env.Config.GithubSvcAccntEmail); err != nil {
		return ctrl.Result{}, err
	}

	githubApi := github.NewHttpRepositoryClient(env.Config.GitHubAuthToken, ctx)
	_, err := github.CreateRepoWebhook(r.Log, githubApi, teamRepo, env.Config.ArgocdHookUrl, env.Config.GitHubWebhookSecret)
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
		For(&stablev1alpha1.Team{}).
		Complete(r)
}

func delayTeamReconcileOnInitialRun(log logr.Logger, seconds int64) {
	if teamReconcileInitialRun.Load() == true {
		log.Info(
			"Delaying Team reconcile on initial run to wait for Company operator",
			"duration", fmt.Sprintf("%ds", seconds * 1000),
		)
		time.Sleep(time.Duration(seconds) * time.Second)
		teamReconcileInitialRun.Store(false)
	}
}

func generateAndSaveTeamApp(team *stablev1alpha1.Team, teamYAML string) error {
	teamApp := argocd.GenerateTeamApp(*team)
	fileUtil := &file.UtilFileService{}

	if err := fileUtil.SaveYamlFile(*teamApp, il.Config.TeamDirectory, teamYAML); err != nil {
		return err
	}

	return nil
}

func generateAndSaveConfigWatchers(team *stablev1alpha1.Team, teamYAML string) error {
	teamConfigWatcherApp := argocd.GenerateTeamConfigWatcherApp(*team)
	fileUtil := &file.UtilFileService{}

	if err := fileUtil.SaveYamlFile(*teamConfigWatcherApp, il.Config.ConfigWatcherDirectory, teamYAML); err != nil {
		return err
	}

	return nil
}
