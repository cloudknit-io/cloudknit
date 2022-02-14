package gotfvars

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"

	v1 "github.com/compuzest/zlifecycle-il-operator/api/v1"
	"github.com/compuzest/zlifecycle-il-operator/controllers/gitreconciler"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/file"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/git"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func SaveTfVarsToFile(fs file.FSAPI, vars []*v1.Variable, folderName string, fileName string) error {
	variables := make([]*v1.Variable, 0, len(vars))
	for _, v := range vars {
		// TODO: This is a hack to just to make it work, needs to be revisited
		v.Value = fmt.Sprintf("%q", v.Value)
		variables = append(variables, v)
	}

	return fs.SaveVarsToFile(variables, folderName, fileName)
}

func GetVariablesFromTfvarsFile(
	ctx context.Context,
	gitReconciler gitreconciler.API,
	log *logrus.Entry,
	e *v1.Environment,
	ec *v1.EnvironmentComponent,
) (string, error) {
	gitAPI, err := git.NewGoGit(ctx)
	if err != nil {
		return "", err
	}
	tempRepoDir, cleanup, err := git.CloneTemp(gitAPI, ec.VariablesFile.Source)
	if err != nil {
		return "", err
	}
	defer cleanup()

	path := filepath.Join(tempRepoDir, ec.VariablesFile.Path)
	log.WithFields(logrus.Fields{
		"component": ec.Name,
		"path":      path,
	}).Info("Reading tfvars file contents")
	buff, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	log.WithFields(logrus.Fields{
		"component":  ec.Name,
		"type":       ec.Type,
		"repository": ec.VariablesFile.Source,
	}).Info("Subscribing to config repository in git reconciler")
	envKey := client.ObjectKey{Name: e.Name, Namespace: e.Namespace}
	subscribed := gitReconciler.Subscribe(ec.VariablesFile.Source, envKey)
	if subscribed {
		log.WithFields(logrus.Fields{
			"component":  ec.Name,
			"type":       ec.Type,
			"repository": ec.VariablesFile.Source,
		}).Info("Already subscribed in git reconciler to repository")
	}

	tfvars := string(buff)

	return tfvars, nil
}
