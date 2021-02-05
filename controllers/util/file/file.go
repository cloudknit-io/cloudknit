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

	stablev1alpha1 "github.com/compuzest/zlifecycle-il-operator/api/v1alpha1"
	"github.com/ghodss/yaml"
	"k8s.io/apimachinery/pkg/util/json"
)

func SaveYamlFile(obj interface{}, folderName string, fileName string) error {
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

func SaveVarsToFile(variables []*stablev1alpha1.Variable, fileName string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("error: failed to create vars file: %s", err.Error())
	}

	defer file.Close()

	for _, variable := range variables {
		fmt.Fprintf(file, "%s = \"%s\"\n", variable.Name, variable.Value)
	}

	return nil
}
