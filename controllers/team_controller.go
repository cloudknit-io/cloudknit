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
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sync"
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

var (
	initArgocdAdminRbacLock     sync.Once
	teamReconcileInitialRunLock = atomic.NewBool(true)
)

// Reconcile method called everytime there is a change in Team Custom Resource
func (r *TeamReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	delayTeamReconcileOnInitialRun(r.Log, 15)
	start := time.Now()

	var initError error
	initArgocdAdminRbacLock.Do(func() {
		initError = r.initArgocdAdminRbac(ctx)
	})
	if initError != nil {
		return ctrl.Result{}, initError
	}

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

	if err := r.updateArgocdRbac(ctx, team); err != nil {
		return ctrl.Result{}, err
	}

	argocdApi := argocd.NewHttpClient(r.Log, env.Config.ArgocdServerUrl)
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

func (r *TeamReconciler) initArgocdAdminRbac(ctx context.Context) error {
	rbacCm := v1.ConfigMap{}
	if err := r.Client.Get(ctx, types.NamespacedName{Name: "argocd-rbac-cm", Namespace: "argocd"}, &rbacCm); err != nil {
		return err
	}
	if rbacCm.Data == nil {
		rbacCm.Data = make(map[string]string)
	}
	admin := "admin"
	oldPolicyCsv := rbacCm.Data["policy.csv"]
	oidcGroup := fmt.Sprintf("%s:%s", env.Config.GitHubOrg, admin)
	newPolicyCsv, err := argocd.GenerateAdminRbacConfig(r.Log, oldPolicyCsv, oidcGroup, admin)
	if err != nil {
		return err
	}
	rbacCm.Data["policy.csv"] = newPolicyCsv
	if err := r.Client.Update(ctx, &rbacCm); err != nil {
		return err
	}
	return nil
}

func (r *TeamReconciler) updateArgocdRbac(ctx context.Context, t *stablev1.Team) error {
	rbacCm := v1.ConfigMap{}
	if err := r.Client.Get(ctx, types.NamespacedName{Name: "argocd-rbac-cm", Namespace: "argocd"}, &rbacCm); err != nil {
		return err
	}
	if rbacCm.Data == nil {
		rbacCm.Data = make(map[string]string)
	}
	teamName := t.Spec.TeamName
	oldPolicyCsv := rbacCm.Data["policy.csv"]
	oidcGroup := fmt.Sprintf("%s:%s", env.Config.GitHubOrg, teamName)
	newPolicyCsv, err := argocd.GenerateNewRbacConfig(r.Log, oldPolicyCsv, oidcGroup, teamName, t.Spec.Permissions)
	if err != nil {
		return err
	}
	rbacCm.Data["policy.csv"] = newPolicyCsv
	if err := r.Client.Update(ctx, &rbacCm); err != nil {
		return err
	}
	return nil
}

func delayTeamReconcileOnInitialRun(log logr.Logger, seconds int64) {
	if teamReconcileInitialRunLock.Load() {
		log.Info(
			"Delaying Team reconcile on initial run to wait for Company operator",
			"duration", fmt.Sprintf("%ds", seconds),
		)
		time.Sleep(time.Duration(seconds) * time.Second)
		teamReconcileInitialRunLock.Store(false)
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
