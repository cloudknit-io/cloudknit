package git

import (
	"errors"
	"time"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func (g *GoGit) Commit(nfo *CommitInfo) (*object.Commit, error) {
	if g.r == nil {
		return nil, errors.New("repo not cloned")
	}
	w, err := g.r.Worktree()
	if err != nil {
		return nil, err
	}

	commitHash, err := w.Commit(nfo.Msg, &gogit.CommitOptions{
		All: true,
		Author: &object.Signature{
			Name:  nfo.Author,
			Email: nfo.Email,
			When:  time.Now(),
		},
	})
	if err != nil {
		return nil, err
	}

	commit, err := g.r.CommitObject(commitHash)
	if err != nil {
		return nil, err
	}

	return commit, nil
}
