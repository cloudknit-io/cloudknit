package overlay

import (
	stablev1 "github.com/compuzest/zlifecycle-il-operator/api/v1"
	"github.com/compuzest/zlifecycle-il-operator/controllers/filereconciler"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/file"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/github"
	"github.com/go-logr/logr"
	"io/ioutil"
	"strings"

	kClient "sigs.k8s.io/controller-runtime/pkg/client"
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
			tokens := strings.Split(overlay.Path, "/")
			name := tokens[len(tokens)-1]
			log.Info(
				"Generating overlay file from git file",
				"overlay", name,
				"folder", componentFolder,
				"source", overlay.Source,
				"path", overlay.Path,
				"component", ec.Name,
			)
			if err := saveOverlayFileFromGit(fileService, repoAPI, overlay, componentFolder, name); err != nil {
				return err
			}
			fm := filereconciler.FileMeta{
				Environment:       e.Spec.EnvName,
				Component:         ec.Name,
				Source:            ec.VariablesFile.Source,
				Path:              ec.VariablesFile.Path,
				Ref:               ec.VariablesFile.Ref,
				EnvironmentObject: kClient.ObjectKey{Name: e.Name, Namespace: e.Namespace},
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
