package github

import (
	"context"
	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"
	"io"
)

func NewHttpRepositoryClient(token string, ctx context.Context) RepositoryApi {
	client := createGithubClient(token, ctx).Repositories
	return HttpRepositoryApi{Client: client, Ctx: ctx}
}

func NewHttpGitClient(token string, ctx context.Context) GitApi {
	client := createGithubClient(token, ctx).Git
	return HttpGitApi{Client: client, Ctx: ctx}
}

func createGithubClient(token string, ctx context.Context) *github.Client {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}

func (api HttpGitApi) GetRef(owner string, repo string, ref string) (*github.Reference, *github.Response, error) {
	return api.Client.GetRef(api.Ctx, owner, repo, ref)
}

func (api HttpGitApi) UpdateRef(owner string, repo string, ref *github.Reference, force bool) (*github.Reference, *github.Response, error) {
	return api.Client.UpdateRef(api.Ctx, owner, repo, ref, force)
}

func (api HttpGitApi) GetCommit(owner string, repo string, sha string) (*github.Commit, *github.Response, error) {
	return api.Client.GetCommit(api.Ctx, owner, repo, sha)
}

func (api HttpGitApi) CreateCommit(owner string, repo string, commit *github.Commit) (*github.Commit, *github.Response, error) {
	return api.Client.CreateCommit(api.Ctx, owner, repo, commit)
}

func (api HttpGitApi) GetTree(owner string, repo string, sha string, recursive bool) (*github.Tree, *github.Response, error) {
	return api.Client.GetTree(api.Ctx, owner, repo, sha, recursive)
}

func (api HttpGitApi) CreateTree(owner string, repo string, baseTree string, entries []*github.TreeEntry) (*github.Tree, *github.Response, error) {
	return api.Client.CreateTree(api.Ctx, owner, repo, baseTree, entries)
}

func (api HttpRepositoryApi) CreateRepository(owner string, repo string) (*github.Repository, *github.Response, error) {
	r := github.Repository{Name: github.String(repo), Private: github.Bool(true)}
	return api.Client.Create(api.Ctx, owner, &r)
}

func (api HttpRepositoryApi) GetRepository(owner string, repo string) (*github.Repository, *github.Response, error) {
	return api.Client.Get(api.Ctx, owner, repo)
}

func (api HttpRepositoryApi) ListHooks(owner string, repo string, opts *github.ListOptions) ([]*github.Hook, *github.Response, error) {
	return api.Client.ListHooks(api.Ctx, owner, repo, opts)
}

func (api HttpRepositoryApi) CreateHook(owner string, repo string, hook *github.Hook) (*github.Hook, *github.Response, error) {
	return api.Client.CreateHook(api.Ctx, owner, repo, hook)
}

func (api HttpRepositoryApi) DownloadContents(owner string, repo string, ref string, path string) (io.ReadCloser, error) {
	opts := &github.RepositoryContentGetOptions{Ref: ref}
	return api.Client.DownloadContents(api.Ctx, owner, repo, path, opts)
}
