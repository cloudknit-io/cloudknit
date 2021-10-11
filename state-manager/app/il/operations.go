package il

import (
	"context"
	"github.com/compuzest/zlifecycle-state-manager/app/git"
	"github.com/compuzest/zlifecycle-state-manager/app/terraform"
	"github.com/compuzest/zlifecycle-state-manager/app/util"
	"github.com/compuzest/zlifecycle-state-manager/app/zlog"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"time"
)

func FetchState(ctx context.Context, zs *ZState) (*terraform.StateWrapper, error) {
	start := time.Now()

	workdir, err := getTerraformWorkdir(zs)
	if err != nil {
		return nil, err
	}

	state, err := terraform.GetState(ctx, workdir)
	if err != nil {
		return nil, err
	}

	zlog.Logger.WithFields(
		logrus.Fields{
			"il": zs.Meta.IL,
			"team": zs.Meta.Team,
			"environment": zs.Meta.Environment,
			"component": zs.Meta.Component,
			"duration": time.Since(start),
		},
	).Info("Successfully fetched terraform state")

	return state, nil
}

func RemoveStateResources(ctx context.Context, zs *ZState, resources []string) (*terraform.StateWrapper, error) {
	start := time.Now()

	workdir, err := getTerraformWorkdir(zs)
	if err != nil {
		return nil, err
	}

	state, err := terraform.RemoveResources(ctx, workdir, resources)
	if err != nil {
		return nil, err
	}

	zlog.Logger.WithFields(
		logrus.Fields{
			"il": zs.Meta.IL,
			"team": zs.Meta.Team,
			"environment": zs.Meta.Environment,
			"component": zs.Meta.Component,
			"resources": resources,
			"duration": time.Since(start),
		},
	).Info("Successfully removed resources from terraform state")

	return state, nil
}

func getTerraformWorkdir(zs *ZState) (string, error) {
	basedir := path.Join("/tmp", zs.Meta.IL)
	exists, err := util.DirExists(basedir)
	if err != nil {
		return "", err
	}

	if !exists {
		err = os.MkdirAll(basedir, 0755)
		if err != nil {
			return "", err
		}
	}

	_, _, err = git.GetRepository(zs.RepoURL, basedir)
	if err != nil {
		return "", err
	}

	workdir, err := BuildILComponentPath(zs.Meta, basedir)
	if err != nil {
		return "", err
	}

	return workdir, nil
}
