package overlay

import (
	stablev1 "github.com/compuzest/zlifecycle-il-operator/api/v1"
	"github.com/compuzest/zlifecycle-il-operator/controllers/filereconciler"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/common"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/file"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/github"
	"github.com/go-logr/logr"
	"io/ioutil"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func GenerateOverlayFiles(
	log logr.Logger,
	fileService file.Service,
	repoAPI github.RepositoryAPI,
	e *stablev1.Environment,
	ec *stablev1.EnvironmentComponent,
	componentFolder string,
) error {
	if ec.OverlayFiles != nil {
		for _, overlay := range ec.OverlayFiles {
			name := common.ExtractNameFromPath(overlay.Path)
			log.Info(
				"Generating overlay file from git file",
				"overlay", name,
				"folder", componentFolder,
				"source", overlay.Source,
				"path", overlay.Path,
				"component", ec.Name,
				"environment", e.Spec.EnvName,
				"team", e.Spec.TeamName,
			)
			if err := saveOverlayFileFromGit(fileService, repoAPI, overlay, componentFolder, name); err != nil {
				return err
			}
			log.Info(
				"Submitting overlay file to file reconciler",
				"filename", overlay.Path,
				"component", ec.Name,
				"environment", e.Spec.EnvName,
				"team", e.Spec.TeamName,
			)
			fm := filereconciler.FileMeta{
				Type:           "overlay",
				Filename:       overlay.Path,
				Environment:    e.Spec.EnvName,
				Component:      ec.Name,
				Source:         overlay.Source,
				Path:           overlay.Path,
				Ref:            overlay.Ref,
				EnvironmentKey: client.ObjectKey{Name: e.Name, Namespace: e.Namespace},
			}
			if _, err := filereconciler.GetReconciler().Submit(&fm); err != nil {
				return err
			}
		}
	}
	if ec.OverlayData != nil {
		for _, overlay := range ec.OverlayData {
			log.Info(
				"Generating overlay file from data field",
				"overlay", overlay.Name,
				"folder", componentFolder,
				"component", ec.Name,
			)
			if err := fileService.SaveFileFromString(overlay.Data, componentFolder, overlay.Name); err != nil {
				return err
			}
		}
	}
	return nil
}

func saveOverlayFileFromGit(
	fileUtil file.Service,
	repoAPI github.RepositoryAPI,
	of *stablev1.OverlayFile,
	folderName string,
	fileName string,
) error {
	ref := of.Ref
	if ref == "" {
		ref = "HEAD"
	}
	f, _, err := github.DownloadFile(repoAPI, of.Source, ref, of.Path)
	if err != nil {
		return err
	}

	buff, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	return fileUtil.SaveFileFromByteArray(buff, folderName, fileName)
}
