package git

import (
	"context"
	"errors"

	"github.com/compuzest/zlifecycle-state-manager/app/env"
	"github.com/compuzest/zlifecycle-state-manager/app/util"
	"github.com/compuzest/zlifecycle-state-manager/app/zlog"
	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/sirupsen/logrus"
)

func auth() *http.BasicAuth {
	token := env.Config().GitToken
	return &http.BasicAuth{
		Username: "zlifecycle", // yes, this can be anything except an empty string
		Password: token,
	}
}

func openRepo(workdir string) (repo *gogit.Repository, err error) {
	repo, err = gogit.PlainOpen(workdir)
	return
}

func cloneRepoFS(url string, workdir string) (*gogit.Repository, error) {
	repo, err := gogit.PlainClone(workdir, false, &gogit.CloneOptions{
		Auth: auth(),
		URL:  url,
	})
	if err != nil {
		return nil, err
	}
	return repo, nil
}

func pullRepo(repo *gogit.Repository) (dirty bool, err error) {
	w, err := repo.Worktree()
	if err != nil {
		return false, err
	}

	err = w.Pull(&gogit.PullOptions{RemoteName: "origin", Auth: auth()})

	if errors.Is(err, gogit.NoErrAlreadyUpToDate) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, err
}

func GetRepository(ctx context.Context, url string, workdir string) (repo *gogit.Repository, dirty bool, err error) {
	var r *gogit.Repository
	dirty = false

	exists, err := util.DirExists(workdir)
	if err != nil {
		return nil, false, err
	}
	empty, err := util.IsDirEmpty(workdir)
	if err != nil {
		return nil, false, err
	}
	if exists && !empty {
		zlog.CtxLogger(ctx).WithFields(
			logrus.Fields{"url": url, "workdir": workdir},
		).Info("Opening existing repo from filesystem")
		r, err = openRepo(workdir)
		if err != nil {
			return nil, false, err
		}

		zlog.CtxLogger(ctx).WithFields(
			logrus.Fields{"url": url, "workdir": workdir},
		).Info("Pulling git changes")
		dirty, err = pullRepo(r)
		if err != nil {
			return nil, false, err
		}
	} else {
		zlog.CtxLogger(ctx).WithFields(
			logrus.Fields{"url": url, "workdir": workdir},
		).Info("Cloning repo")
		r, err = cloneRepoFS(url, workdir)
		if err != nil {
			return nil, false, err
		}
	}

	return r, dirty, nil
}
