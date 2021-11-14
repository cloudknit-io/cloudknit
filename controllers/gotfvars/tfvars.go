package gotfvars

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	v1 "github.com/compuzest/zlifecycle-il-operator/api/v1"
	"github.com/compuzest/zlifecycle-il-operator/controllers/gitreconciler"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/file"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/git"
	"github.com/go-logr/logr"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

func SaveTfVarsToFile(fileService file.Service, ecVars []*v1.Variable, folderName string, fileName string) error {
	variables := make([]*v1.Variable, 0, len(ecVars))
	for _, v := range ecVars {
		// TODO: This is a hack to just to make it work, needs to be revisited
		v.Value = fmt.Sprintf("\"%s\"", v.Value)
		variables = append(variables, v)
	}

	return fileService.SaveVarsToFile(variables, folderName, fileName)
}

func GetVariablesFromTfvarsFile(
	ctx context.Context,
	log logr.Logger,
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
	log.Info(
		"Reading tfvars file contents",
		"team", e.Spec.TeamName,
		"environment", e.Spec.EnvName,
		"component", ec.Name,
		"path", path,
	)
	buff, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	log.Info(
		"Subscribing to config repository in git reconciler",
		"environment", e.Spec.EnvName,
		"team", e.Spec.TeamName,
		"component", ec.Name,
		"type", ec.Type,
		"repository", ec.VariablesFile.Source,
	)
	envKey := client.ObjectKey{Name: e.Name, Namespace: e.Namespace}
	subscribed := gitreconciler.GetReconciler().Subscribe(ec.VariablesFile.Source, envKey)
	if subscribed {
		log.Info(
			"Already subscribed in git reconciler to repository",
			"environment", e.Spec.EnvName,
			"team", e.Spec.TeamName,
			"component", ec.Name,
			"type", ec.Type,
			"repository", ec.VariablesFile.Source,
		)
	}

	tfvars := string(buff)

	return tfvars, nil
}
