package git

import (
	gogit "github.com/go-git/go-git/v5"
	"github.com/pkg/errors"
)

func (g *GoGit) Pull() (updated bool, err error) {
	if g.r == nil {
		return false, errors.Wrapf(ErrRepoNotCloned, "cannot pull")
	}

	auth, err := g.getAuthOptions()
	if err != nil {
		return false, errors.Wrap(err, "error getting auth options")
	}

	w, err := g.r.Worktree()
	if err != nil {
		return false, err
	}
	oldHead, err := g.r.Head()
	if err != nil {
		return false, err
	}

	if err := w.Pull(
		&gogit.PullOptions{
			RemoteName:   "origin",
			SingleBranch: true,
			Auth:         auth,
		},
	); err != nil {
		return false, err
	}

	newHead, err := g.r.Head()
	if err != nil {
		return false, err
	}

	// if HEAD refs differ, then we pulled changes for the repo
	return oldHead != newHead, nil
}
