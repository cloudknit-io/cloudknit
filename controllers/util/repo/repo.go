package repo

import (
	"context"

	"github.com/compuzest/zlifecycle-il-operator/controllers/util/env"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/github"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/compuzest/zlifecycle-il-operator/controllers/argocd"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/common"
	kClient "sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	AuthMethodSSH       = "ssh"
	AuthMethodGithubApp = "githubApp"
)

//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=../../../mocks/mock_repo_api.go -package=mocks "github.com/compuzest/zlifecycle-il-operator/controllers/util/repo" Registration

type Registration interface {
	TryRegisterRepo(repoURL string) error
}

func NewRegistration(ctx context.Context, c kClient.Client, mode string, log *logrus.Entry) (Registration, error) {
	argocdAPI := argocd.NewHTTPClient(ctx, logr.FromContextOrDiscard(ctx), env.Config.ArgocdServerURL)
	switch mode {
	case AuthMethodGithubApp:
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
		return NewGitHubAppService(ctx, c, gitAppAPI, argocdAPI, log)
	case AuthMethodSSH:
		return NewSSHService(ctx, c, argocdAPI, log), nil
	default:
		return nil, errors.Errorf("invalid mode: %s", mode)
	}
}

type GitHubAppService struct {
	ctx       context.Context
	c         kClient.Client
	log       *logrus.Entry
	argocdAPI argocd.API
	gitAppAPI github.AppAPI
}

func NewGitHubAppService(
	ctx context.Context,
	c kClient.Client,
	gitAppAPI github.AppAPI,
	argocdAPI argocd.API,
	log *logrus.Entry,
) (*GitHubAppService, error) {
	return &GitHubAppService{
		ctx:       ctx,
		c:         c,
		log:       log,
		argocdAPI: argocdAPI,
		gitAppAPI: gitAppAPI,
	}, nil
}

var _ Registration = (*GitHubAppService)(nil)

// TryRegisterRepo registers a git repo in ArgoCD using a GitHubApp auth.
func (s *GitHubAppService) TryRegisterRepo(repoURL string) error {
	key := kClient.ObjectKey{Namespace: env.Config.GitHubAppSecretNamespace, Name: env.Config.GitHubAppSecretName}
	privateKey, err := common.GetSSHPrivateKey(s.ctx, s.c, key)
	if err != nil {
		return err
	}

	owner, repo, err := common.ParseRepositoryInfo(repoURL)
	if err != nil {
		return errors.Wrapf(err, "error parsing owner and repo name from git url: [%s]", repoURL)
	}

	installationID, appID, err := github.GetAppInstallationID(s.log, s.gitAppAPI, owner, repo)
	if err != nil {
		return errors.Wrapf(err, "error getting github app installation id for repo %s/%s", owner, repo)
	}
	if installationID == nil {
		return errors.New("invalid installation id: nil")
	}

	repoOpts := argocd.RepoOpts{
		RepoURL:                 repoURL,
		GitHubAppID:             *appID,
		GitHubAppInstallationID: *installationID,
		GitHubAppPrivateKey:     privateKey,
		Mode:                    argocd.ModeGitHubApp,
	}
	if _, err := argocd.RegisterRepo(s.log, s.argocdAPI, &repoOpts); err != nil {
		return err
	}

	return nil
}

type SSHService struct {
	ctx       context.Context
	c         kClient.Client
	log       *logrus.Entry
	argocdAPI argocd.API
}

func NewSSHService(ctx context.Context, c kClient.Client, argocdAPI argocd.API, log *logrus.Entry) *SSHService {
	return &SSHService{
		ctx:       ctx,
		c:         c,
		log:       log,
		argocdAPI: argocdAPI,
	}
}

// TryRegisterRepo registers a git repo in ArgoCD using a private SSH key for auth.
func (s *SSHService) TryRegisterRepo(repoURL string) error {
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

var _ Registration = (*SSHService)(nil)
