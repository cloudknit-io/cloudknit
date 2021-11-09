package git

import (
	"time"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// Commit will add all files changed files to the staging area and commit them.
// Throws an ErrEmptyCommit if no files are changed.
// It returns the commit object.
func (g *GoGit) Commit(nfo *CommitInfo) (commit *object.Commit, err error) {
	if g.r == nil {
		return nil, ErrRepoNotCloned
	}
	w, err := g.r.Worktree()
	if err != nil {
		return nil, err
	}

	_, err = w.Add(".")
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

	commit, err = g.r.CommitObject(commitHash)
	if err != nil {
		return nil, err
	}

	fileStats, err := commit.StatsContext(g.ctx)
	if err != nil {
		return nil, err
	}

	if len(fileStats) == 0 {
		return nil, ErrEmptyCommit
	}

	return commit, nil
}

func (g *GoGit) HeadCommitHash() (hash string, err error) {
	ref, err := g.r.Head()
	if err != nil {
		return "", err
	}

	return ref.Hash().String(), nil
}
