package apps

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/compuzest/zlifecycle-il-operator/controllers/external/argocd"

	"github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	stablev1 "github.com/compuzest/zlifecycle-il-operator/api/v1"
	"github.com/compuzest/zlifecycle-il-operator/controllers/codegen/file"
	"github.com/compuzest/zlifecycle-il-operator/controllers/external/git"
	"github.com/compuzest/zlifecycle-il-operator/controllers/gitreconciler"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util"
	perrors "github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"k8s.io/apimachinery/pkg/util/yaml"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

func GenerateArgocdApps(
	log *logrus.Entry,
	fileAPI file.API,
	gitClient git.API,
	gitReconciler gitreconciler.API,
	key *client.ObjectKey,
	e *stablev1.Environment,
	destinationFolder string,
) error {
	for _, ec := range e.Spec.Components {
		if ec.Type != "argocd" {
			continue
		}

		var apps []*v1alpha1.Application

		if err := func() error {
			tempDir, cleanup, err := git.CloneTemp(gitClient, ec.Module.Source)
			defer cleanup()
			if err != nil {
				return err
			}

			log.WithFields(logrus.Fields{
				"source":      ec.Module.Source,
				"version":     ec.Module.Version,
				"path":        ec.Module.Path,
				"destination": destinationFolder,
				"component":   ec.Name,
			}).Info("Generating argocd app")
			absolutePath := filepath.Join(tempDir, ec.Module.Path)
			if util.IsDir(absolutePath) {
				apps, err = parseArgocdApplicationFolder(absolutePath)
				if err != nil {
					return err
				}
			} else {
				app, err := parseArgocdApplicationYAML(absolutePath)
				if err != nil {
					return err
				}
				apps = append(apps, app)
			}

			if err := saveToGitRepo(apps, destinationFolder, fileAPI); err != nil {
				return err
			}

			submitToGitReconciler(gitReconciler, key, ec, log)

			return nil
		}(); err != nil {
			return err
		}
	}

	return nil
}

func saveToGitRepo(apps []*v1alpha1.Application, folderName string, fileAPI file.API) error {
	for _, app := range apps {
		fileName := app.Name + ".yaml"
		if err := fileAPI.SaveYamlFile(app, folderName, fileName); err != nil {
			return perrors.Wrapf(err, "error saving file to %s/%s", folderName, fileName)
		}
	}

	return nil
}

func submitToGitReconciler(gitReconciler gitreconciler.API, key *client.ObjectKey, ec *stablev1.EnvironmentComponent, log *logrus.Entry) {
	log.WithFields(logrus.Fields{
		"component":  ec.Name,
		"type":       ec.Type,
		"repository": ec.Module.Source,
	}).Info("Subscribing to config repository in git reconciler")
	subscribed := gitReconciler.Subscribe(ec.Module.Source, *key)
	if subscribed {
		log.WithFields(logrus.Fields{
			"component":  ec.Name,
			"type":       ec.Type,
			"repository": ec.Module.Source,
		}).Info("Already subscribed in git reconciler to repository")
	}
}

func parseArgocdApplicationFolder(path string) (apps []*v1alpha1.Application, err error) {
	walkF := func(path string, info fs.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if info.IsDir() {
			return nil
		}
		isYAML := strings.HasSuffix(info.Name(), ".yaml") || strings.HasSuffix(info.Name(), ".yml")
		if !isYAML {
			return nil
		}
		app, err := parseArgocdApplicationYAML(path)
		if err != nil {
			return perrors.Wrapf(err, "error parsing argocd application yaml from %s", path)
		}
		apps = append(apps, app)

		return nil
	}
	if err = filepath.Walk(path, walkF); err != nil {
		return nil, perrors.Wrapf(err, "error parsing argocd applications in folder %s", path)
	}

	return apps, nil
}

func parseArgocdApplicationYAML(path string) (*v1alpha1.Application, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, perrors.Wrap(err, "error reading argocd app file")
	}
	var app v1alpha1.Application
	if err = yaml.Unmarshal(data, &app); err != nil {
		return nil, perrors.Wrapf(err, "error unmarshalling argocd application yaml")
	}
	argocd.AddLabelsToCustomerApp(&app)

	return &app, nil
}
