package terraform

import (
	"context"

	"github.com/compuzest/zlifecycle-state-manager/app/zlog"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func GetState(ctx context.Context, workdir string) (*StateWrapper, error) {
	execPath, err := GetExecPath(ctx)
	if err != nil {
		return nil, err
	}

	w, err := InitTerraform(ctx, workdir, execPath)
	if err != nil {
		return nil, errors.Wrap(err, "error getting terraform state")
	}

	zlog.CtxLogger(ctx).WithFields(
		logrus.Fields{"workdir": workdir},
	).Info("Fetching terraform state")
	state, err := w.State()
	if err != nil {
		return nil, err
	}

	return state, nil
}

func RemoveResources(ctx context.Context, workdir string, resources []string) (*StateWrapper, error) {
	execPath, err := GetExecPath(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "error get terraform executable path")
	}

	w, err := InitTerraform(ctx, workdir, execPath)
	if err != nil {
		return nil, errors.Wrap(err, "error initializing terraform")
	}

	zlog.CtxLogger(ctx).WithFields(
		logrus.Fields{"workdir": workdir, "resources": resources},
	).Info("Removing resources from terraform state")
	state, err := w.RemoveStateResources(resources)
	if err != nil {
		return nil, errors.Wrap(err, "error removing state resources")
	}

	return state, nil
}
