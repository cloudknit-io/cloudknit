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
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/env"
	"github.com/sirupsen/logrus"
	"sync"
	"time"

	"github.com/compuzest/zlifecycle-il-operator/controllers/util/git"

	"github.com/compuzest/zlifecycle-il-operator/controllers/util/file"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/repo"
	"k8s.io/apimachinery/pkg/api/errors"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	stablev1 "github.com/compuzest/zlifecycle-il-operator/api/v1"
	"github.com/compuzest/zlifecycle-il-operator/controllers/argocd"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/github"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/il"
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

	// init logic
	var initOperatorError error
	initOperatorLock.Do(func() {
		initOperatorError = r.initCompany(ctx)
	})
	if initOperatorError != nil {
		return ctrl.Result{}, initOperatorError
	}

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
		r.LogV2.WithField("name", txName).Info("Creating APM transaction")
		defer tx.End()
	}

	// services init
	fileAPI := file.NewOsFileService()
	argocdAPI := argocd.NewHTTPClient(apmCtx, r.Log, argocdServerURL)
	repoAPI := github.NewHTTPRepositoryAPI(apmCtx)
	gitAPI, err := git.NewGoGit(apmCtx)
	if err != nil {
		return ctrl.Result{}, err
	}

	// temp clone IL repo
	tempILRepoDir, cleanup, err := git.CloneTemp(gitAPI, zlILRepoURL)
	if err != nil {
		return ctrl.Result{}, err
	}
	defer cleanup()

	// reconcile logic
	companyRepo := company.Spec.ConfigRepo.Source
	if err := repo.TryRegisterRepo(apmCtx, r.Client, r.Log, argocdAPI, companyRepo, operatorNamespace, operatorSSHSecret); err != nil {
		return ctrl.Result{}, err
	}

	if err := generateAndSaveCompanyApp(fileAPI, company, tempILRepoDir); err != nil {
		return ctrl.Result{}, err
	}

	if err := generateAndSaveCompanyConfigWatcher(fileAPI, company, tempILRepoDir); err != nil {
		return ctrl.Result{}, err
	}

	_, err = github.TryCreateRepository(r.Log, repoAPI, ilRepoOwner, zlILRepoName)
	if err != nil {
		return ctrl.Result{}, err
	}

	if err := repo.TryRegisterRepo(apmCtx, r.Client, r.Log, argocdAPI, zlILRepoURL, operatorNamespace, operatorSSHSecret); err != nil {
		return ctrl.Result{}, err
	}

	commitInfo := git.CommitInfo{
		Author: githubSvcAccntName,
		Email:  githubSvcAccntEmail,
		Msg:    fmt.Sprintf("Reconciling company %s", company.Spec.CompanyName),
	}
	pushed, err := gitAPI.CommitAndPush(&commitInfo)
	if err != nil {
		return ctrl.Result{}, err
	}
	if pushed {
		r.LogV2.Info("Committed new changes to IL repo")
	} else {
		r.LogV2.Info("No git changes to commit, no-op reconciliation.")
	}

	if err := argocd.TryCreateBootstrapApps(apmCtx, r.Log); err != nil {
		return ctrl.Result{}, err
	}

	// try create a webhook, it will fail if git service account does not have permissions to create it
	_, err = github.CreateRepoWebhook(r.Log, repoAPI, companyRepo, argocdHookURL, gitHubWebhookSecret)
	if err != nil {
		r.LogV2.WithError(err).Error("error creating Company webhook")
	}

	duration := time.Since(start)
	r.LogV2.WithField("duration", duration).Info("Reconcile finished")

	return ctrl.Result{}, nil
}

func (r *CompanyReconciler) initCompany(ctx context.Context) error {
	r.LogV2.Info("Running company operator init")

	r.LogV2.Info("Creating webhook for IL repo")
	repoAPI := github.NewHTTPRepositoryAPI(ctx)
	if _, err := github.CreateRepoWebhook(r.Log, repoAPI, zlILRepoURL, argocdHookURL, gitHubWebhookSecret); err != nil {
		r.LogV2.WithError(err).WithField("repo", zlILRepoURL).Error("error creating Company IL webhook")
	}

	r.LogV2.Info("Registering helm chart repo")
	argocdAPI := argocd.NewHTTPClient(ctx, r.Log, argocdServerURL)

	return repo.TryRegisterRepo(ctx, r.Client, r.Log, argocdAPI, helmChartsRepo, operatorNamespace, operatorSSHSecret)
}

func (r *CompanyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&stablev1.Company{}).
		Complete(r)
}

func generateAndSaveCompanyApp(fileAPI file.FSAPI, company *stablev1.Company, ilRepoDir string) error {
	companyApp := argocd.GenerateCompanyApp(company)

	return fileAPI.SaveYamlFile(*companyApp, il.CompanyDirectoryAbsolutePath(ilRepoDir), company.Spec.CompanyName+".yaml")
}

func generateAndSaveCompanyConfigWatcher(fileAPI file.FSAPI, company *stablev1.Company, ilRepoDir string) error {
	companyConfigWatcherApp := argocd.GenerateCompanyConfigWatcherApp(company.Spec.CompanyName, company.Spec.ConfigRepo.Source)

	return fileAPI.SaveYamlFile(*companyConfigWatcherApp, il.ConfigWatcherDirectoryAbsolutePath(ilRepoDir), company.Spec.CompanyName+".yaml")
}
