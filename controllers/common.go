package controllers

import (
	"context"
	"strconv"
	"strings"

	"github.com/compuzest/zlifecycle-il-operator/controllers/external/argocd"
	"github.com/compuzest/zlifecycle-il-operator/controllers/external/github"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util"
	"github.com/compuzest/zlifecycle-il-operator/controllers/watcher"
	perrors "github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	kClient "sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/compuzest/zlifecycle-il-operator/controllers/env"
)

func checkIsNamespaceWatched(namespace string) bool {
	watchedNamespace := env.Config.KubernetesOperatorWatchedNamespace
	return namespace == watchedNamespace
}

func checkIsResourceWatched(resource string) bool {
	watchedResources := strings.Split(env.Config.KubernetesOperatorWatchedResources, ",")

	for _, r := range watchedResources {
		if strings.EqualFold(strings.TrimSpace(r), resource) {
			return true
		}
	}

	return false
}

type WatcherServices struct {
	argocdClient      argocd.API
	companyGitClient  github.API
	internalGitClient github.API
	ilGitClient       github.API
	companyWatcher    watcher.API
	internalWatcher   watcher.API
	ilWatcher         watcher.API
}

func newGitHubServices(ctx context.Context, client kClient.Client, gitOrg string, log *logrus.Entry) (*WatcherServices, error) {
	argocdClient := argocd.NewHTTPClient(ctx, log, env.Config.ArgocdServerURL)

	// get private keys
	companyPrivateKey, err := util.GetAuthTierModePrivateKey(ctx, client, env.Config.GitHubCompanyAuthMethod, util.AuthTierCompany)
	if err != nil {
		return nil, perrors.Wrap(err, "error getting company private key")
	}
	internalPrivateKey, err := util.GetAuthTierModePrivateKey(ctx, client, env.Config.GitHubInternalAuthMethod, util.AuthTierInternal)
	if err != nil {
		return nil, perrors.Wrap(err, "error getting internal private key")
	}

	companyGitClient, err := newGitHubClient(ctx, util.AuthTierCompany, env.Config.GitHubCompanyAuthMethod, companyPrivateKey, gitOrg, log)
	if err != nil {
		return nil, perrors.Wrap(err, "error instantiating company git client")
	}
	internalGitClient, err := newGitHubClient(
		ctx, util.AuthTierInternal, env.Config.GitHubInternalAuthMethod, internalPrivateKey, env.Config.GitCoreRepositoryOwner, log,
	)
	if err != nil {
		return nil, perrors.Wrap(err, "error instantiating internal git client")
	}
	ilGitClient, err := newGitHubClient(
		ctx, util.AuthTierInternal, env.Config.GitHubInternalAuthMethod, internalPrivateKey, env.Config.GitILRepositoryOwner, log,
	)
	if err != nil {
		return nil, perrors.Wrap(err, "error instantiating il git client")
	}

	// watchers
	companyWatcher, err := newWatcher(ctx, argocdClient, companyPrivateKey, env.Config.GitHubCompanyAuthMethod, util.AuthTierCompany, gitOrg, log)
	if err != nil {
		return nil, perrors.Wrap(err, "error instantiating company watcher")
	}
	internalWatcher, err := newWatcher(
		ctx, argocdClient, internalPrivateKey, env.Config.GitHubInternalAuthMethod, util.AuthTierInternal, env.Config.GitCoreRepositoryOwner, log,
	)
	if err != nil {
		return nil, perrors.Wrap(err, "error instantiating internal watcher")
	}
	ilWatcher, err := newWatcher(
		ctx, argocdClient, internalPrivateKey, env.Config.GitHubInternalAuthMethod, util.AuthTierInternal, env.Config.GitILRepositoryOwner, log,
	)
	if err != nil {
		return nil, perrors.Wrap(err, "error instantiating il watcher")
	}

	return &WatcherServices{
		argocdClient:      argocdClient,
		companyGitClient:  companyGitClient,
		internalGitClient: internalGitClient,
		ilGitClient:       ilGitClient,
		companyWatcher:    companyWatcher,
		internalWatcher:   internalWatcher,
		ilWatcher:         ilWatcher,
	}, nil
}

func newGitHubClient(
	ctx context.Context,
	tier util.AuthTier,
	mode util.AuthMode,
	privateKey []byte,
	gitOrg string,
	log *logrus.Entry,
) (github.API, error) {
	b := github.NewClientBuilder()

	switch mode {
	case util.AuthModeGitHubApp:
		appID, err := getGitHubAppID(tier)
		if err != nil {
			return nil, perrors.Wrap(err, "error getting github app id")
		}
		client, err := b.WithGitHubApp(ctx, privateKey, appID).Build()
		if err != nil {
			return nil, perrors.Wrap(err, "error instantiating github client with github app auth method")
		}
		if gitOrg == "" {
			return client, nil
		}
		return newGitHubAppClientWithInstallationID(ctx, client, privateKey, gitOrg, log)
	case util.AuthModeSSH:
		return b.WithToken(ctx, env.Config.GitToken).Build()
	default:
		return nil, perrors.Errorf("invalid auth mode [%s]", mode)
	}
}

func newGitHubAppClientWithInstallationID(
	ctx context.Context,
	client github.API,
	privateKey []byte,
	gitOrg string,
	log *logrus.Entry,
) (github.API, error) {
	installationID, appID, err := github.GetAppInstallationID(log, client, gitOrg)
	if err != nil {
		return nil, perrors.Wrapf(err, "error getting github app installation id for git organization [%s]", gitOrg)
	}

	return github.NewClientBuilder().WithGitHubApp(
		ctx, privateKey, strconv.FormatInt(*appID, 10),
	).WithInstallationID(
		strconv.FormatInt(*installationID, 10),
	).Build()
}

func getGitHubAppID(tier util.AuthTier) (string, error) {
	switch tier {
	case util.AuthTierCompany:
		return env.Config.GitHubAppIDCompany, nil
	case util.AuthTierInternal:
		return env.Config.GitHubAppIDInternal, nil
	default:
		return "", perrors.Errorf("invalid auth tier: %s", tier)
	}
}

func newWatcher(
	ctx context.Context,
	argocdClient argocd.API,
	privateKey []byte,
	mode util.AuthMode,
	tier util.AuthTier,
	gitOrg string,
	log *logrus.Entry,
) (watcher.API, error) {
	switch mode {
	case util.AuthModeGitHubApp:
		gitClient, err := newGitHubClient(ctx, tier, mode, privateKey, "", log)
		if err != nil {
			return nil, perrors.Wrapf(err, "error instantiating git client for mode [%s] and tier [%s]", mode, tier)
		}
		installation, resp, err := gitClient.FindOrganizationInstallation(gitOrg)
		if err != nil {
			return nil, perrors.Wrapf(err, "error finding installation for organization [%s]", gitOrg)
		}
		defer util.CloseBody(resp.Body)
		return watcher.NewGitHubAppWatcher(ctx, *installation.AppID, *installation.ID, argocdClient, privateKey, log)
	case util.AuthModeSSH:
		return watcher.NewSSHWatcher(ctx, argocdClient, privateKey, log), nil
	default:
		return nil, perrors.Errorf("invalid auth mode: %s", mode)
	}
}
