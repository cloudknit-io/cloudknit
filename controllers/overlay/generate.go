package overlay

import (
	"context"
	"path/filepath"

	"github.com/compuzest/zlifecycle-il-operator/controllers/util/common"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/git"

	stablev1 "github.com/compuzest/zlifecycle-il-operator/api/v1"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/file"
	"github.com/go-logr/logr"
)

func GenerateOverlayFiles(
	ctx context.Context,
	log logr.Logger,
	fileAPI file.Service,
	e *stablev1.Environment,
	ec *stablev1.EnvironmentComponent,
	destinationFolder string,
) error {
	gitAPI, err := git.NewGoGit(ctx)
	if err != nil {
		return err
	}
	// generate overlay files from config repo
	for _, overlay := range ec.OverlayFiles {
		tempDir, cleanup, err := git.CloneTemp(gitAPI, overlay.Source)
		defer cleanup()
		if err != nil {
			return err
		}
		for _, path := range overlay.Paths {
			log.Info(
				"Generating overlay file",
				"repo", overlay.Source,
				"ref", overlay.Ref,
				"source", path,
				"destination", destinationFolder,
				"team", e.Spec.TeamName,
				"environment", e.Spec.EnvName,
				"component", ec.Name,
			)
			absolutePath := filepath.Join(tempDir, path)
			if common.IsDir(absolutePath) {
				if err := common.CopyDirContent(absolutePath, destinationFolder); err != nil {
					return err
				}
			} else {
				if err := common.CopyFile(absolutePath, destinationFolder); err != nil {
					return err
				}
			}
		}
	}
	// Generate inline overlay files
	for _, overlay := range ec.OverlayData {
		log.Info(
			"Generating overlay file from data field",
			"overlay", overlay.Name,
			"destination", destinationFolder,
			"team", e.Spec.TeamName,
			"environment", e.Spec.EnvName,
			"component", ec.Name,
		)
		if err := fileAPI.SaveFileFromString(overlay.Data, destinationFolder, overlay.Name); err != nil {
			return err
		}
	}

	return nil
}
