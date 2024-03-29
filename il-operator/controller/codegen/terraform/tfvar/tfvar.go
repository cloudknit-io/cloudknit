package tfvar

import (
	"fmt"
	"os"
	"path/filepath"

	git2 "github.com/compuzest/zlifecycle-il-operator/controller/common/git"
	"github.com/compuzest/zlifecycle-il-operator/controller/services/operations/git"

	"github.com/compuzest/zlifecycle-il-operator/controller/services/gitreconciler"

	"github.com/compuzest/zlifecycle-il-operator/controller/codegen/interpolator"

	"github.com/pkg/errors"

	"github.com/compuzest/zlifecycle-il-operator/controller/codegen/file"

	"github.com/sirupsen/logrus"

	v1 "github.com/compuzest/zlifecycle-il-operator/api/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type TFVars struct {
	Variables []*Variable
}

type Variable struct {
	Name  string
	Value any
}

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
	gitClient git2.API,
	log *logrus.Entry,
	key *client.ObjectKey,
	ec *v1.EnvironmentComponent,
	zlocals []*v1.LocalVariable,
) (string, error) {
	tempRepoDir, cleanup, err := git.CloneTemp(gitClient, ec.VariablesFile.Source, log)
	if err != nil {
		return "", errors.Wrapf(err, "error temp cloning repo [%s]", ec.VariablesFile.Source)
	}
	defer cleanup()

	path := filepath.Join(tempRepoDir, ec.VariablesFile.Path)
	log.WithFields(logrus.Fields{
		"component": ec.Name,
		"path":      path,
	}).Infof("Reading tfvars file contents for environment component %s", ec.Name)
	buff, err := os.ReadFile(path)
	if err != nil {
		return "", errors.Wrapf(err, "error reading file [%s]", path)
	}

	tfvars := string(buff)
	interpolated, err := interpolator.InterpolateTFVars(tfvars, zlocals)
	if err != nil {
		return "", err
	}

	submitToGitReconciler(gitReconciler, key, ec, log)

	return interpolated, nil
}

func submitToGitReconciler(gitReconciler gitreconciler.API, key *client.ObjectKey, ec *v1.EnvironmentComponent, log *logrus.Entry) {
	log.WithFields(logrus.Fields{
		"component":  ec.Name,
		"type":       ec.Type,
		"repository": ec.VariablesFile.Source,
	}).Infof("Subscribing to config repository %s in git reconciler", ec.VariablesFile.Source)
	subscribed := gitReconciler.Subscribe(ec.VariablesFile.Source, *key)
	if subscribed {
		log.WithFields(logrus.Fields{
			"component":  ec.Name,
			"type":       ec.Type,
			"repository": ec.VariablesFile.Source,
		}).Info("Already subscribed in git reconciler to repository")
	}
}
