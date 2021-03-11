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
	"errors"
	"fmt"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"

	stablev1alpha1 "github.com/compuzest/zlifecycle-il-operator/api/v1alpha1"
	env "github.com/compuzest/zlifecycle-il-operator/controllers/util/env"
	github "github.com/compuzest/zlifecycle-il-operator/controllers/util/github"
	il "github.com/compuzest/zlifecycle-il-operator/controllers/util/il"

	coreV1 "k8s.io/api/core/v1"

	argocd "github.com/compuzest/zlifecycle-il-operator/controllers/argocd"
	fileutil "github.com/compuzest/zlifecycle-il-operator/controllers/util/file"
)

// TeamReconciler reconciles a Team object
type TeamReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups="",resources=configmaps;secrets,verbs=get;list;watch;update;patch
// +kubebuilder:rbac:groups=stable.compuzest.com,resources=teams,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=stable.compuzest.com,resources=teams/status,verbs=get;update;patch

// Reconcile method called everytime there is a change in Team Custom Resource
func (r *TeamReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()

	team := &stablev1alpha1.Team{}
	r.Get(ctx, req.NamespacedName, team)

	teamYAML := fmt.Sprintf("%s-team.yaml", team.Spec.TeamName)

	if err := fileutil.CreateEmptyDirectory(il.EnvironmentDirectory(team.Spec.TeamName)); err != nil {
		return ctrl.Result{}, err
	}

	teamRepo := team.Spec.ConfigRepo.Source
	repoSecret := team.Spec.ConfigRepo.RepoSecret
	if err := tryRegisterTeamRepo(r.Client, r.Log, ctx, teamRepo, req.Namespace, repoSecret); err != nil {
		return ctrl.Result{}, err
	}

	if err := generateAndSaveTeamApp(team, teamYAML); err != nil {
		return ctrl.Result{}, err
	}

	if err := generateAndSaveConfigWatchers(team, teamYAML); err != nil {
		return ctrl.Result{}, err
	}

	// Avoid race condition on initial Reconcile, collides with Team controller commit
	time.Sleep(5 * time.Second)

	if err := github.CommitAndPushFiles(
		env.Config.ILRepoSourceOwner,
		env.Config.ILRepoName,
		[]string{il.Config.TeamDirectory, il.Config.ConfigWatcherDirectory},
		env.Config.RepoBranch,
		fmt.Sprintf("Reconciling team %s", team.Spec.TeamName),
		env.Config.GithubSvcAccntName,
		env.Config.GithubSvcAccntEmail); err != nil {
		return ctrl.Result{}, err
	}

	githubApi := github.NewHttpClient(env.Config.GitHubAuthToken, ctx)
	_, err := github.CreateRepoWebhook(r.Log, githubApi, teamRepo, env.Config.ArgocdHookUrl)
	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func tryRegisterTeamRepo(
	c client.Client,
	log logr.Logger,
	ctx context.Context,
	teamRepo string,
	namespace string,
	repoSecret string,
) error {
	secret := &coreV1.Secret{}
	secretNamespacedName :=
		types.NamespacedName{Namespace: namespace, Name: repoSecret}
	if err := c.Get(ctx, secretNamespacedName, secret); err != nil {
		log.Info(
			"Secret %s does not exist in namespace %s\n",
			repoSecret,
			namespace,
		)
		return err
	}

	sshPrivateKeyField := "sshPrivateKey"
	sshPrivateKey := string(secret.Data[sshPrivateKeyField])
	if sshPrivateKey == "" {
		errMsg := fmt.Sprintf("Secret is missing %s data field!", sshPrivateKeyField)
		err := errors.New(errMsg)
		log.Error(err, errMsg)
		return err
	}

	repoOpts := argocd.RepoOpts{
		RepoUrl:       teamRepo,
		SshPrivateKey: sshPrivateKey,
	}

	argocdApi := argocd.NewHttpClient(log, argocd.GetArgocdServerAddr())
	if _, err := argocd.RegisterRepo(log, argocdApi, repoOpts); err != nil {
		log.Error(err, "Error while calling ArgoCD Repo API")
		return err
	}

	return nil
}

// SetupWithManager sets up the Company Controller with Manager
func (r *TeamReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&stablev1alpha1.Team{}).
		Complete(r)
}

func generateAndSaveTeamApp(team *stablev1alpha1.Team, teamYAML string) error {
	teamApp := argocd.GenerateTeamApp(*team)

	if err := fileutil.SaveYamlFile(*teamApp, il.Config.TeamDirectory, teamYAML); err != nil {
		return err
	}

	return nil
}

func generateAndSaveConfigWatchers(team *stablev1alpha1.Team, teamYAML string) error {
	teamConfigWatcherApp := argocd.GenerateTeamConfigWatcherApp(*team)

	if err := fileutil.SaveYamlFile(*teamConfigWatcherApp, il.Config.ConfigWatcherDirectory, teamYAML); err != nil {
		return err
	}

	return nil
}
