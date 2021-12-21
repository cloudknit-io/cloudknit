package terraform

import (
	"context"

	"github.com/compuzest/zlifecycle-state-manager/app/zlog"
	"github.com/hashicorp/terraform-exec/tfinstall"
	"github.com/pkg/errors"
)

var globalExecPath = ""

// findExecutable searches the file system for the terraform binary.
func findExecutable(ctx context.Context) (execPath string, err error) {
	execPath, err = tfinstall.Find(ctx, tfinstall.LookPath())
	return
}

// GetExecPath should always be used to get the terraform executable path as it reuses the global singleton variable
func GetExecPath(ctx context.Context) (execPath string, err error) {
	if globalExecPath == "" {
		zlog.CtxLogger(ctx).Info("Searching file system for terraform executable")
		globalExecPath, err = findExecutable(ctx)
		if err != nil {
			return "", errors.Wrap(err, "error looking up terraform executable path")
		}
	}
	return globalExecPath, nil
}
