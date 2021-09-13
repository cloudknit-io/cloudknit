/* Copyright (C) 2020 CompuZest, Inc. - All Rights Reserved
 *
 * Unauthorized copying of this file, via any medium, is strictly prohibited
 * Proprietary and confidential
 *
 * NOTICE: All information contained herein is, and remains the property of
 * CompuZest, Inc. The intellectual and technical concepts contained herein are
 * proprietary to CompuZest, Inc. and are protected by trade secret or copyright
 * law. Dissemination of this information or reproduction of this material is
 * strictly forbidden unless prior written permission is obtained from CompuZest, Inc.
 */

package file

import (
	"fmt"
	"io/ioutil"
	"os"
	"text/template"

	stablev1 "github.com/compuzest/zlifecycle-il-operator/api/v1"
	"github.com/ghodss/yaml"
	"github.com/go-logr/logr"
	_ "github.com/golang/mock/mockgen/model"
	"k8s.io/apimachinery/pkg/util/json"
)

//go:generate go run -mod=mod github.com/golang/mock/mockgen -source=./file.go -destination=../../../mocks/mock_file.go -package=mocks github.com/compuzest/zlifecycle-il-operator/mocks

type UtilFile interface {
	SaveFileFromString(input string, folderName string, fileName string) error
	SaveFileFromByteArray(input []byte, folderName string, fileName string) error
	SaveYamlFile(obj interface{}, folderName string, fileName string) error
	SaveVarsToFile(variables []*stablev1.Variable, folderName string, fileName string) error
	CreateEmptyDirectory(folderName string) error
	SaveFileFromTemplate(t *template.Template, vars interface{}, folderName string, fileName string) error
	RemoveAll(path string) error
}

type UtilFileService struct {
	UtilFile
	log logr.Logger
}

func (f UtilFileService) SaveVarsToFile(variables []*stablev1.Variable, folderName string, fileName string) error {
	file, err := createFileRecursive(folderName, fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, variable := range variables {
		if _, err := fmt.Fprintf(file, "%s = \"%s\"\n", variable.Name, variable.Value); err != nil {
			return err
		}
	}

	return nil
}

// SaveFileFromTemplate creates a file with variables
func (f UtilFileService) SaveFileFromTemplate(t *template.Template, vars interface{}, folderName string, fileName string) error {
	file, err := createFileRecursive(folderName, fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	err = t.Execute(file, vars)

	if err != nil {
		return err
	}
	return nil
}

func createFileRecursive(folderName string, fileName string) (*os.File, error) {
	if err := os.MkdirAll(folderName, os.ModePerm); err != nil {
		return nil, fmt.Errorf("error: failed to create directory: %w", err)
	}
	file, err := os.Create(fmt.Sprintf("%s/%s", folderName, fileName))

	if err != nil {
		return nil, err
	}

	return file, nil
}

// SaveFileFromString Create file
func (f UtilFileService) SaveFileFromString(input string, folderName string, fileName string) error {
	if err := saveBytesToFile(folderName, fileName, []byte(input)); err != nil {
		return err
	}

	return nil
}

func (f UtilFileService) SaveFileFromByteArray(input []byte, folderName string, fileName string) error {
	if err := saveBytesToFile(folderName, fileName, input); err != nil {
		return err
	}

	return nil
}

// CreateEmptyDirectory creates empty directory with a .keep file
func (f UtilFileService) CreateEmptyDirectory(folderName string) error {
	if err := saveBytesToFile(folderName, ".keep", nil); err != nil {
		return err
	}

	return nil
}

// SaveYamlFile creates file and directory, does not validate yaml
func (f UtilFileService) SaveYamlFile(obj interface{}, folderName string, fileName string) error {
	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		return fmt.Errorf("error: failed to marshal json: %s", err.Error())
	}

	bytes, err := yaml.JSONToYAML(jsonBytes)
	if err != nil {
		return fmt.Errorf("error: failed to convert json to yaml: %s", err.Error())
	}

	if err := saveBytesToFile(folderName, fileName, bytes); err != nil {
		return err
	}

	return nil
}

func saveBytesToFile(folderName string, fileName string, bytes []byte) error {
	if err := os.MkdirAll(folderName, os.ModePerm); err != nil {
		return fmt.Errorf("error: failed to create directory: %w", err)
	}

	if err := ioutil.WriteFile(fmt.Sprintf("%s/%s", folderName, fileName), bytes, 0644); err != nil {
		return fmt.Errorf("error: failed to write file: %w", err)
	}

	return nil
}

// RemoveAll removes path and any children it contains.
// It removes everything it can but returns the first error
// it encounters. If the path does not exist, RemoveAll
// returns nil (no error).
func (f UtilFileService) RemoveAll(path string) error {
	return os.RemoveAll(path)
}
