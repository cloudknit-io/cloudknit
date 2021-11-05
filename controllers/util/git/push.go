package git

import (
	"errors"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

func (g *GoGit) Push() error {
	if g.r == nil {
		return errors.New("repo not cloned")
	}
	return g.r.Push(&gogit.PushOptions{
		Auth: &http.BasicAuth{
			Username: "zlifecycle",
			Password: g.token,
		},
	})
}

func (g *GoGit) CommitAndPush(nfo *CommitInfo) error {
	if g.r == nil {
		return errors.New("repo not cloned")
	}
	if _, err := g.Commit(nfo); err != nil {
		return err
	}

	return g.Push()
}
