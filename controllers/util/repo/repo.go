package repo

import (
	"context"
	"strconv"

	"github.com/compuzest/zlifecycle-il-operator/controllers/util/env"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/github"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/compuzest/zlifecycle-il-operator/controllers/argocd"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/common"
	kClient "sigs.k8s.io/controller-runtime/pkg/client"
)

// TryRegisterRepoViaSSHPrivateKey registers a git repo in ArgoCD using a private SSH key for auth.
func TryRegisterRepoViaSSHPrivateKey(
	ctx context.Context,
	c kClient.Client,
	log *logrus.Entry,
	api argocd.API,
	repoURL string,
	secretNamespace string,
	secretName string,
) error {
	key := kClient.ObjectKey{Namespace: secretNamespace, Name: secretName}
	privateKey, err := common.GetSSHPrivateKey(ctx, c, key)
	if err != nil {
		return err
	}

	repoOpts := argocd.RepoOpts{
		RepoURL:       repoURL,
		SSHPrivateKey: string(privateKey),
	}
	if _, err := argocd.RegisterRepo(log, api, &repoOpts); err != nil {
		return err
	}

	return nil
}

// TryRegisterRepoViaGitHubApp registers a git repo in ArgoCD using a GitHubApp auth.
func TryRegisterRepoViaGitHubApp(
	ctx context.Context,
	c kClient.Client,
	log *logrus.Entry,
	argocdAPI argocd.API,
	gitAppAPI github.AppAPI,
	repoURL string,
	secretNamespace string,
	secretName string,
) error {
	key := kClient.ObjectKey{Namespace: secretNamespace, Name: secretName}
	privateKey, err := common.GetSSHPrivateKey(ctx, c, key)
	if err != nil {
		return err
	}

	owner, repo, err := common.ParseRepositoryInfo(repoURL)
	if err != nil {
		return errors.Wrapf(err, "error parsing owner and repo name from git url: [%s]", repoURL)
	}

	installationID, err := github.GetAppInstallationID(log, gitAppAPI, owner, repo)
	if err != nil {
		return errors.Wrapf(err, "error getting github app installation id for repo %s/%s", owner, repo)
	}
	if installationID == nil {
		return errors.New("invalid installation id: nil")
	}

	repoOpts := argocd.RepoOpts{
		RepoURL:                 repoURL,
		GitHubAppID:             env.Config.GitHubAppID,
		GitHubAppInstallationID: strconv.FormatInt(*installationID, 10),
		GitHubAppPrivateKey:     privateKey,
	}
	if _, err := argocd.RegisterRepo(log, argocdAPI, &repoOpts); err != nil {
		return err
	}

	return nil
}
