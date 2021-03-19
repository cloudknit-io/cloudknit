package github

import (
	"context"
	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"
)

func NewHttpClient(token string, ctx context.Context) RepositoryApi {
	return HttpRepositoryApi{Client: createGithubClient(token, ctx), Ctx: ctx}
}

func createGithubClient(token string, ctx context.Context) *github.Client {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}

func (api HttpRepositoryApi) CreateRepository(owner string, repo string) (*github.Repository, *github.Response, error) {
	r := github.Repository{Name: github.String(repo), Private: github.Bool(true)}
	return api.Client.Repositories.Create(api.Ctx, owner, &r)
}

func (api HttpRepositoryApi) GetRepository(owner string, repo string) (*github.Repository, *github.Response, error) {
	return api.Client.Repositories.Get(api.Ctx, owner, repo)
}

func (api HttpRepositoryApi) ListHooks(owner string, repo string, opts *github.ListOptions) ([]*github.Hook, *github.Response, error) {
	return api.Client.Repositories.ListHooks(api.Ctx, owner, repo, opts)
}

func (api HttpRepositoryApi) CreateHook(owner string, repo string, hook *github.Hook) (*github.Hook, *github.Response, error) {
	return api.Client.Repositories.CreateHook(api.Ctx, owner, repo, hook)
}
