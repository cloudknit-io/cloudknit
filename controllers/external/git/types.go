package git

import (
	"context"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/pkg/errors"
)

var (
	ErrEmptyCommit   = errors.New("commit is empty")
	ErrRepoNotCloned = errors.New("repo not cloned")
)

type API interface {
	Clone(repo string, dir string) error
	Open(path string) error
	Commit(nfo *CommitInfo) (*object.Commit, error)
	HeadCommitHash() (hash string, err error)
	CommitAndPush(nfo *CommitInfo) (pushed bool, err error)
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

func NewGoGit(ctx context.Context, token string) (API, error) {
	return &GoGit{token: token, ctx: ctx}, nil
}

var _ API = (*GoGit)(nil)
