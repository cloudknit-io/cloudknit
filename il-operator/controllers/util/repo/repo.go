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
	AuthTierCompany     = "company"
	AuthTierInternal    = "internal"
)

//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=../../../mocks/mock_repo_api.go -package=mocks "github.com/compuzest/zlifecycle-il-operator/controllers/util/repo" Registration

type Registration interface {
	TryRegisterRepo(repoURL string) error
}

func NewRegistration(ctx context.Context, c kClient.Client, mode string, tier string, log *logrus.Entry) (Registration, error) {
	argocdAPI := argocd.NewHTTPClient(ctx, logr.FromContextOrDiscard(ctx), env.Config.ArgocdServerURL)

	switch mode {
	case AuthMethodGithubApp:
		nfo, err := getGitHubAppCreds(tier)
		if err != nil {
			return nil, errors.Wrap(err, "error getting github app creds")
		}
		// github app svc
		gitAppAPI, err := newGitAppAPI(ctx, c, nfo)
		if err != nil {
			return nil, errors.Wrap(err, "error creating github app api")
		}
		return NewGitHubAppService(ctx, c, gitAppAPI, argocdAPI, log)
	case AuthMethodSSH:
		return NewSSHService(ctx, c, argocdAPI, log), nil
	default:
		return nil, errors.Errorf("invalid mode: %s", mode)
	}
}

func getGitHubAppCreds(tier string) (*secretInfo, error) {
	var secretName string
	var secretNamespace string
	var appID string

	switch tier {
	case AuthTierCompany:
		secretName = env.Config.GitHubAppSecretNameCompany
		secretNamespace = env.Config.GitHubAppSecretNamespaceCompany
		appID = env.Config.GitHubAppIDCompany
	case AuthTierInternal:
		secretName = env.Config.GitHubAppSecretNameInternal
		secretNamespace = env.Config.GitHubAppSecretNamespaceInternal
		appID = env.Config.GitHubAppIDInternal
	default:
		return nil, errors.Errorf("invalid tier: %s", tier)
	}

	return &secretInfo{
		secretName:      secretName,
		secretNamespace: secretNamespace,
		appID:           appID,
	}, nil
}

func newGitAppAPI(ctx context.Context, c kClient.Client, nfo *secretInfo) (github.AppAPI, error) {
	privateKey, err := common.GetSSHPrivateKey(
		ctx,
		c,
		kClient.ObjectKey{Name: nfo.secretName, Namespace: nfo.secretNamespace},
	)
	if err != nil {
		return nil, errors.Wrap(err, "error getting ssh private key for github app")
	}
	gitAppAPI, err := github.NewHTTPAppClient(ctx, privateKey, nfo.appID)
	if err != nil {
		return nil, errors.Wrap(err, "error instantiating github apps client")
	}

	return gitAppAPI, nil
}

type secretInfo struct {
	secretName      string
	secretNamespace string
	appID           string
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
	key := kClient.ObjectKey{Namespace: env.Config.GitHubAppSecretNamespaceCompany, Name: env.Config.GitHubAppSecretNameCompany}
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
