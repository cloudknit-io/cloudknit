package github

import (
	"context"
	"io"

	_ "github.com/golang/mock/mockgen/model" // workaround for mockgen failing
	"github.com/google/go-github/v32/github"
)

//go:generate go run --mod=mod github.com/golang/mock/mockgen -destination=../../../mocks/mock_github_api.go -package=mocks "github.com/compuzest/zlifecycle-il-operator/controllers/util/github" GitAPI,RepositoryAPI

type GitAPI interface {
	GetRef(owner string, repo string, ref string) (*github.Reference, *github.Response, error)
	UpdateRef(owner string, repo string, ref *github.Reference, force bool) (*github.Reference, *github.Response, error)
	GetCommit(owner string, repo string, sha string) (*github.Commit, *github.Response, error)
	CreateCommit(owner string, repo string, commit *github.Commit) (*github.Commit, *github.Response, error)
	GetTree(owner string, repo string, sha string, recursive bool) (*github.Tree, *github.Response, error)
	CreateTree(owner string, repo string, baseTree string, entries []*github.TreeEntry) (*github.Tree, *github.Response, error)
}

type HTTPGitAPI struct {
	Ctx    context.Context
	Client *github.GitService
}

type RepositoryAPI interface {
	CreateRepository(owner string, repo string) (*github.Repository, *github.Response, error)
	GetRepository(owner string, repo string) (*github.Repository, *github.Response, error)
	ListHooks(owner string, repo string, opts *github.ListOptions) ([]*github.Hook, *github.Response, error)
	CreateHook(owner string, repo string, hook *github.Hook) (*github.Hook, *github.Response, error)
	DownloadContents(owner string, repo string, ref string, path string) (io.ReadCloser, error)
}

type HTTPRepositoryAPI struct {
	Ctx    context.Context
	Client *github.RepositoriesService
}

type Package struct {
	FullName      string
	Description   string
	StarsCount    int
	ForksCount    int
	LastUpdatedBy string
}

type HookCfg struct {
	URL         string `json:"url"`
	ContentType string `json:"content_type"`
}

type RepoOpts struct {
	Owner string `json:"owner"`
	Repo  string `json:"repo"`
	Ref   string `json:"ref"`
}

type Owner = string

type Repo = string
