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

package controller

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/compuzest/zlifecycle-il-operator/controller/lib/apm"
	"github.com/compuzest/zlifecycle-il-operator/controller/lib/watcherservices"
	"github.com/compuzest/zlifecycle-il-operator/controller/lib/zerrors"

	"github.com/compuzest/zlifecycle-il-operator/controller/util"

	"github.com/compuzest/zlifecycle-il-operator/controller/codegen/file"
	"github.com/compuzest/zlifecycle-il-operator/controller/env"
	"github.com/compuzest/zlifecycle-il-operator/controller/external/argocd"
	"github.com/compuzest/zlifecycle-il-operator/controller/external/git"
	"github.com/compuzest/zlifecycle-il-operator/controller/external/github"
	perrors "github.com/pkg/errors"

	"github.com/sirupsen/logrus"

	"k8s.io/apimachinery/pkg/api/errors"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	stablev1 "github.com/compuzest/zlifecycle-il-operator/api/v1"
)

var initOperatorLock sync.Once

// CompanyReconciler reconciles a Company object.
type CompanyReconciler struct {
	client.Client
	Log    logr.Logger
	LogV2  *logrus.Entry
	Scheme *runtime.Scheme
	APM    apm.APM
}

// +kubebuilder:rbac:groups=stable.compuzest.com,resources=companies,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=stable.compuzest.com,resources=companies/status,verbs=get;update;patch

func (r *CompanyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	if shouldEndReconcile("company", r.LogV2) {
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
	if resource := "company"; !checkIsResourceWatched(resource) {
		r.LogV2.WithFields(logrus.Fields{
			"object":           fmt.Sprintf("%s/%s", req.NamespacedName.Namespace, req.NamespacedName.Name),
			"resource":         resource,
			"watchedResources": env.Config.KubernetesOperatorWatchedResources,
		}).Info("Resource is not configured to be watched by operator")
		return ctrl.Result{}, nil
	}

	start := time.Now()

	// get company resource from k8s cache
	company := &stablev1.Company{}
	if err := r.Get(ctx, req.NamespacedName, company); err != nil {
		if errors.IsNotFound(err) {
			r.LogV2.WithFields(logrus.Fields{
				"name":      req.Name,
				"namespace": req.Namespace,
			}).Info("Company missing from cache, ending reconcile")
			return ctrl.Result{}, nil
		}
		r.LogV2.WithError(err).WithFields(logrus.Fields{
			"name":      req.Name,
			"namespace": req.Namespace,
		}).Error("Error occurred while getting Company")

		return ctrl.Result{}, err
	}

	// start apm transaction
	txName := fmt.Sprintf("companyreconciler.%s", company.Spec.CompanyName)
	tx := r.APM.StartTransaction(txName)
	apmCtx := ctx
	if tx != nil {
		tx.AddAttribute("company", company.Spec.CompanyName)
		apmCtx = r.APM.NewContext(ctx, tx)
		r.LogV2.WithField("name", txName).Infof("Creating APM transaction for company %s", company.Spec.CompanyName)
		defer tx.End()
	}

	// vars
	companyRepoURL := company.Spec.ConfigRepo.Source
	ilZLRepoURL := env.Config.ILZLifecycleRepositoryURL

	// services init
	fileAPI := file.NewOSFileService()
	watcherServices, err := watcherservices.NewGitHubServices(apmCtx, r.Client, env.Config.GitHubCompanyOrganization, r.LogV2)
	if err != nil {
		companyErr := zerrors.NewCompanyError(
			company.Spec.CompanyName,
			perrors.Wrap(err, "error instantiating watcher services"),
		)
		return ctrl.Result{}, r.APM.NoticeError(tx, r.LogV2, companyErr)
	}
	token, err := github.GenerateInstallationToken(r.LogV2, watcherServices.ILGitClient, env.Config.GitILRepositoryOwner)
	if err != nil {
		companyErr := zerrors.NewCompanyError(
			company.Spec.CompanyName,
			perrors.Wrapf(err, "error generating installation token for git organization [%s]", env.Config.GitILRepositoryOwner),
		)
		return ctrl.Result{}, r.APM.NoticeError(tx, r.LogV2, companyErr)
	}
	gitClient, err := git.NewGoGit(apmCtx, &git.GoGitOptions{Mode: git.ModeToken, Token: token})
	if err != nil {
		companyErr := zerrors.NewCompanyError(
			company.Spec.CompanyName,
			perrors.Wrap(err, "error instantiating git client"),
		)
		return ctrl.Result{}, r.APM.NoticeError(tx, r.LogV2, companyErr)
	}

	// init logic
	var initOperatorError error
	initOperatorLock.Do(func() {
		initOperatorError = r.initCompany(ctx, watcherServices)
	})
	if initOperatorError != nil {
		companyErr := zerrors.NewCompanyError(company.Spec.CompanyName, perrors.Wrap(err, "error initializing company"))
		return ctrl.Result{}, r.APM.NoticeError(tx, r.LogV2, companyErr)
	}

	// temp clone IL repo
	tempILRepoDir, cleanup, err := git.CloneTemp(gitClient, ilZLRepoURL, r.LogV2)
	if err != nil {
		companyErr := zerrors.NewCompanyError(company.Spec.CompanyName, perrors.Wrapf(err, "error cloning temp dir for repo [%s]", env.Config.GitILRepositoryOwner))
		return ctrl.Result{}, r.APM.NoticeError(tx, r.LogV2, companyErr)
	}
	defer cleanup()

	// reconcile logic

	if err := watcherServices.CompanyWatcher.Watch(companyRepoURL); err != nil {
		companyErr := zerrors.NewCompanyError(company.Spec.CompanyName, perrors.Wrap(err, "error registering company config repo in argocd using github app auth"))
		return ctrl.Result{}, r.APM.NoticeError(tx, r.LogV2, companyErr)
	}

	if err := generateAndSaveCompanyApp(fileAPI, company, tempILRepoDir); err != nil {
		companyErr := zerrors.NewCompanyError(company.Spec.CompanyName, perrors.Wrap(err, "error generating and saving company argocd app"))
		return ctrl.Result{}, r.APM.NoticeError(tx, r.LogV2, companyErr)
	}

	if err := generateAndSaveCompanyConfigWatcher(fileAPI, company, tempILRepoDir); err != nil {
		companyErr := zerrors.NewCompanyError(company.Spec.CompanyName, perrors.Wrap(err, "error generating and saving company config watcher argocd app"))
		return ctrl.Result{}, r.APM.NoticeError(tx, r.LogV2, companyErr)
	}

	if err := watcherServices.ILWatcher.Watch(ilZLRepoURL); err != nil {
		companyErr := zerrors.NewCompanyError(company.Spec.CompanyName, perrors.Wrap(err, "error registering company IL repo in argocd"))
		return ctrl.Result{}, r.APM.NoticeError(tx, r.LogV2, companyErr)
	}

	commitInfo := git.CommitInfo{
		Author: env.Config.GitServiceAccountName,
		Email:  env.Config.GitServiceAccountEmail,
		Msg:    fmt.Sprintf("Reconciling company %s", company.Spec.CompanyName),
	}
	pushed, err := gitClient.CommitAndPush(&commitInfo)
	if err != nil {
		companyErr := zerrors.NewCompanyError(company.Spec.CompanyName, perrors.Wrap(err, "error running commit and push company IL changes"))
		return ctrl.Result{}, r.APM.NoticeError(tx, r.LogV2, companyErr)
	}
	if pushed {
		r.LogV2.Infof("Committed new changes for company %s to IL repo", company.Spec.CompanyName)
	} else {
		r.LogV2.Infof("No git changes to commit for company %s, no-op reconciliation.", company.Spec.CompanyName)
	}

	if err := argocd.TryCreateBootstrapApps(apmCtx, watcherServices.ArgocdClient, r.Log); err != nil {
		companyErr := zerrors.NewCompanyError(company.Spec.CompanyName, perrors.Wrap(err, "error creating company bootstrap argocd app"))
		return ctrl.Result{}, r.APM.NoticeError(tx, r.LogV2, companyErr)
	}

	if !util.IsGitLabURL(companyRepoURL) {
		// try create a webhook, it will fail if git service account does not have permissions to create it
		_, err = github.CreateRepoWebhook(r.LogV2, watcherServices.CompanyGitClient, companyRepoURL, env.Config.ArgocdWebhookURL, env.Config.GitHubWebhookSecret)
		if err != nil {
			r.LogV2.WithError(err).Error("error creating Company webhook")
		}
	}

	duration := time.Since(start)
	r.LogV2.WithField("duration", duration).Infof("Reconcile finished for company %s", company.Spec.CompanyName)

	return ctrl.Result{}, nil
}

func (r *CompanyReconciler) initCompany(ctx context.Context, services *watcherservices.WatcherServices) error {
	r.LogV2.Info("Running company operator init")

	r.LogV2.Info("Creating webhook for IL repo")
	if _, err := github.CreateRepoWebhook(
		r.LogV2,
		services.ILGitClient,
		env.Config.ILZLifecycleRepositoryURL,
		env.Config.ArgocdWebhookURL,
		env.Config.GitHubWebhookSecret,
	); err != nil {
		r.LogV2.WithError(err).WithField("repo", env.Config.ILZLifecycleRepositoryURL).Error("error creating Company IL ZL webhook")
	}
	if _, err := github.CreateRepoWebhook(
		r.LogV2,
		services.ILGitClient,
		env.Config.ILTerraformRepositoryURL,
		env.Config.ArgocdWebhookURL,
		env.Config.GitHubWebhookSecret,
	); err != nil {
		r.LogV2.WithError(err).WithField("repo", env.Config.ILTerraformRepositoryURL).Error("error creating Company IL TF webhook")
	}

	r.LogV2.Info("Updating default argocd cluster namespaces")
	if err := argocd.UpdateDefaultClusterNamespaces(
		r.LogV2,
		services.ArgocdClient,
		[]string{env.ArgocdNamespace(), env.ConfigNamespace(), env.ExecutorNamespace()},
	); err != nil {
		r.LogV2.Fatalf("error updating argocd cluster namespaces: %v", err)
	}

	r.LogV2.Info("Registering helm chart repo")
	return services.InternalWatcher.Watch(env.Config.GitHelmChartsRepository)
}

func (r *CompanyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&stablev1.Company{}).
		Complete(r)
}
