package github

import (
	"context"
	"io"
	"strconv"

	perrors "github.com/pkg/errors"

	"github.com/google/go-github/v42/github"
)

//go:generate mockgen --build_flags=--mod=mod -destination=./mock_github_api.go -package=github "github.com/compuzest/zlifecycle-il-operator/controllers/external/github" API
type API interface {
	FindRepositoryInstallation(owner string, repo string) (*github.Installation, *github.Response, error)
	FindOrganizationInstallation(org string) (*github.Installation, *github.Response, error)
	CreateInstallationToken(installationID int64) (*github.InstallationToken, *github.Response, error)
	CreateRepository(owner string, repo string) (*github.Repository, *github.Response, error)
	GetRepository(owner string, repo string) (*github.Repository, *github.Response, error)
	ListHooks(owner string, repo string, opts *github.ListOptions) ([]*github.Hook, *github.Response, error)
	CreateHook(owner string, repo string, hook *github.Hook) (*github.Hook, *github.Response, error)
	DownloadContents(owner string, repo string, ref string, path string) (file io.ReadCloser, exists bool, err error)
}

type Client struct {
	ctx       context.Context
	appClient *github.Client
	client    *github.Client
}

type ClientBuilder struct {
	ctx            context.Context
	token          string
	privateKey     []byte
	appID          string
	installationID string
	mode           string
}

func NewClientBuilder() *ClientBuilder {
	return &ClientBuilder{}
}

func (b *ClientBuilder) WithGitHubApp(ctx context.Context, privateKey []byte, appID string) *ClientBuilder {
	b.ctx = ctx
	b.privateKey = privateKey
	b.appID = appID
	b.mode = modeGitHubApp
	return b
}

func (b *ClientBuilder) WithInstallationID(installationID string) *ClientBuilder {
	b.installationID = installationID
	return b
}

func (b *ClientBuilder) WithToken(ctx context.Context, token string) *ClientBuilder {
	b.ctx = ctx
	b.token = token
	b.mode = modeToken
	return b
}

func (b *ClientBuilder) Build() (*Client, error) {
	switch b.mode {
	case modeGitHubApp:
		int64AppID, err := strconv.ParseInt(b.appID, 10, 64)
		if err != nil {
			return nil, perrors.Wrap(err, "error parsing app id as int64")
		}

		client := Client{
			ctx: b.ctx,
		}

		githubClient, err := newGitHubAppClient(b.privateKey, int64AppID)
		if err != nil {
			return nil, perrors.Wrap(err, "error instantiating github client without installation id")
		}
		client.appClient = githubClient
		if b.installationID != "" {
			int64InstallationID, nestedErr := strconv.ParseInt(b.installationID, 10, 64)
			if nestedErr != nil {
				return nil, perrors.Wrap(nestedErr, "error parsing app id as int64")
			}
			githubClientWithInstallation, nestedErr := newGitHubAppClientWithInstallation(b.privateKey, int64AppID, int64InstallationID)
			if nestedErr != nil {
				return nil, perrors.Wrapf(nestedErr, "error instantiating github client with installation id %s", b.installationID)
			}
			client.client = githubClientWithInstallation
		}

		if err != nil {
			return nil, perrors.Wrap(err, "error creating github app client")
		}

		return &client, nil
	case modeToken:
		return &Client{
			ctx:    b.ctx,
			client: newTokenGitHubClient(b.ctx, b.token),
		}, nil
	default:
		return nil, perrors.New("invalid configuration")
	}
}

type RepositoryAPI interface {
	CreateRepository(owner string, repo string) (*github.Repository, *github.Response, error)
	GetRepository(owner string, repo string) (*github.Repository, *github.Response, error)
	ListHooks(owner string, repo string, opts *github.ListOptions) ([]*github.Hook, *github.Response, error)
	CreateHook(owner string, repo string, hook *github.Hook) (*github.Hook, *github.Response, error)
	DownloadContents(owner string, repo string, ref string, path string) (file io.ReadCloser, exists bool, err error)
}

type HTTPRepositoryAPI struct {
	ctx context.Context
	c   *github.RepositoriesService
}

// GitHub models

type HookCfg struct {
	URL         string `json:"url"`
	ContentType string `json:"content_type"`
}

type RepoOpts struct {
	Owner string `json:"owner"`
	Repo  string `json:"repo"`
	Ref   string `json:"ref"`
}
