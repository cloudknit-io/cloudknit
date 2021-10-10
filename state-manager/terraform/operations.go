package terraform

import (
	"context"
	"github.com/compuzest/zlifecycle-state-manager/zlog"
	"github.com/sirupsen/logrus"
)

func GetState(ctx context.Context, workdir string) (*StateWrapper, error) {
	execPath, err := GetExecPath(ctx)
	if err != nil {
		return nil, err
	}

	w, err := InitTerraform(ctx, workdir, execPath)
	if err != nil {
		return nil, err
	}

	zlog.Logger.WithFields(
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
		return nil, err
	}

	w, err := InitTerraform(ctx, workdir, execPath)
	if err != nil {
		return nil, err
	}

	zlog.Logger.WithFields(
		logrus.Fields{"workdir": workdir, "resources": resources},
	).Info("Removing resources from terraform state")
	state, err := w.RemoveStateResources(resources)
	if err != nil {
		return nil, err
	}

	return state, nil
}
