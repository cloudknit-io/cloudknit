package github

import (
	"context"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/bradleyfalzon/ghinstallation"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/common"
	"github.com/pkg/errors"

	"github.com/compuzest/zlifecycle-il-operator/controllers/util/env"

	"github.com/google/go-github/v42/github"
	"golang.org/x/oauth2"
)

func NewHTTPRepositoryAPI(ctx context.Context) RepositoryAPI {
	client := newGitHubClient(ctx).Repositories
	return &HTTPRepositoryAPI{c: client, ctx: ctx}
}

func NewHTTPGitClient(ctx context.Context) GitAPI {
	client := newGitHubClient(ctx).Git
	return &HTTPGitAPI{c: client, ctx: ctx}
}

func NewHTTPAppClient(ctx context.Context, privateKey []byte, appID string) (AppAPI, error) {
	int64AppID, err := strconv.ParseInt(appID, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "error parsing app id as int64")
	}
	client, err := newGitHubAppClient(privateKey, int64AppID)
	if err != nil {
		return nil, errors.Wrap(err, "error creating github app client")
	}
	return &HTTPAppAPI{ctx: ctx, c: client.Apps}, nil
}

func newGitHubAppClient(privateKey []byte, appID int64) (*github.Client, error) {
	tr := http.DefaultTransport

	// Wrap the shared transport for use with the app ID 1 authenticating with installation ID 99.
	itr, err := ghinstallation.NewAppsTransport(tr, appID, privateKey)
	if err != nil {
		return nil, errors.Wrap(err, "error creating github app transport")
	}

	// Use installation transport with github.com/google/go-github
	client := github.NewClient(&http.Client{Transport: itr})

	return client, nil
}

func newGitHubClient(ctx context.Context) *github.Client {
	token := env.Config.GitToken
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}

// Git API

func (api *HTTPGitAPI) GetRef(owner string, repo string, ref string) (*github.Reference, *github.Response, error) {
	return api.c.GetRef(api.ctx, owner, repo, ref)
}

func (api *HTTPGitAPI) UpdateRef(owner string, repo string, ref *github.Reference, force bool) (*github.Reference, *github.Response, error) {
	return api.c.UpdateRef(api.ctx, owner, repo, ref, force)
}

func (api *HTTPGitAPI) GetCommit(owner string, repo string, sha string) (*github.Commit, *github.Response, error) {
	return api.c.GetCommit(api.ctx, owner, repo, sha)
}

func (api *HTTPGitAPI) CreateCommit(owner string, repo string, commit *github.Commit) (*github.Commit, *github.Response, error) {
	return api.c.CreateCommit(api.ctx, owner, repo, commit)
}

func (api *HTTPGitAPI) GetTree(owner string, repo string, sha string, recursive bool) (*github.Tree, *github.Response, error) {
	return api.c.GetTree(api.ctx, owner, repo, sha, recursive)
}

func (api *HTTPGitAPI) CreateTree(owner string, repo string, baseTree string, entries []*github.TreeEntry) (*github.Tree, *github.Response, error) {
	return api.c.CreateTree(api.ctx, owner, repo, baseTree, entries)
}

// Repository API

func (api *HTTPRepositoryAPI) CreateRepository(owner string, repo string) (*github.Repository, *github.Response, error) {
	r := github.Repository{Name: github.String(repo), Private: github.Bool(true)}
	return api.c.Create(api.ctx, owner, &r)
}

func (api *HTTPRepositoryAPI) GetRepository(owner string, repo string) (*github.Repository, *github.Response, error) {
	return api.c.Get(api.ctx, owner, repo)
}

func (api *HTTPRepositoryAPI) ListHooks(owner string, repo string, opts *github.ListOptions) ([]*github.Hook, *github.Response, error) {
	return api.c.ListHooks(api.ctx, owner, repo, opts)
}

func (api *HTTPRepositoryAPI) CreateHook(owner string, repo string, hook *github.Hook) (*github.Hook, *github.Response, error) {
	return api.c.CreateHook(api.ctx, owner, repo, hook)
}

func (api *HTTPRepositoryAPI) DownloadContents(owner string, repo string, ref string, path string) (file io.ReadCloser, exists bool, err error) {
	opts := &github.RepositoryContentGetOptions{Ref: ref}
	f, resp, err := api.c.DownloadContents(api.ctx, owner, repo, path, opts)
	if err != nil {
		if strings.HasPrefix(err.Error(), "No file named") {
			return nil, false, nil
		}
		return nil, false, err
	}
	defer common.CloseBody(resp.Body)
	return f, true, nil
}

// GitHub App API

func (api *HTTPAppAPI) FindRepositoryInstallation(owner string, repo string) (*github.Installation, *github.Response, error) {
	return api.c.FindRepositoryInstallation(api.ctx, owner, repo)
}

var (
	_ RepositoryAPI = (*HTTPRepositoryAPI)(nil)
	_ GitAPI        = (*HTTPGitAPI)(nil)
	_ AppAPI        = (*HTTPAppAPI)(nil)
)
