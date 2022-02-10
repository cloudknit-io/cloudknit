package repo

import (
	"context"
	"github.com/go-logr/logr"
	"strconv"

	"github.com/compuzest/zlifecycle-il-operator/controllers/util/env"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/github"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/compuzest/zlifecycle-il-operator/controllers/argocd"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/common"
	kClient "sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	GitAuthMethodSSH    = "ssh"
	AuthMethodGithubApp = "githubApp"
)

//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=../../mocks/mock_repo_api.go -package=mocks "github.com/compuzest/zlifecycle-il-operator/controllers/util/repo" API

type API interface {
	TryRegisterRepo(repoURL string) error
}

func New(ctx context.Context, c kClient.Client, mode string, log *logrus.Entry) (API, error) {
	switch mode {
	case AuthMethodGithubApp:
		return newGitHubAppService(ctx, c, log)
	case GitAuthMethodSSH:
		return newSSHService(ctx, c, log), nil
	default:
		return nil, errors.Errorf("invalid mode: %s", mode)
	}
}

type githubAppService struct {
	ctx       context.Context
	c         kClient.Client
	log       *logrus.Entry
	argocdAPI argocd.API
	gitAppAPI github.AppAPI
}

func newGitHubAppService(
	ctx context.Context,
	c kClient.Client,
	log *logrus.Entry,
) (*githubAppService, error) {
	// github app svc
	privateKey, err := common.GetSSHPrivateKey(
		ctx,
		c,
		kClient.ObjectKey{Name: env.Config.GitHubAppSecretName, Namespace: env.Config.GitHubAppSecretNamespace},
	)
	if err != nil {
		return nil, errors.Wrap(err, "error getting ssh private key for github app")
	}
	gitAppAPI, err := github.NewHTTPAppClient(ctx, privateKey, env.Config.GitHubAppID)
	if err != nil {
		return nil, errors.Wrap(err, "error instantiating github apps client")
	}
	// argocd
	argocdAPI := argocd.NewHTTPClient(ctx, logr.FromContext(ctx), env.Config.ArgocdServerURL)
	return &githubAppService{
		ctx:       ctx,
		c:         c,
		log:       log,
		argocdAPI: argocdAPI,
		gitAppAPI: gitAppAPI,
	}, nil
}

var _ API = (*githubAppService)(nil)

// TryRegisterRepo registers a git repo in ArgoCD using a GitHubApp auth.
func (s *githubAppService) TryRegisterRepo(repoURL string) error {
	key := kClient.ObjectKey{Namespace: env.Config.GitHubAppSecretNamespace, Name: env.Config.GitHubAppSecretName}
	privateKey, err := common.GetSSHPrivateKey(s.ctx, s.c, key)
	if err != nil {
		return err
	}

	int64AppID, err := strconv.ParseInt(env.Config.GitHubAppID, 10, 64)
	if err != nil {
		return errors.Wrap(err, "error parsing app id as int64")
	}

	owner, repo, err := common.ParseRepositoryInfo(repoURL)
	if err != nil {
		return errors.Wrapf(err, "error parsing owner and repo name from git url: [%s]", repoURL)
	}

	installationID, err := github.GetAppInstallationID(s.log, s.gitAppAPI, owner, repo)
	if err != nil {
		return errors.Wrapf(err, "error getting github app installation id for repo %s/%s", owner, repo)
	}
	if installationID == nil {
		return errors.New("invalid installation id: nil")
	}

	repoOpts := argocd.RepoOpts{
		RepoURL:                 repoURL,
		GitHubAppID:             int64AppID,
		GitHubAppInstallationID: *installationID,
		GitHubAppPrivateKey:     privateKey,
		Mode:                    argocd.ModeGitHubApp,
	}
	if _, err := argocd.RegisterRepo(s.log, s.argocdAPI, &repoOpts); err != nil {
		return err
	}

	return nil
}

type sshService struct {
	ctx       context.Context
	c         kClient.Client
	log       *logrus.Entry
	argocdAPI argocd.API
}

func newSSHService(ctx context.Context, c kClient.Client, log *logrus.Entry) *sshService {
	argocdAPI := argocd.NewHTTPClient(ctx, logr.FromContext(ctx), env.Config.ArgocdServerURL)
	return &sshService{
		ctx:       ctx,
		c:         c,
		log:       log,
		argocdAPI: argocdAPI,
	}
}

// TryRegisterRepo registers a git repo in ArgoCD using a private SSH key for auth.
func (s *sshService) TryRegisterRepo(repoURL string) error {
	key := kClient.ObjectKey{Namespace: env.Config.KubernetesServiceNamespace, Name: env.Config.GitSSHSecretName}
	privateKey, err := common.GetSSHPrivateKey(s.ctx, s.c, key)
	if err != nil {
		return err
	}

	repoOpts := argocd.RepoOpts{
		RepoURL:       repoURL,
		SSHPrivateKey: string(privateKey),
	}
	if _, err := argocd.RegisterRepo(s.log, s.argocdAPI, &repoOpts); err != nil {
		return err
	}

	return nil
}

var _ API = (*sshService)(nil)
