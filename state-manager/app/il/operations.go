package il

import (
	"context"
	"os"
	"path"
	"time"

	"github.com/compuzest/zlifecycle-state-manager/app/git"
	"github.com/compuzest/zlifecycle-state-manager/app/terraform"
	"github.com/compuzest/zlifecycle-state-manager/app/util"
	"github.com/compuzest/zlifecycle-state-manager/app/zlog"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func FetchState(ctx context.Context, zs *ZState) (*terraform.StateWrapper, error) {
	start := time.Now()

	workdir, err := getTerraformWorkdir(ctx, zs)
	if err != nil {
		return nil, errors.Wrap(err, "error getting terraform workdir")
	}

	state, err := terraform.GetState(ctx, workdir)
	if err != nil {
		return nil, errors.Wrap(err, "error executing terraform state ls operation")
	}

	zlog.CtxLogger(ctx).WithFields(
		logrus.Fields{
			"il":          zs.Meta.IL,
			"team":        zs.Meta.Team,
			"environment": zs.Meta.Environment,
			"component":   zs.Meta.Component,
			"duration":    time.Since(start),
		},
	).Info("Successfully fetched terraform state")

	return state, nil
}

func RemoveStateResources(ctx context.Context, zs *ZState, resources []string) (*terraform.StateWrapper, error) {
	start := time.Now()

	workdir, err := getTerraformWorkdir(ctx, zs)
	if err != nil {
		return nil, errors.Wrap(err, "error getting terraform workdir")
	}

	state, err := terraform.RemoveResources(ctx, workdir, resources)
	if err != nil {
		return nil, err
	}

	zlog.CtxLogger(ctx).WithFields(
		logrus.Fields{
			"il":          zs.Meta.IL,
			"team":        zs.Meta.Team,
			"environment": zs.Meta.Environment,
			"component":   zs.Meta.Component,
			"resources":   resources,
			"duration":    time.Since(start),
		},
	).Info("Successfully removed resources from terraform state")

	return state, nil
}

func getTerraformWorkdir(ctx context.Context, zs *ZState) (string, error) {
	basedir := path.Join("/tmp", zs.Meta.IL)
	exists, err := util.DirExists(basedir)
	if err != nil {
		return "", err
	}

	if !exists {
		err = os.MkdirAll(basedir, 0o755)
		if err != nil {
			return "", err
		}
	}

	_, _, err = git.GetRepository(ctx, zs.RepoURL, basedir)
	if err != nil {
		return "", errors.Wrapf(err, "error getting repository: %s", zs.RepoURL)
	}

	workdir, err := BuildILComponentPath(zs.Meta, basedir)
	if err != nil {
		return "", err
	}

	return workdir, nil
}
