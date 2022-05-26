package git

import (
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/pkg/errors"
)

var (
	ErrEmptyCommit   = errors.New("commit is empty")
	ErrRepoNotCloned = errors.New("repo not cloned")
)

//go:generate mockgen --build_flags=--mod=mod -destination=./mock_git_api.go -package=git "github.com/compuzest/zlifecycle-il-operator/controller/common/git" API
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
