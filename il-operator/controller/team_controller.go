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
	"github.com/compuzest/zlifecycle-il-operator/controller/services/operations/git"
	"github.com/compuzest/zlifecycle-il-operator/controller/services/operations/github"

	"github.com/compuzest/zlifecycle-il-operator/controller/services/watcherservices"
	"github.com/compuzest/zlifecycle-il-operator/controller/services/zerrors"

	"github.com/compuzest/zlifecycle-il-operator/controller/util"

	"github.com/compuzest/zlifecycle-il-operator/controller/codegen/file"
	"github.com/compuzest/zlifecycle-il-operator/controller/codegen/il"

	"github.com/compuzest/zlifecycle-il-operator/controller/env"
	"github.com/sirupsen/logrus"

	perrors "github.com/pkg/errors"

	"github.com/go-logr/logr"
	"go.uber.org/atomic"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
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

	r.LogV2 = r.LogV2.WithField("team", team.Spec.TeamName)

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

	// services init
	fileAPI := file.NewOSFileService()
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
	gitClient, err := gogit.NewGoGit(apmCtx, &gogit.Options{Mode: gogit.ModeToken, Token: token})
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

	if err := fileAPI.CreateEmptyDirectory(il.EnvironmentDirectoryPath(team.Spec.TeamName)); err != nil {
		teamErr := zerrors.NewTeamError(team.Spec.TeamName, perrors.Wrap(err, "error creating team dir"))
		return ctrl.Result{}, r.APM.NoticeError(tx, r.LogV2, teamErr)
	}

	teamRepoURL := team.Spec.ConfigRepo.Source

	if err := watcherServices.CompanyWatcher.Watch(teamRepoURL); err != nil {
		teamErr := zerrors.NewTeamError(team.Spec.TeamName, perrors.Wrap(err, "error registering argocd team repo via github app auth"))
		return ctrl.Result{}, r.APM.NoticeError(tx, r.LogV2, teamErr)
	}

	teamAppFilename := fmt.Sprintf("%s-team.yaml", team.Spec.TeamName)

	cloudKnitServiceClient := cloudknitservice.NewService(env.Config.ZLifecycleAPIURL)
	err = cloudKnitServiceClient.PostTeam(ctx, env.Config.CompanyName, *team, r.LogV2)

	if err != nil {
		teamErr := zerrors.NewTeamError(team.Spec.TeamName, perrors.Wrap(err, "error with the POST team call"))
		return ctrl.Result{}, r.APM.NoticeError(tx, r.LogV2, teamErr)
	}

	if err := generateAndSaveConfigWatchers(fileAPI, team, teamAppFilename, tempILRepoDir); err != nil {
		teamErr := zerrors.NewTeamError(team.Spec.TeamName, perrors.Wrap(err, "error generating team config watchers"))
		return ctrl.Result{}, r.APM.NoticeError(tx, r.LogV2, teamErr)
	}

	commitInfo := git2.CommitInfo{
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
		webhookURL := env.Config.ArgocdWebhookURL + "?t=" + team.Spec.TeamName
		_, err = github.CreateRepoWebhook(r.LogV2, watcherServices.CompanyGitClient, teamRepoURL, webhookURL, env.Config.GitHubWebhookSecret)
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

func delayTeamReconcileOnInitialRun(log *logrus.Entry, seconds int64) {
	if teamReconcileInitialRunLock.Load() {
		log.WithField("duration", fmt.Sprintf("%ds", seconds)).Info("Delaying Team reconcile on initial run to wait for Company operator")
		time.Sleep(time.Duration(seconds) * time.Second)
		teamReconcileInitialRunLock.Store(false)
	}
}
