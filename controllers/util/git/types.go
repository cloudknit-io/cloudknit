package git

import (
	"errors"
	"os"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type Git interface {
	Clone(repo string, dir string) error
	Commit(nfo *CommitInfo) (*object.Commit, error)
	CommitAndPush(nfo *CommitInfo) error
	Push() error
}

type GoGit struct {
	r     *gogit.Repository
	token string
}

type CommitInfo struct {
	Author string
	Email  string
	Msg    string
}

func NewGoGit() (*GoGit, error) {
	t := os.Getenv("GITHUB_AUTH_TOKEN")
	if t == "" {
		return nil, errors.New("missing github auth token")
	}

	return &GoGit{token: t}, nil
}

var _ Git = (*GoGit)(nil)
