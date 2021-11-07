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
	"strings"
	"sync"
	"time"

	"github.com/compuzest/zlifecycle-il-operator/controllers/util/git"

	"k8s.io/apiserver/pkg/registry/generic/registry"

	"sigs.k8s.io/controller-runtime/pkg/reconcile"

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

	stablev1 "github.com/compuzest/zlifecycle-il-operator/api/v1"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/env"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/github"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/il"

	"github.com/compuzest/zlifecycle-il-operator/controllers/argocd"
)

// TeamReconciler reconciles a Team object.
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

// Reconcile method called everytime there is a change in Team Custom Resource.
func (r *TeamReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// delay Team Reconcile so Company reconciles finish first
	delayTeamReconcileOnInitialRun(r.Log, 15)
	start := time.Now()

	// init logic
	var initError error
	initArgocdAdminRbacLock.Do(func() {
		initError = r.initArgocdAdminRbac(ctx)
	})
	if initError != nil {
		if strings.Contains(initError.Error(), registry.OptimisticLockErrorMsg) {
			// do manual retry without error
			return reconcile.Result{RequeueAfter: time.Second * 1}, nil
		}
		return ctrl.Result{}, initError
	}

	// fetch Team resource from k8s cache
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

	// services init
	fileAPI := &file.OsFileService{}
	argocdAPI := argocd.NewHTTPClient(ctx, r.Log, env.Config.ArgocdServerURL)
	repoAPI := github.NewHTTPRepositoryAPI(ctx)
	gitAPI, err := git.NewGoGit(ctx)
	if err != nil {
		return ctrl.Result{}, err
	}

	// temp clone IL repo
	tempILRepoDir, cleanup, err := git.CloneTemp(gitAPI, ilRepoURL)
	if err != nil {
		return ctrl.Result{}, err
	}
	defer cleanup()

	if err := r.updateArgocdRbac(ctx, team); err != nil {
		if strings.Contains(err.Error(), registry.OptimisticLockErrorMsg) {
			// do manual retry without error
			return reconcile.Result{RequeueAfter: time.Second * 1}, nil
		}
		return ctrl.Result{}, err
	}

	if _, err := argocd.TryCreateProject(ctx, r.Log, team.Spec.TeamName, env.Config.GitHubOrg); err != nil {
		return ctrl.Result{}, err
	}

	if err := fileAPI.CreateEmptyDirectory(il.EnvironmentDirectoryPath(team.Spec.TeamName)); err != nil {
		return ctrl.Result{}, err
	}

	teamRepo := team.Spec.ConfigRepo.Source
	operatorSSHSecret := env.Config.ZlifecycleMasterRepoSSHSecret
	operatorNamespace := env.Config.ZlifecycleOperatorNamespace

	if err := repo.TryRegisterRepo(ctx, r.Client, r.Log, argocdAPI, teamRepo, operatorNamespace, operatorSSHSecret); err != nil {
		return ctrl.Result{}, err
	}

	teamAppFilename := fmt.Sprintf("%s-team.yaml", team.Spec.TeamName)

	if err := generateAndSaveTeamApp(fileAPI, team, teamAppFilename, tempILRepoDir); err != nil {
		return ctrl.Result{}, err
	}

	if err := generateAndSaveConfigWatchers(fileAPI, team, teamAppFilename, tempILRepoDir); err != nil {
		return ctrl.Result{}, err
	}

	commitInfo := git.CommitInfo{
		Author: githubSvcAccntName,
		Email:  githubSvcAccntEmail,
		Msg:    fmt.Sprintf("Reconciling team %s", team.Spec.TeamName),
	}
	pushed, err := gitAPI.CommitAndPush(&commitInfo)
	if err != nil {
		return ctrl.Result{}, err
	}
	if pushed {
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

	_, err = github.CreateRepoWebhook(r.Log, repoAPI, teamRepo, env.Config.ArgocdHookURL, env.Config.GitHubWebhookSecret)
	if err != nil {
		r.Log.Error(err, "error creating Team webhook", "team", team.Spec.TeamName)
	}

	duration := time.Since(start)
	r.Log.Info(
		"Reconcile finished",
		"duration", duration,
		"team", team.Spec.TeamName,
	)

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the Company Controller with Manager.
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

	return r.Client.Update(ctx, &rbacCm)
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
	return r.Client.Update(ctx, &rbacCm)
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

func generateAndSaveTeamApp(fileAPI file.Service, team *stablev1.Team, filename string, ilRepoDir string) error {
	teamApp := argocd.GenerateTeamApp(team)

	return fileAPI.SaveYamlFile(*teamApp, il.TeamDirectoryAbsolutePath(ilRepoDir), filename)
}

func generateAndSaveConfigWatchers(fileAPI file.Service, team *stablev1.Team, filename string, ilRepoDir string) error {
	teamConfigWatcherApp := argocd.GenerateTeamConfigWatcherApp(team)

	return fileAPI.SaveYamlFile(*teamConfigWatcherApp, il.ConfigWatcherDirectoryAbsolutePath(ilRepoDir), filename)
}
