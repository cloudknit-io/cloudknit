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
	"github.com/compuzest/zlifecycle-il-operator/controllers/util"
	"strings"
	"sync"
	"time"

	"github.com/compuzest/zlifecycle-il-operator/controllers/watcherservices"

	"github.com/compuzest/zlifecycle-il-operator/controllers/codegen/file"
	"github.com/compuzest/zlifecycle-il-operator/controllers/codegen/il"

	"github.com/compuzest/zlifecycle-il-operator/controllers/apm"
	"github.com/compuzest/zlifecycle-il-operator/controllers/env"
	"github.com/compuzest/zlifecycle-il-operator/controllers/external/argocd"
	"github.com/compuzest/zlifecycle-il-operator/controllers/external/git"
	"github.com/compuzest/zlifecycle-il-operator/controllers/external/github"
	"github.com/sirupsen/logrus"

	"github.com/compuzest/zlifecycle-il-operator/controllers/zerrors"
	perrors "github.com/pkg/errors"

	"k8s.io/apiserver/pkg/registry/generic/registry"

	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/go-logr/logr"
	"go.uber.org/atomic"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	stablev1 "github.com/compuzest/zlifecycle-il-operator/api/v1"
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
	if shouldEndReconcile("team", r.LogV2) {
		return ctrl.Result{}, nil
	}

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
		r.LogV2.WithField("name", txName).Infof("Creating APM transaction for team %s", team.Spec.TeamName)
		defer tx.End()
	}

	teamRepoURL := team.Spec.ConfigRepo.Source

	// services init
	fileAPI := &file.OsFileService{}
	watcherServices, err := watcherservices.NewGitHubServices(apmCtx, r.Client, env.Config.GitHubCompanyOrganization, r.LogV2)
	if err != nil {
		teamErr := zerrors.NewTeamError(team.Spec.TeamName, perrors.Wrap(err, "error instantiating watcher services"))
		return ctrl.Result{}, r.APM.NoticeError(tx, r.LogV2, teamErr)
	}
	token, err := github.GenerateInstallationToken(r.LogV2, watcherServices.ILGitClient, env.Config.GitILRepositoryOwner)
	if err != nil {
		teamErr := zerrors.NewTeamError(
			team.Spec.TeamName, perrors.Wrapf(err, "error generating installation token for git organization [%s]", env.Config.GitILRepositoryOwner),
		)
		return ctrl.Result{}, r.APM.NoticeError(tx, r.LogV2, teamErr)
	}
	gitClient, err := git.NewGoGit(apmCtx, &git.GoGitOptions{Mode: git.ModeToken, Token: token})
	if err != nil {
		teamErr := zerrors.NewTeamError(team.Spec.TeamName, perrors.Wrap(err, "error instantiating git client"))
		return ctrl.Result{}, r.APM.NoticeError(tx, r.LogV2, teamErr)
	}

	// temp clone IL repo
	tempILRepoDir, cleanup, err := git.CloneTemp(gitClient, env.Config.ILZLifecycleRepositoryURL, r.LogV2)
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

	if _, err := argocd.TryCreateProject(apmCtx, watcherServices.ArgocdClient, r.Log, team.Spec.TeamName, env.Config.GitHubCompanyOrganization); err != nil {
		teamErr := zerrors.NewTeamError(team.Spec.TeamName, perrors.Wrap(err, "error trying to create argocd project"))
		return ctrl.Result{}, r.APM.NoticeError(tx, r.LogV2, teamErr)
	}

	if err := fileAPI.CreateEmptyDirectory(il.EnvironmentDirectoryPath(team.Spec.TeamName)); err != nil {
		teamErr := zerrors.NewTeamError(team.Spec.TeamName, perrors.Wrap(err, "error creating team dir"))
		return ctrl.Result{}, r.APM.NoticeError(tx, r.LogV2, teamErr)
	}

	if err := watcherServices.CompanyWatcher.Watch(teamRepoURL); err != nil {
		teamErr := zerrors.NewTeamError(team.Spec.TeamName, perrors.Wrap(err, "error registering argocd team repo via github app auth"))
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
		Author: env.Config.GitServiceAccountName,
		Email:  env.Config.GitServiceAccountEmail,
		Msg:    fmt.Sprintf("Reconciling team %s", team.Spec.TeamName),
	}
	pushed, err := gitClient.CommitAndPush(&commitInfo)
	if err != nil {
		teamErr := zerrors.NewTeamError(team.Spec.TeamName, perrors.Wrap(err, "error running commit and push"))
		return ctrl.Result{}, r.APM.NoticeError(tx, r.LogV2, teamErr)
	}
	if pushed {
		r.LogV2.Infof("Committed new changes for team %s to IL repo", team.Spec.TeamName)
	} else {
		r.LogV2.Infof("No git changes to commit for team %s, no-op reconciliation.", team.Spec.TeamName)
	}

	if !util.IsGitLabURL(teamRepoURL) {
		_, err = github.CreateRepoWebhook(r.LogV2, watcherServices.CompanyGitClient, teamRepoURL, env.Config.ArgocdWebhookURL, env.Config.GitHubWebhookSecret)
		if err != nil {
			r.LogV2.WithError(err).Error("error creating Team webhook")
		}
	}

	duration := time.Since(start)
	r.LogV2.WithField("duration", duration).Infof("Reconcile finished for team %s", team.Spec.TeamName)
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
	nn := types.NamespacedName{Name: env.Config.ArgocdRBACConfigMap, Namespace: env.ArgocdNamespace()}
	if err := r.Client.Get(ctx, nn, &rbacCm); err != nil {
		return perrors.Wrap(err, "error getting argocd rbacm configmap from k8s")
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
	if err := r.Client.Get(ctx, types.NamespacedName{Name: env.Config.ArgocdRBACConfigMap, Namespace: env.ArgocdNamespace()}, &rbacCm); err != nil {
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
