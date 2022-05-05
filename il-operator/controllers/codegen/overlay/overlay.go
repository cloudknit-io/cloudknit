package overlay

import (
	"path/filepath"

	"github.com/compuzest/zlifecycle-il-operator/controllers/lib/gitreconciler"

	"github.com/compuzest/zlifecycle-il-operator/controllers/codegen/file"

	"github.com/compuzest/zlifecycle-il-operator/controllers/external/git"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util"

	"github.com/sirupsen/logrus"

	"sigs.k8s.io/controller-runtime/pkg/client"

	stablev1 "github.com/compuzest/zlifecycle-il-operator/api/v1"
)

func GenerateOverlayFiles(
	log *logrus.Entry,
	fileAPI file.API,
	gitClient git.API,
	gitReconciler gitreconciler.API,
	key *client.ObjectKey,
	ec *stablev1.EnvironmentComponent,
	destinationFolder string,
) error {
	// generate overlay files from config repo
	for _, overlay := range ec.OverlayFiles {
		if err := func() error {
			tempDir, cleanup, err := git.CloneTemp(gitClient, overlay.Source, log)
			defer cleanup()
			if err != nil {
				return err
			}
			for _, path := range overlay.Paths {
				log.WithFields(logrus.Fields{
					"repo":        overlay.Source,
					"ref":         overlay.Ref,
					"source":      path,
					"destination": destinationFolder,
					"component":   ec.Name,
				}).Infof("Generating overlay file(s) for environment component %s", ec.Name)
				absolutePath := filepath.Join(tempDir, path)
				if util.IsDir(absolutePath) {
					if err := util.CopyDirContent(absolutePath, destinationFolder); err != nil {
						return err
					}
				} else {
					if err := util.CopyFile(absolutePath, destinationFolder); err != nil {
						return err
					}
				}
				submitToGitReconciler(gitReconciler, key, ec, log)
			}

			return nil
		}(); err != nil {
			return err
		}
	}
	// Generate inline overlay files
	for _, overlay := range ec.OverlayData {
		log.WithFields(logrus.Fields{
			"overlay":     overlay.Name,
			"destination": destinationFolder,
			"component":   ec.Name,
		}).Infof("Generating overlay file from data field for environment component %s", ec.Name)
		if err := fileAPI.SaveFileFromString(overlay.Data, destinationFolder, overlay.Name); err != nil {
			return err
		}
	}

	return nil
}

func submitToGitReconciler(gitReconciler gitreconciler.API, key *client.ObjectKey, ec *stablev1.EnvironmentComponent, log *logrus.Entry) {
	log.WithFields(logrus.Fields{
		"component":  ec.Name,
		"type":       ec.Type,
		"repository": ec.Module.Source,
	}).Infof("Subscribing to config repository %s in git reconciler", ec.Module.Source)
	subscribed := gitReconciler.Subscribe(ec.Module.Source, *key)
	if subscribed {
		log.WithFields(logrus.Fields{
			"component":  ec.Name,
			"type":       ec.Type,
			"repository": ec.Module.Source,
		}).Info("Already subscribed in git reconciler to repository")
	}
}
