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
	"k8s.io/apimachinery/pkg/api/errors"
	"time"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	stablev1 "github.com/compuzest/zlifecycle-il-operator/api/v1"
	"github.com/compuzest/zlifecycle-il-operator/controllers/argocd"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/env"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/github"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/il"
)

// CompanyReconciler reconciles a Company object
type CompanyReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=stable.compuzest.com,resources=companies,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=stable.compuzest.com,resources=companies/status,verbs=get;update;patch

func (r *CompanyReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	start := time.Now()
	ctx := context.Background()

	if err := r.initOperator(ctx); err != nil {
		r.Log.Error(err, "error running init function")
	}

	company := &stablev1.Company{}
	if err := r.Get(ctx, req.NamespacedName, company); err != nil {
		if errors.IsNotFound(err) {
			r.Log.Info(
				"Company missing from cache, ending reconcile...",
				"name", req.Name,
				"namespace", req.Namespace,
			)
			return ctrl.Result{}, nil
		}
		r.Log.Error(
			err,
			"Error occurred while getting Company...",
			"name", req.Name,
			"namespace", req.Namespace,
		)

		return ctrl.Result{}, err
	}

	companyRepo := company.Spec.ConfigRepo.Source
	operatorSshSecret := env.Config.ZlifecycleMasterRepoSshSecret
	operatorNamespace := env.Config.ZlifecycleOperatorNamespace

	argocdApi := argocd.NewHttpClient(r.Log, env.Config.ArgocdServerUrl)

	if err := repo.TryRegisterRepo(r.Client, r.Log, ctx, argocdApi, companyRepo, operatorNamespace, operatorSshSecret); err != nil {
		return ctrl.Result{}, err
	}

	if err := generateAndSaveCompanyApp(company); err != nil {
		return ctrl.Result{}, err
	}

	if err := generateAndSaveCompanyConfigWatcher(company); err != nil {
		return ctrl.Result{}, err
	}

	owner := env.Config.ZlifecycleOwner
	ilRepo := il.RepoName(company.Name)
	githubRepositoryAPI := github.NewHttpRepositoryClient(env.Config.GitHubAuthToken, ctx)
	_, err := github.TryCreateRepository(r.Log, githubRepositoryAPI, owner, ilRepo)
	if err != nil {
		return ctrl.Result{}, err
	}

	ilRepoURL := il.RepoURL(owner, company.Name)
	if err := repo.TryRegisterRepo(r.Client, r.Log, ctx, argocdApi, ilRepoURL, operatorNamespace, operatorSshSecret); err != nil {
		return ctrl.Result{}, err
	}

	dirty, err := github.CommitAndPushFiles(
		env.Config.ILRepoSourceOwner,
		ilRepo,
		[]string{il.Config.CompanyDirectory, il.Config.ConfigWatcherDirectory},
		env.Config.RepoBranch,
		fmt.Sprintf("Reconciling company %s", company.Spec.CompanyName),
		env.Config.GithubSvcAccntName,
		env.Config.GithubSvcAccntEmail,
	)
	if err != nil {
		return ctrl.Result{}, err
	}
	if dirty {
		r.Log.Info(
			"Committed new changes to IL repo",
			"company", company.Spec.CompanyName,
		)
	} else {
		r.Log.Info(
			"No git changes to commit, no-op reconciliation.",
			"company", company.Spec.CompanyName,
		)
	}

	if err := argocd.TryCreateBootstrapApps(r.Log); err != nil {
		return ctrl.Result{}, err
	}

	_, err = github.CreateRepoWebhook(r.Log, githubRepositoryAPI, companyRepo, env.Config.ArgocdHookUrl, env.Config.GitHubWebhookSecret)
	if err != nil {
		return ctrl.Result{}, err
	}

	duration := time.Since(start)
	r.Log.Info(
		"Reconcile finished",
		"duration", duration,
		"company", company.Spec.CompanyName,
	)

	return ctrl.Result{}, nil
}

func (r *CompanyReconciler) initOperator(ctx context.Context) error {
	r.Log.Info("running company operator init")
	r.Log.Info("registering helm chart repo")
	argocdApi := argocd.NewHttpClient(r.Log, env.Config.ArgocdServerUrl)
	helmChartsRepo := env.Config.HelmChartsRepo
	operatorSshSecret := env.Config.ZlifecycleMasterRepoSshSecret
	operatorNamespace := env.Config.ZlifecycleOperatorNamespace
	if err := repo.TryRegisterRepo(r.Client, r.Log, ctx, argocdApi, helmChartsRepo, operatorNamespace, operatorSshSecret); err != nil {
		return err
	}
	return nil
}

func (r *CompanyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&stablev1.Company{}).
		Complete(r)
}

func generateAndSaveCompanyApp(company *stablev1.Company) error {
	companyApp := argocd.GenerateCompanyApp(*company)
	fileUtil := &file.UtilFileService{}

	if err := fileUtil.SaveYamlFile(*companyApp, il.Config.CompanyDirectory, company.Spec.CompanyName+".yaml"); err != nil {
		return err
	}

	return nil
}

func generateAndSaveCompanyConfigWatcher(company *stablev1.Company) error {
	companyConfigWatcherApp := argocd.GenerateCompanyConfigWatcherApp(company.Spec.CompanyName, company.Spec.ConfigRepo.Source)
	fileUtil := &file.UtilFileService{}

	if err := fileUtil.SaveYamlFile(*companyConfigWatcherApp, il.Config.ConfigWatcherDirectory, company.Spec.CompanyName+".yaml"); err != nil {
		return err
	}

	return nil
}
