package github

import (
	"context"
	"io"
	"net/http"
	"strings"

	"github.com/compuzest/zlifecycle-il-operator/controller/util"

	"github.com/bradleyfalzon/ghinstallation"
	"github.com/pkg/errors"

	"github.com/google/go-github/v42/github"
	"golang.org/x/oauth2"
)

const (
	modeGitHubApp = "githubApp"
	modeToken     = "token"
)

func newGitHubAppClient(privateKey []byte, appID int64) (*github.Client, error) {
	tr := http.DefaultTransport

	itr, err := ghinstallation.NewAppsTransport(tr, appID, privateKey)
	if err != nil {
		return nil, errors.Wrap(err, "error creating github app transport")
	}

	// Use installation transport with github.com/google/go-github
	client := github.NewClient(&http.Client{Transport: itr})

	return client, nil
}

func newGitHubAppClientWithInstallation(privateKey []byte, appID, installationID int64) (*github.Client, error) {
	tr := http.DefaultTransport

	itr, err := ghinstallation.New(tr, appID, installationID, privateKey)
	if err != nil {
		return nil, errors.Wrap(err, "error creating github app transport")
	}

	// Use installation transport with github.com/google/go-github
	client := github.NewClient(&http.Client{Transport: itr})

	return client, nil
}

func newTokenGitHubClient(ctx context.Context, token string) *github.Client {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
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
	defer util.CloseBody(resp.Body)
	return f, true, nil
}

// GitHub Client

func (c *Client) FindRepositoryInstallation(owner string, repo string) (*github.Installation, *github.Response, error) {
	return c.appClient.Apps.FindRepositoryInstallation(c.ctx, owner, repo)
}

func (c *Client) FindOrganizationInstallation(org string) (*github.Installation, *github.Response, error) {
	return c.appClient.Apps.FindOrganizationInstallation(c.ctx, org)
}

func (c *Client) CreateInstallationToken(installationID int64) (*github.InstallationToken, *github.Response, error) {
	return c.appClient.Apps.CreateInstallationToken(c.ctx, installationID, nil)
}

func (c *Client) CreateRepository(owner string, repo string) (*github.Repository, *github.Response, error) {
	r := github.Repository{Name: github.String(repo), Private: github.Bool(true)}
	return c.client.Repositories.Create(c.ctx, owner, &r)
}

func (c *Client) GetRepository(owner string, repo string) (*github.Repository, *github.Response, error) {
	return c.client.Repositories.Get(c.ctx, owner, repo)
}

func (c *Client) ListHooks(owner string, repo string, opts *github.ListOptions) ([]*github.Hook, *github.Response, error) {
	return c.client.Repositories.ListHooks(c.ctx, owner, repo, opts)
}

func (c *Client) CreateHook(owner string, repo string, hook *github.Hook) (*github.Hook, *github.Response, error) {
	return c.client.Repositories.CreateHook(c.ctx, owner, repo, hook)
}

func (c *Client) DownloadContents(owner string, repo string, ref string, path string) (file io.ReadCloser, exists bool, err error) {
	opts := &github.RepositoryContentGetOptions{Ref: ref}
	f, resp, err := c.client.Repositories.DownloadContents(c.ctx, owner, repo, path, opts)
	if err != nil {
		if strings.HasPrefix(err.Error(), "No file named") {
			return nil, false, nil
		}
		return nil, false, err
	}
	defer util.CloseBody(resp.Body)
	return f, true, nil
}

var (
	_ RepositoryAPI = (*HTTPRepositoryAPI)(nil)
	_ API           = (*Client)(nil)
)
