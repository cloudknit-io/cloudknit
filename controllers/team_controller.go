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
	"github.com/compuzest/zlifecycle-il-operator/controllers/apm"
	"github.com/sirupsen/logrus"
	"strings"
	"sync"
	"time"

	"github.com/compuzest/zlifecycle-il-operator/controllers/zerrors"
	perrors "github.com/pkg/errors"

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
	LogV2  *logrus.Entry
	Scheme *runtime.Scheme
	APM    apm.APM
}

// +kubebuilder:rbac:groups=stable.compuzest.com,resources=teams,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=stable.compuzest.com,resources=teams/status,verbs=get;update;patch

var (
	initArgocdAdminRbacLock     sync.Once
	teamReconcileInitialRunLock = atomic.NewBool(true)
)

// Reconcile method called everytime there is a change in Team Custom Resource.
func (r *TeamReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	if !checkIsNamespaceWatched(req.NamespacedName.Namespace) {
		r.LogV2.WithFields(logrus.Fields{
			"object":           fmt.Sprintf("%s/%s", req.NamespacedName.Namespace, req.NamespacedName.Name),
			"namespace":        req.NamespacedName.Namespace,
			"watchedNamespace": env.Config.KubernetesOperatorWatchedNamespace,
		}).Info("Namespace is not configured to be watched by operator")
		return ctrl.Result{}, nil
	}
	if resource := "team"; !checkIsResourceWatched(resource) {
		r.LogV2.WithFields(logrus.Fields{
			"object":           fmt.Sprintf("%s/%s", req.NamespacedName.Namespace, req.NamespacedName.Name),
			"resource":         resource,
			"watchedResources": env.Config.KubernetesOperatorWatchedResources,
		}).Info("Resource is not configured to be watched by operator")
		return ctrl.Result{}, nil
	}

	// delay Team Reconcile so Company reconciles finish first
	delayTeamReconcileOnInitialRun(r.LogV2, 15)
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
			r.LogV2.WithFields(logrus.Fields{
				"name":      req.Name,
				"namespace": req.Namespace,
			}).Info("team missing from cache, ending reconcile")
			return ctrl.Result{}, nil
		}
		teamErr := zerrors.NewTeamError(team.Spec.TeamName, perrors.Wrap(err, "error getting team from k8s cache"))
		return ctrl.Result{}, teamErr
	}

	// start apm transaction
	txName := fmt.Sprintf("teamreconciler.%s", team.Spec.TeamName)
	tx := r.APM.StartTransaction(txName)
	apmCtx := ctx
	if tx != nil {
		tx.AddAttribute("company", env.Config.CompanyName)
		tx.AddAttribute("team", team.Spec.TeamName)
		apmCtx = r.APM.NewContext(ctx, tx)
		r.LogV2 = r.LogV2.WithField("team", team.Spec.TeamName).WithContext(apmCtx)
		r.LogV2.WithField("name", txName).Info("Creating APM transaction")
		defer tx.End()
	}

	// services init
	fileAPI := &file.OsFileService{}
	argocdAPI := argocd.NewHTTPClient(apmCtx, r.Log, env.Config.ArgocdServerURL)
	repoAPI := github.NewHTTPRepositoryAPI(apmCtx)
	gitAPI, err := git.NewGoGit(apmCtx)
	if err != nil {
		teamErr := zerrors.NewTeamError(team.Spec.TeamName, perrors.Wrap(err, "error instantiating git API"))
		return ctrl.Result{}, r.APM.NoticeError(tx, r.LogV2, teamErr)
	}

	// temp clone IL repo
	tempILRepoDir, cleanup, err := git.CloneTemp(gitAPI, zlILRepoURL)
	if err != nil {
		teamErr := zerrors.NewTeamError(team.Spec.TeamName, perrors.Wrap(err, "error running git temp clone"))
		return ctrl.Result{}, r.APM.NoticeError(tx, r.LogV2, teamErr)
	}
	defer cleanup()

	if err := r.updateArgocdRbac(apmCtx, team); err != nil {
		if strings.Contains(err.Error(), registry.OptimisticLockErrorMsg) {
			// do manual retry without error
			return reconcile.Result{RequeueAfter: time.Second * 1}, nil
		}
		teamErr := zerrors.NewTeamError(team.Spec.TeamName, perrors.Wrap(err, "error updating argocd rbac"))
		return ctrl.Result{}, r.APM.NoticeError(tx, r.LogV2, teamErr)
	}

	if _, err := argocd.TryCreateProject(apmCtx, r.Log, team.Spec.TeamName, env.Config.GitHubCompanyOrganization); err != nil {
		teamErr := zerrors.NewTeamError(team.Spec.TeamName, perrors.Wrap(err, "error trying to create argocd project"))
		return ctrl.Result{}, r.APM.NoticeError(tx, r.LogV2, teamErr)
	}

	if err := fileAPI.CreateEmptyDirectory(il.EnvironmentDirectoryPath(team.Spec.TeamName)); err != nil {
		teamErr := zerrors.NewTeamError(team.Spec.TeamName, perrors.Wrap(err, "error creating team dir"))
		return ctrl.Result{}, r.APM.NoticeError(tx, r.LogV2, teamErr)
	}

	teamRepo := team.Spec.ConfigRepo.Source
	operatorSSHSecret := env.Config.GitSSHSecretName
	operatorNamespace := env.Config.KubernetesServiceNamespace

	if err := repo.TryRegisterRepo(apmCtx, r.Client, r.Log, argocdAPI, teamRepo, operatorNamespace, operatorSSHSecret); err != nil {
		teamErr := zerrors.NewTeamError(team.Spec.TeamName, perrors.Wrap(err, "error registering argocd team repo"))
		return ctrl.Result{}, r.APM.NoticeError(tx, r.LogV2, teamErr)
	}

	teamAppFilename := fmt.Sprintf("%s-team.yaml", team.Spec.TeamName)

	if err := generateAndSaveTeamApp(fileAPI, team, teamAppFilename, tempILRepoDir); err != nil {
		teamErr := zerrors.NewTeamError(team.Spec.TeamName, perrors.Wrap(err, "error generating team argocd app"))
		return ctrl.Result{}, r.APM.NoticeError(tx, r.LogV2, teamErr)
	}

	if err := generateAndSaveConfigWatchers(fileAPI, team, teamAppFilename, tempILRepoDir); err != nil {
		teamErr := zerrors.NewTeamError(team.Spec.TeamName, perrors.Wrap(err, "error generating team config watchers"))
		return ctrl.Result{}, r.APM.NoticeError(tx, r.LogV2, teamErr)
	}

	commitInfo := git.CommitInfo{
		Author: githubSvcAccntName,
		Email:  githubSvcAccntEmail,
		Msg:    fmt.Sprintf("Reconciling team %s", team.Spec.TeamName),
	}
	pushed, err := gitAPI.CommitAndPush(&commitInfo)
	if err != nil {
		teamErr := zerrors.NewTeamError(team.Spec.TeamName, perrors.Wrap(err, "error running commit and push"))
		return ctrl.Result{}, r.APM.NoticeError(tx, r.LogV2, teamErr)
	}
	if pushed {
		r.LogV2.Info("Committed new changes to IL repo")
	} else {
		r.LogV2.Info("No git changes to commit, no-op reconciliation.")
	}

	_, err = github.CreateRepoWebhook(r.Log, repoAPI, teamRepo, env.Config.ArgocdWebhookURL, env.Config.GitHubWebhookSecret)
	if err != nil {
		r.LogV2.WithError(err).Error("error creating Team webhook")
	}

	duration := time.Since(start)
	r.LogV2.WithField("duration", duration).Info("Reconcile finished")
	attrs := map[string]interface{}{
		"duration": duration,
		"team":     team.Spec.TeamName,
	}
	r.APM.RecordCustomEvent("teamreconciler", attrs)

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
	if err := r.Client.Get(ctx, types.NamespacedName{Name: "argocd-rbac-cm", Namespace: env.ArgoNamespace()}, &rbacCm); err != nil {
		return err
	}
	if rbacCm.Data == nil {
		rbacCm.Data = make(map[string]string)
	}
	admin := "admin"
	oldPolicyCsv := rbacCm.Data["policy.csv"]
	oidcGroup := fmt.Sprintf("%s:%s", env.Config.GitHubCompanyOrganization, admin)
	newPolicyCsv, err := argocd.GenerateAdminRbacConfig(r.Log, oldPolicyCsv, oidcGroup, admin)
	if err != nil {
		return err
	}
	rbacCm.Data["policy.csv"] = newPolicyCsv

	return r.Client.Update(ctx, &rbacCm)
}

func (r *TeamReconciler) updateArgocdRbac(ctx context.Context, t *stablev1.Team) error {
	rbacCm := v1.ConfigMap{}
	if err := r.Client.Get(ctx, types.NamespacedName{Name: "argocd-rbac-cm", Namespace: env.ArgoNamespace()}, &rbacCm); err != nil {
		return err
	}
	if rbacCm.Data == nil {
		rbacCm.Data = make(map[string]string)
	}
	teamName := t.Spec.TeamName
	oldPolicyCsv := rbacCm.Data["policy.csv"]
	oidcGroup := fmt.Sprintf("%s:%s", env.Config.GitHubCompanyOrganization, teamName)
	newPolicyCsv, err := argocd.GenerateNewRbacConfig(r.Log, oldPolicyCsv, oidcGroup, teamName, t.Spec.Permissions)
	if err != nil {
		return err
	}
	rbacCm.Data["policy.csv"] = newPolicyCsv
	return r.Client.Update(ctx, &rbacCm)
}

func delayTeamReconcileOnInitialRun(log *logrus.Entry, seconds int64) {
	if teamReconcileInitialRunLock.Load() {
		log.WithField("duration", fmt.Sprintf("%ds", seconds)).Info("Delaying Team reconcile on initial run to wait for Company operator")
		time.Sleep(time.Duration(seconds) * time.Second)
		teamReconcileInitialRunLock.Store(false)
	}
}

func generateAndSaveTeamApp(fileAPI file.FSAPI, team *stablev1.Team, filename string, ilRepoDir string) error {
	teamApp := argocd.GenerateTeamApp(team)

	return fileAPI.SaveYamlFile(*teamApp, il.TeamDirectoryAbsolutePath(ilRepoDir), filename)
}

func generateAndSaveConfigWatchers(fileAPI file.FSAPI, team *stablev1.Team, filename string, ilRepoDir string) error {
	teamConfigWatcherApp := argocd.GenerateTeamConfigWatcherApp(team)

	return fileAPI.SaveYamlFile(*teamConfigWatcherApp, il.ConfigWatcherDirectoryAbsolutePath(ilRepoDir), filename)
}
