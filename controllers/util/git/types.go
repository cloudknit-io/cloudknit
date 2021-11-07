package git

import (
	"context"
	"errors"
	"os"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

var (
	ErrEmptyCommit   = errors.New("commit is empty")
	ErrRepoNotCloned = errors.New("repo not cloned")
)

type Git interface {
	Clone(repo string, dir string) error
	Open(path string) error
	Commit(nfo *CommitInfo) (*object.Commit, error)
	CommitAndPush(nfo *CommitInfo) (empty bool, err error)
	Push() error
}

type GoGit struct {
	ctx   context.Context
	r     *gogit.Repository
	token string
}

type CommitInfo struct {
	Author string
	Email  string
	Msg    string
}

func NewGoGit(ctx context.Context) (*GoGit, error) {
	t := os.Getenv("GITHUB_AUTH_TOKEN")
	if t == "" {
		return nil, errors.New("missing github auth token")
	}

	return &GoGit{token: t, ctx: ctx}, nil
}

var _ Git = (*GoGit)(nil)
