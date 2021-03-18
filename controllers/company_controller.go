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
	argocd "github.com/compuzest/zlifecycle-il-operator/controllers/argocd"
	env "github.com/compuzest/zlifecycle-il-operator/controllers/util/env"
	fileutil "github.com/compuzest/zlifecycle-il-operator/controllers/util/file"
	github "github.com/compuzest/zlifecycle-il-operator/controllers/util/github"
	il "github.com/compuzest/zlifecycle-il-operator/controllers/util/il"
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
	ctx := context.Background()

	company := &stablev1alpha1.Company{}
	r.Get(ctx, req.NamespacedName, company)

	companyRepo := company.Spec.ConfigRepo.Source
	repoSecret := company.Spec.ConfigRepo.RepoSecret
	if err := argocd.TryRegisterRepo(r.Client, r.Log, ctx, companyRepo, req.Namespace, repoSecret); err != nil {
		return ctrl.Result{}, err
	}

	if err := generateAndSaveCompanyApp(company); err != nil {
		return ctrl.Result{}, err
	}

	if err := generateAndSaveCompanyConfigWatcher(company); err != nil {
		return ctrl.Result{}, err
	}

	ilRepo := company.Name + "-il"
	githubRepositoryApi := github.NewHttpClient(env.Config.GitHubAuthToken, ctx)
	_, err := github.TryCreateRepository(r.Log, githubRepositoryApi, env.Config.ZlifecycleOwner, ilRepo)
	if err != nil {
		return ctrl.Result{}, err
	}

	if err := argocd.TryRegisterRepo(r.Client, r.Log, ctx, ilRepo, req.Namespace, repoSecret); err != nil {
		return ctrl.Result{}, err
	}

	github.CommitAndPushFiles(
		env.Config.ILRepoSourceOwner,
		env.Config.ILRepoName,
		[]string{il.Config.CompanyDirectory, il.Config.ConfigWatcherDirectory},
		env.Config.RepoBranch,
		fmt.Sprintf("Reconciling company %s", company.Spec.CompanyName),
		env.Config.GithubSvcAccntName,
		env.Config.GithubSvcAccntEmail)

	return ctrl.Result{}, nil
}

func (r *CompanyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&stablev1alpha1.Company{}).
		Complete(r)
}

func generateAndSaveCompanyApp(company *stablev1alpha1.Company) error {
	companyApp := argocd.GenerateCompanyApp(*company)

	if err := fileutil.SaveYamlFile(*companyApp, il.Config.CompanyDirectory, company.Spec.CompanyName+".yaml"); err != nil {
		return err
	}

	return nil
}

func generateAndSaveCompanyConfigWatcher(company *stablev1alpha1.Company) error {
	companyConfigWatcherApp := argocd.GenerateCompanyConfigWatcherApp(company.Spec.CompanyName, company.Spec.ConfigRepo.Source)

	if err := fileutil.SaveYamlFile(*companyConfigWatcherApp, il.Config.ConfigWatcherDirectory, company.Spec.CompanyName+".yaml"); err != nil {
		return err
	}

	return nil
}
