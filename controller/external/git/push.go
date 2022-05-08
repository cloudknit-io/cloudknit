package git

import (
	"github.com/pkg/errors"

	gogit "github.com/go-git/go-git/v5"
)

func (g *GoGit) Push() error {
	if g.r == nil {
		return errors.Wrapf(ErrRepoNotCloned, "cannot push")
	}

	auth, err := g.getAuthOptions()
	if err != nil {
		return errors.Wrap(err, "error getting auth options")
	}

	return g.r.Push(&gogit.PushOptions{
		Auth: auth,
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
