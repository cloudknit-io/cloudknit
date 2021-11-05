package git

import (
	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

func (g *GoGit) Clone(repo string, directory string) error {
	r, err := gogit.PlainClone(directory, false, &gogit.CloneOptions{
		URL: repo,
		Auth: &http.BasicAuth{
			Username: "zlifecycle",
			Password: g.token,
		},
	})
	if err != nil {
		return err
	}
	g.r = r
	return nil
}
