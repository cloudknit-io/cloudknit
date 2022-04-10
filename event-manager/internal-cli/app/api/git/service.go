package git

import (
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

type CommitInfo struct {
	Author string
	Email  string
	Msg    string
}
