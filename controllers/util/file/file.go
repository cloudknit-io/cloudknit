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
	_ "github.com/golang/mock/mockgen/model" // workaround for mockgen failing
	"k8s.io/apimachinery/pkg/util/json"
)

//go:generate go run -mod=mod github.com/golang/mock/mockgen -source=./file.go -destination=../../../mocks/mock_file.go -package=mocks github.com/compuzest/zlifecycle-il-operator/mocks

type Service interface {
	SaveFileFromString(input string, folderName string, fileName string) error
	SaveFileFromByteArray(input []byte, folderName string, fileName string) error
	SaveYamlFile(obj interface{}, folderName string, fileName string) error
	SaveVarsToFile(variables []*stablev1.Variable, folderName string, fileName string) error
	CreateEmptyDirectory(folderName string) error
	SaveFileFromTemplate(t *template.Template, vars interface{}, folderName string, fileName string) error
	NewFile(folderName string, fileName string) (*os.File, error)
	RemoveAll(path string) error
}

type OsFileService struct{}

func NewOsFileService() *OsFileService {
	return &OsFileService{}
}

func (f *OsFileService) NewFile(folderName string, fileName string) (*os.File, error) {
	file, err := createFileRecursive(folderName, fileName)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (f *OsFileService) SaveVarsToFile(variables []*stablev1.Variable, folderName string, fileName string) error {
	file, err := createFileRecursive(folderName, fileName)
	if err != nil {
		return err
	}
	defer closeFile(file)

	for _, variable := range variables {
		if _, err := fmt.Fprintf(file, "%s = \"%s\"\n", variable.Name, variable.Value); err != nil {
			return err
		}
	}

	return nil
}

// SaveFileFromTemplate creates a file with variables
func (f *OsFileService) SaveFileFromTemplate(t *template.Template, vars interface{}, folderName string, fileName string) error {
	file, err := createFileRecursive(folderName, fileName)
	if err != nil {
		return err
	}
	defer closeFile(file)
	err = t.Execute(file, vars)

	if err != nil {
		return err
	}
	return nil
}

func closeFile(file *os.File) {
	_ = file.Close()
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
func (f *OsFileService) SaveFileFromString(input string, folderName string, fileName string) error {
	return saveBytesToFile([]byte(input), folderName, fileName)
}

func (f *OsFileService) SaveFileFromByteArray(input []byte, folderName string, fileName string) error {
	return saveBytesToFile(input, folderName, fileName)
}

// CreateEmptyDirectory creates empty directory with a .keep file
func (f *OsFileService) CreateEmptyDirectory(folderName string) error {
	return saveBytesToFile(nil, folderName, ".keep")
}

// SaveYamlFile creates file and directory, does not validate yaml
func (f *OsFileService) SaveYamlFile(obj interface{}, folderName string, fileName string) error {
	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		return fmt.Errorf("error: failed to marshal json: %w", err)
	}

	bytes, err := yaml.JSONToYAML(jsonBytes)
	if err != nil {
		return fmt.Errorf("error: failed to convert json to yaml: %w", err)
	}

	return saveBytesToFile(bytes, folderName, fileName)
}

func saveBytesToFile(bytes []byte, folderName string, fileName string) error {
	if err := os.MkdirAll(folderName, os.ModePerm); err != nil {
		return fmt.Errorf("error: failed to create directory: %w", err)
	}

	return ioutil.WriteFile(fmt.Sprintf("%s/%s", folderName, fileName), bytes, 0644)
}

// RemoveAll removes path and any children it contains.
// It removes everything it can but returns the first error
// it encounters. If the path does not exist, RemoveAll
// returns nil (no error).
func (f *OsFileService) RemoveAll(path string) error {
	return os.RemoveAll(path)
}
