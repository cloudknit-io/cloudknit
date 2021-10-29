package github

import (
	"context"
	"io"
	"strings"

	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"
)

func NewHTTPRepositoryAPI(ctx context.Context, token string) RepositoryAPI {
	client := createGithubClient(ctx, token).Repositories
	return HTTPRepositoryAPI{Client: client, Ctx: ctx}
}

func NewHTTPGitClient(ctx context.Context, token string) GitAPI {
	client := createGithubClient(ctx, token).Git
	return HTTPGitAPI{Client: client, Ctx: ctx}
}

func createGithubClient(ctx context.Context, token string) *github.Client {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}

func (api HTTPGitAPI) GetRef(owner string, repo string, ref string) (*github.Reference, *github.Response, error) {
	return api.Client.GetRef(api.Ctx, owner, repo, ref)
}

func (api HTTPGitAPI) UpdateRef(owner string, repo string, ref *github.Reference, force bool) (*github.Reference, *github.Response, error) {
	return api.Client.UpdateRef(api.Ctx, owner, repo, ref, force)
}

func (api HTTPGitAPI) GetCommit(owner string, repo string, sha string) (*github.Commit, *github.Response, error) {
	return api.Client.GetCommit(api.Ctx, owner, repo, sha)
}

func (api HTTPGitAPI) CreateCommit(owner string, repo string, commit *github.Commit) (*github.Commit, *github.Response, error) {
	return api.Client.CreateCommit(api.Ctx, owner, repo, commit)
}

func (api HTTPGitAPI) GetTree(owner string, repo string, sha string, recursive bool) (*github.Tree, *github.Response, error) {
	return api.Client.GetTree(api.Ctx, owner, repo, sha, recursive)
}

func (api HTTPGitAPI) CreateTree(owner string, repo string, baseTree string, entries []*github.TreeEntry) (*github.Tree, *github.Response, error) {
	return api.Client.CreateTree(api.Ctx, owner, repo, baseTree, entries)
}

func (api HTTPRepositoryAPI) CreateRepository(owner string, repo string) (*github.Repository, *github.Response, error) {
	r := github.Repository{Name: github.String(repo), Private: github.Bool(true)}
	return api.Client.Create(api.Ctx, owner, &r)
}

func (api HTTPRepositoryAPI) GetRepository(owner string, repo string) (*github.Repository, *github.Response, error) {
	return api.Client.Get(api.Ctx, owner, repo)
}

func (api HTTPRepositoryAPI) ListHooks(owner string, repo string, opts *github.ListOptions) ([]*github.Hook, *github.Response, error) {
	return api.Client.ListHooks(api.Ctx, owner, repo, opts)
}

func (api HTTPRepositoryAPI) CreateHook(owner string, repo string, hook *github.Hook) (*github.Hook, *github.Response, error) {
	return api.Client.CreateHook(api.Ctx, owner, repo, hook)
}

func (api HTTPRepositoryAPI) DownloadContents(owner string, repo string, ref string, path string) (file io.ReadCloser, exists bool, err error) {
	opts := &github.RepositoryContentGetOptions{Ref: ref}
	f, err := api.Client.DownloadContents(api.Ctx, owner, repo, path, opts)
	if err != nil {
		if strings.HasPrefix(err.Error(), "No file named") {
			return nil, false, nil
		}
		return nil, false, err
	}
	return f, true, nil
}
