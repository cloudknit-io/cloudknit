package gotfvars

import (
	"fmt"
	"os"
	"path/filepath"

	perrors "github.com/pkg/errors"

	"github.com/compuzest/zlifecycle-il-operator/controllers/codegen/file"

	"github.com/compuzest/zlifecycle-il-operator/controllers/external/git"
	"github.com/sirupsen/logrus"

	v1 "github.com/compuzest/zlifecycle-il-operator/api/v1"
	"github.com/compuzest/zlifecycle-il-operator/controllers/gitreconciler"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func GenerateTFVarsFile(fs file.API, vars []*v1.Variable, folderName string, fileName string) error {
	variables := make([]*v1.Variable, 0, len(vars))
	for _, v := range vars {
		// TODO: This is a hack to just to make it work, needs to be revisited
		v.Value = fmt.Sprintf("%q", v.Value)
		variables = append(variables, v)
	}

	return fs.SaveVarsToFile(variables, folderName, fileName)
}

func GetVariablesFromTfvarsFile(
	gitReconciler gitreconciler.API,
	gitClient git.API,
	log *logrus.Entry,
	key *client.ObjectKey,
	ec *v1.EnvironmentComponent,
) (string, error) {
	tempRepoDir, cleanup, err := git.CloneTemp(gitClient, ec.VariablesFile.Source)
	if err != nil {
		return "", perrors.Wrapf(err, "error temp cloning repo [%s]", ec.VariablesFile.Source)
	}
	defer cleanup()

	path := filepath.Join(tempRepoDir, ec.VariablesFile.Path)
	log.WithFields(logrus.Fields{
		"component": ec.Name,
		"path":      path,
	}).Info("Reading tfvars file contents")
	buff, err := os.ReadFile(path)
	if err != nil {
		return "", perrors.Wrapf(err, "error reading file [%s]", path)
	}

	log.WithFields(logrus.Fields{
		"component":  ec.Name,
		"type":       ec.Type,
		"repository": ec.VariablesFile.Source,
	}).Info("Subscribing to config repository in git reconciler")
	subscribed := gitReconciler.Subscribe(ec.VariablesFile.Source, *key)
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
