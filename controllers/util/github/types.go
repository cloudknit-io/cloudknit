package github

import (
	"context"
	"github.com/google/go-github/v32/github"
)

type RepositoryApi interface {
	ListHooks(owner string, repo string, opts *github.ListOptions) ([]*github.Hook, *github.Response, error)
	CreateHook(owner string, repo string, hook *github.Hook) (*github.Hook, *github.Response, error)
}

type HttpRepositoryApi struct {
	Ctx    context.Context
	Client *github.Client
}

type Package struct {
	FullName      string
	Description   string
	StarsCount    int
	ForksCount    int
	LastUpdatedBy string
}

type HookCfg struct {
	Url         string `json:"url"`
	ContentType string `json:"content_type"`
}

type Owner = string

type Repo = string
