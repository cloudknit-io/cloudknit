package github

import (
	"io"

	"github.com/google/go-github/v42/github"
)

//go:generate mockgen --build_flags=--mod=mod -destination=./mock_github_api.go -package=github "github.com/compuzest/zlifecycle-il-operator/controller/common/github" API
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
