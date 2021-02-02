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

package controllers

import (
	"fmt"
	stablev1alpha1 "github.com/compuzest/zlifecycle-il-operator/api/v1alpha1"
	"github.com/ghodss/yaml"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/util/json"
	"os"
)

func SaveYamlFile(obj interface{}, fileName string) {
	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}

	bytes, err2 := yaml.JSONToYAML(jsonBytes)
	if err2 != nil {
		panic(err2)
	}

	err3 := ioutil.WriteFile(fileName, bytes, 0644)
	if err3 != nil {
		panic(err3)
	}
}

func SaveVarsToFile(variables []*stablev1alpha1.Variable, fileName string) {
	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	for _, variable := range variables {
		fmt.Fprintf(file, "%s = \"%s\"\n", variable.Name, variable.Value)
	}
}
