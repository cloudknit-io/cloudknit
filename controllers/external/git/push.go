package git

import (
	"github.com/pkg/errors"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

func (g *GoGit) Push() error {
	if g.r == nil {
		return errors.Wrapf(ErrRepoNotCloned, "cannot push")
	}

	return g.r.Push(&gogit.PushOptions{
		Auth: &http.BasicAuth{
			Username: "zlifecycle",
			Password: g.token,
		},
	})
}

func (g *GoGit) CommitAndPush(nfo *CommitInfo) (pushed bool, err error) {
	if g.r == nil {
		return false, errors.Wrapf(ErrRepoNotCloned, "cannot commit and push")
	}

	if _, err := g.Commit(nfo); err != nil {
		if errors.Is(err, ErrEmptyCommit) {
			return false, nil
		}
		return false, err
	}

	return true, g.Push()
}
