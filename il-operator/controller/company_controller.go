package controller

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/compuzest/zlifecycle-il-operator/controller/common/apm"
	"github.com/compuzest/zlifecycle-il-operator/controller/common/cloudknitservice"
	git2 "github.com/compuzest/zlifecycle-il-operator/controller/common/git"
	"github.com/compuzest/zlifecycle-il-operator/controller/common/git/gogit"
	"github.com/compuzest/zlifecycle-il-operator/controller/services/operations/argocd"
	"github.com/compuzest/zlifecycle-il-operator/controller/services/operations/git"
	"github.com/compuzest/zlifecycle-il-operator/controller/services/operations/github"

	"github.com/compuzest/zlifecycle-il-operator/controller/services/watcherservices"
	"github.com/compuzest/zlifecycle-il-operator/controller/services/zerrors"

	"github.com/compuzest/zlifecycle-il-operator/controller/util"

	"github.com/compuzest/zlifecycle-il-operator/controller/codegen/file"
	"github.com/compuzest/zlifecycle-il-operator/controller/env"
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

	cloudKnitServiceClient := cloudknitservice.NewService(env.Config.ZLifecycleAPIURL)
	organization, err := cloudKnitServiceClient.Get(ctx, env.Config.CompanyName, r.LogV2)

	if err != nil {
		companyErr := zerrors.NewCompanyError(
			company.Spec.CompanyName,
			perrors.Wrap(err, "error instantiating cloudKnit service"),
		)
		return ctrl.Result{}, r.APM.NoticeError(tx, r.LogV2, companyErr)
	}

	if len(organization.GitHubOrgName) != 0 {
		env.Config.GitHubCompanyOrganization = organization.GitHubOrgName
	}
	if len(organization.GitHubRepo) != 0 {
		env.Config.GitHubRepoURL = organization.GitHubRepo
	} else {
		env.Config.GitHubRepoURL = company.Spec.ConfigRepo.Source
	}

	// vars
	companyRepoURL := env.Config.GitHubRepoURL
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
	gitClient, err := gogit.NewGoGit(apmCtx, &gogit.Options{Mode: gogit.ModeToken, Token: token})
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

	commitInfo := git2.CommitInfo{
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

	if _, err := argocd.TryCreateProject(apmCtx, watcherServices.ArgocdClient, r.LogV2, company.Spec.CompanyName, env.Config.GitHubCompanyOrganization); err != nil {
		companyErr := zerrors.NewCompanyError(company.Spec.CompanyName, perrors.Wrap(err, "error trying to create argocd project"))
		return ctrl.Result{}, r.APM.NoticeError(tx, r.LogV2, companyErr)
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

	r.LogV2.Info("Registering helm chart repo")
	return services.InternalWatcher.Watch(env.Config.GitHelmChartsRepository)
}

func (r *CompanyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&stablev1.Company{}).
		Complete(r)
}
