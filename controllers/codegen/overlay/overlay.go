package overlay

import (
	"path/filepath"

	"github.com/compuzest/zlifecycle-il-operator/controllers/codegen/file"

	"github.com/compuzest/zlifecycle-il-operator/controllers/external/git"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util"

	"github.com/sirupsen/logrus"

	"github.com/compuzest/zlifecycle-il-operator/controllers/gitreconciler"

	"sigs.k8s.io/controller-runtime/pkg/client"

	stablev1 "github.com/compuzest/zlifecycle-il-operator/api/v1"
)

func GenerateOverlayFiles(
	log *logrus.Entry,
	fileAPI file.API,
	gitClient git.API,
	gitReconciler gitreconciler.API,
	e *stablev1.Environment,
	ec *stablev1.EnvironmentComponent,
	destinationFolder string,
) error {
	// generate overlay files from config repo
	for _, overlay := range ec.OverlayFiles {
		if err := func() error {
			tempDir, cleanup, err := git.CloneTemp(gitClient, overlay.Source)
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
				}).Info("Generating overlay file")
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

				// submit repo to git reconciler
				log.WithFields(logrus.Fields{
					"component":  ec.Name,
					"type":       ec.Type,
					"repository": overlay.Source,
				}).Info("Subscribing to config repository in git reconciler")
				envKey := client.ObjectKey{Name: e.Name, Namespace: e.Namespace}
				subscribed := gitReconciler.Subscribe(overlay.Source, envKey)
				if subscribed {
					log.WithFields(logrus.Fields{
						"component":  ec.Name,
						"type":       ec.Type,
						"repository": overlay.Source,
					}).Info("Already subscribed in git reconciler to repository")
				}
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
		}).Info("Generating overlay file from data field")
		if err := fileAPI.SaveFileFromString(overlay.Data, destinationFolder, overlay.Name); err != nil {
			return err
		}
	}

	return nil
}
