package terraform

import (
	"context"
	"github.com/compuzest/zlifecycle-state-manager/app/zlog"
	"github.com/hashicorp/terraform-exec/tfinstall"
)

var (
	globalExecPath = ""
)

// FindExecutable searches the file system for the terraform binary.
func FindExecutable(ctx context.Context) (execPath string, err error) {
	execPath, err = tfinstall.Find(ctx, tfinstall.LookPath())
	return
}

// GetExecPath should always be used to get the terraform executable path as it reuses the global singleton variable
func GetExecPath(ctx context.Context) (execPath string, err error) {
	if globalExecPath == "" {
		zlog.Logger.Info("Searching file system for terraform executable")
		globalExecPath, err = FindExecutable(ctx)
		if err != nil {
			return "", err
		}
	}
	return globalExecPath, nil
}
