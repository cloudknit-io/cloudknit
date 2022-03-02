package watcher

import (
	"context"

	"github.com/compuzest/zlifecycle-il-operator/controllers/external/argocd"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

//go:generate mockgen --build_flags=--mod=mod -destination=./mock_watcher_api.go -package=watcher "github.com/compuzest/zlifecycle-il-operator/controllers/watcher" API
type API interface {
	Watch(repoURL string) error
}

type GitHubAppWatcher struct {
	ctx            context.Context
	log            *logrus.Entry
	argocdAPI      argocd.API
	appID          int64
	installationID int64
	privateKey     []byte
}

func NewGitHubAppWatcher(
	ctx context.Context,
	appID int64,
	installationID int64,
	argocdClient argocd.API,
	privateKey []byte,
	log *logrus.Entry,
) (*GitHubAppWatcher, error) {
	return &GitHubAppWatcher{
		ctx:            ctx,
		log:            log,
		appID:          appID,
		installationID: installationID,
		argocdAPI:      argocdClient,
		privateKey:     privateKey,
	}, nil
}

var _ API = (*GitHubAppWatcher)(nil)

// Watch registers a git repo in ArgoCD using a GitHubApp auth.
func (s *GitHubAppWatcher) Watch(repoURL string) error {
	repoOpts := argocd.RepoOpts{
		RepoURL:                 util.RewriteGitURLToHTTPS(repoURL),
		GitHubAppID:             s.appID,
		GitHubAppInstallationID: s.installationID,
		GitHubAppPrivateKey:     s.privateKey,
		Mode:                    argocd.RegistrationModeGithubApp,
	}
	if _, err := argocd.RegisterRepo(s.log, s.argocdAPI, &repoOpts); err != nil {
		return errors.Wrapf(err, "error registering repo in argocd using github app auth method")
	}

	return nil
}

type SSHWatcher struct {
	ctx        context.Context
	argocdAPI  argocd.API
	privateKey []byte
	log        *logrus.Entry
}

func NewSSHWatcher(ctx context.Context, argocdAPI argocd.API, privateKey []byte, log *logrus.Entry) *SSHWatcher {
	return &SSHWatcher{
		ctx:        ctx,
		privateKey: privateKey,
		argocdAPI:  argocdAPI,
		log:        log,
	}
}

// Watch registers a git repo in ArgoCD using a private SSH key for auth.
func (s *SSHWatcher) Watch(repoURL string) error {
	repoOpts := argocd.RepoOpts{
		RepoURL:       repoURL,
		SSHPrivateKey: string(s.privateKey),
	}
	if _, err := argocd.RegisterRepo(s.log, s.argocdAPI, &repoOpts); err != nil {
		return errors.Wrap(err, "error registering repo in argocd using SSH auth method")
	}

	return nil
}

var _ API = (*SSHWatcher)(nil)