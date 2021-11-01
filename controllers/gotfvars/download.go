package gotfvars

import (
	"fmt"
	v1 "github.com/compuzest/zlifecycle-il-operator/api/v1"
	"github.com/compuzest/zlifecycle-il-operator/controllers/filereconciler"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/file"
	"io/ioutil"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/compuzest/zlifecycle-il-operator/controllers/util/common"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/github"
	"github.com/go-logr/logr"
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
	log logr.Logger,
	api github.RepositoryAPI,
	environment *v1.Environment,
	ec *v1.EnvironmentComponent,
) (string, error) {
	log.Info(
		"Downloading tfvars file",
		"source", ec.VariablesFile.Source,
		"ref", ec.VariablesFile.Ref,
		"path", ec.VariablesFile.Path,
	)
	buff, exists, err := downloadTfvarsFile(api, ec.VariablesFile.Source, ec.VariablesFile.Ref, ec.VariablesFile.Path)
	if err != nil {
		return "", err
	}
	if !exists {
		return "", fmt.Errorf("file does not exist: %s/%s?ref=%s", ec.VariablesFile.Source, ec.VariablesFile.Path, ec.VariablesFile.Ref)
	}

	log.Info(
		"Submitting terraform variables file to file reconciler",
		"environment", environment.Spec.EnvName,
		"team", environment.Spec.TeamName,
		"component", ec.Name,
		"type", ec.Type,
	)
	fm := &filereconciler.FileMeta{
		Type:           "tfvars",
		Filename:       ec.VariablesFile.Path,
		Environment:    environment.Spec.EnvName,
		Component:      ec.Name,
		Source:         ec.VariablesFile.Source,
		Path:           ec.VariablesFile.Path,
		Ref:            ec.VariablesFile.Ref,
		EnvironmentKey: client.ObjectKey{Name: environment.Name, Namespace: environment.Namespace},
	}
	if _, err := filereconciler.GetReconciler().Submit(fm); err != nil {
		return "", err
	}

	tfvars := string(buff)

	return tfvars, nil
}

func downloadTfvarsFile(api github.RepositoryAPI, repoURL string, ref string, path string) (file []byte, exists bool, err error) {
	rc, exists, err := github.DownloadFile(api, repoURL, ref, path)
	if err != nil {
		return nil, false, err
	}
	if !exists {
		return nil, false, nil
	}
	defer common.CloseBody(rc)
	buff, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, false, err
	}
	return buff, true, nil
}
