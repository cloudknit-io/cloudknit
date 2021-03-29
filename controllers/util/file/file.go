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

	"github.com/go-logr/logr"

	stablev1alpha1 "github.com/compuzest/zlifecycle-il-operator/api/v1alpha1"
	"github.com/ghodss/yaml"
	"k8s.io/apimachinery/pkg/util/json"
)

//go:generate mockgen -source=./file.go -destination=../../../mocks/mock_file.go -package=mocks github.com/compuzest/zlifecycle-il-operator/mocks

type UtilFile interface {
	SaveFileFromString(jsonString string, folderName string, fileName string) error
	SaveYamlFile(obj interface{}, folderName string, fileName string) error
	SaveVarsToFile(variables []*stablev1alpha1.Variable, folderName string, fileName string) error
	CreateEmptyDirectory(folderName string) error
	SaveFileFromTemplate(t *template.Template, vars interface{}, folderName string, fileName string) error
}

type UtilFileService struct {
	UtilFile
	log logr.Logger
}

// CreateEmptyDirectory creates empty directory with a .keep file
func (f UtilFileService) CreateEmptyDirectory(folderName string) error {
	if err := os.MkdirAll(folderName, os.ModePerm); err != nil {
		return fmt.Errorf("error: failed to create directory: %s", err.Error())
	}

	if err := ioutil.WriteFile(folderName+"/.keep", nil, 0644); err != nil {
		return fmt.Errorf("error: failed to write .keep file: %s", err.Error())
	}

	return nil
}

// SaveFileFromTemplate creates a file with variables
func (f UtilFileService) SaveFileFromTemplate(t *template.Template, vars interface{}, folderName string, fileName string) error {
	if err := os.MkdirAll(folderName, os.ModePerm); err != nil {
		return fmt.Errorf("error: failed to create directory: %s", err.Error())
	}
	createdFile, err := os.Create(folderName + "/" + fileName)

	if err != nil {
		return err
	}
	err = t.Execute(createdFile, vars)

	if err != nil {
		return err
	}
	createdFile.Close()
	return nil
}

// SaveFileFromString Create file
func (f UtilFileService) SaveFileFromString(str string, folderName string, fileName string) error {
	if err := os.MkdirAll(folderName, os.ModePerm); err != nil {
		return fmt.Errorf("error: failed to create directory: %s", err.Error())
	}

	if err := ioutil.WriteFile(folderName+"/"+fileName, []byte(str), 0644); err != nil {
		return fmt.Errorf("error: failed to write file: %s", err.Error())
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

	if err := os.MkdirAll(folderName, os.ModePerm); err != nil {
		return fmt.Errorf("error: failed to create directory: %s", err.Error())
	}

	if err := ioutil.WriteFile(folderName+"/"+fileName, bytes, 0644); err != nil {
		return fmt.Errorf("error: failed to write file: %s", err.Error())
	}

	return nil
}

func (f UtilFileService) SaveVarsToFile(variables []*stablev1alpha1.Variable, folderName string, fileName string) error {
	if err := os.MkdirAll(folderName, os.ModePerm); err != nil {
		return fmt.Errorf("error: failed to create directory: %s", err.Error())
	}

	file, err := os.Create(folderName + "/" + fileName)
	if err != nil {
		return fmt.Errorf("error: failed to create vars file: %s", err.Error())
	}

	defer file.Close()

	for _, variable := range variables {
		fmt.Fprintf(file, "%s = \"%s\"\n", variable.Name, variable.Value)
	}

	return nil
}
