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
			tempDir, cleanup, err := git.CloneTemp(gitClient, ec.Module.Source, log)
			defer cleanup()
			if err != nil {
				return err
			}

			sourceAbsolutePath := filepath.Join(tempDir, ec.Module.Path)
			if util.IsDir(sourceAbsolutePath) {
				apps, err = parseArgocdApplicationFolder(sourceAbsolutePath, e, ec)
				if err != nil {
					return err
				}
			} else {
				app, err := parseArgocdApplicationYAML(filepath.Base(sourceAbsolutePath), sourceAbsolutePath, e, ec)
				if err != nil {
					return err
				}
				apps = append(apps, app)
			}

			log.WithFields(logrus.Fields{
				"source":      ec.Module.Source,
				"version":     ec.Module.Version,
				"path":        ec.Module.Path,
				"destination": destinationFolder,
				"component":   ec.Name,
			}).Infof("Generating ArgoCD App of Apps Helm chart for environment component %s", ec.Name)
			if err := generateHelmChart(fileAPI, destinationFolder, ec.Name); err != nil {
				return err
			}

			log.WithFields(logrus.Fields{
				"source":      ec.Module.Source,
				"version":     ec.Module.Version,
				"path":        ec.Module.Path,
				"destination": destinationFolder,
				"component":   ec.Name,
			}).Infof("Generating ArgoCD Applications in  App of Apps Helm chart for environment component %s", ec.Name)
			if err := generateArgocdApplications(apps, destinationFolder, fileAPI); err != nil {
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

func generateArgocdApplications(apps []*v1alpha1.Application, folderName string, fileAPI file.API) error {
	templatesDir := filepath.Join(folderName, "templates")
	for _, app := range apps {
		fileName := app.Labels["source_file_name"]
		if fileName == "" {
			fileName = app.Name + ".yaml"
		}
		if err := fileAPI.SaveYamlFile(app, templatesDir, fileName); err != nil {
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

func generateHelmChart(fileAPI file.API, dir string, name string) error {
	chart := NewHelmChart(name)
	if err := fileAPI.SaveYamlFile(chart, dir, "Chart.yaml"); err != nil {
		return perrors.Wrapf(err, "error generating chart yaml")
	}

	return nil
}

func parseArgocdApplicationFolder(path string, e *stablev1.Environment, ec *stablev1.EnvironmentComponent) (apps []*v1alpha1.Application, err error) {
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
		app, err := parseArgocdApplicationYAML(info.Name(), path, e, ec)
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

func parseArgocdApplicationYAML(filename, path string, e *stablev1.Environment, ec *stablev1.EnvironmentComponent) (*v1alpha1.Application, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, perrors.Wrap(err, "error reading argocd app file")
	}
	var app v1alpha1.Application
	if err = yaml.Unmarshal(data, &app); err != nil {
		return nil, perrors.Wrapf(err, "error unmarshalling argocd application yaml")
	}
	argocd.AddLabelsToCustomerApp(&app, e, ec, filename)

	return &app, nil
}
