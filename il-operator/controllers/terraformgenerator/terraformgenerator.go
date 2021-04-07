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

package terraformgenerator

import (
	"os"
	"path/filepath"
	"text/template"

	stablev1alpha1 "github.com/compuzest/zlifecycle-il-operator/api/v1alpha1"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/file"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/il"
)

// UtilTerraformGenerator package interface for generating terraform files
type UtilTerraformGenerator interface {
	GenerateTerraform(fileUtil file.UtilFile, environmentComponent *stablev1alpha1.EnvironmentComponent, environment *stablev1alpha1.Environment, environmentComponentDirectory string) error
	GenerateProvider(file file.UtilFile, environmentComponentDirectory string) error
	GenerateFromTemplate(vars interface{}, environmentComponentDirectory string, fileUtil file.UtilFile, templateName string, filePath string) error
}

type TerraformGenerator struct {
	UtilTerraformGenerator
}

var DefaultTerraformVersion = "0.13.2"

func TerraformIlPath(envComponentDirectory string) string {
	return envComponentDirectory + "/terraform"
}

func (tf TerraformGenerator) GenerateTerraform(fileUtil file.UtilFile, environmentComponent *stablev1alpha1.EnvironmentComponent, environment *stablev1alpha1.Environment, environmentComponentDirectory string) error {
	componentName := environmentComponent.Module.Name

	backendConfig := TerraformBackendConfig{
		Region:        "us-east-1",
		Profile:       "compuzest-shared",
		Version:       DefaultTerraformVersion,
		Bucket:        "compuzest-zlifecycle-tfstate",
		DynamoDBTable: "compuzest-zlifecycle-tflock",
		TeamName:      environment.Spec.TeamName,
		EnvName:       environment.Spec.EnvName,
		ComponentName: componentName,
	}

	moduleConfig := TerraformModuleConfig{
		ComponentName: componentName,
		Source:        il.EnvComponentModuleSource(environmentComponent.Module.Source, componentName),
		Path:          il.EnvComponentModulePath(environmentComponent.Module.Path),
		Variables:     environmentComponent.Variables,
	}

	err := tf.GenerateProvider(fileUtil, environmentComponentDirectory, componentName)
	if err != nil {
		return err
	}

	workingDir, err := os.Getwd()
	err = tf.GenerateFromTemplate(backendConfig, environmentComponentDirectory, componentName, fileUtil, filepath.Join(workingDir, "templates/terraform_backend.tmpl"), "terraform")
	if err != nil {
		return err
	}

	err = tf.GenerateFromTemplate(moduleConfig, environmentComponentDirectory, componentName, fileUtil, filepath.Join(workingDir, "templates/terraform_module.tmpl"), "module")
	if err != nil {
		return err
	}
	return nil
}

// TerraformBackendConfig variables for creating tf backend
type TerraformBackendConfig struct {
	Region        string
	Version       string
	Key           string
	Bucket        string
	DynamoDBTable string
	Profile       string
	TeamName      string
	EnvName       string
	ComponentName string
}

// TerraformModuleConfig variables for creating tf module
type TerraformModuleConfig struct {
	ComponentName string
	Source        string
	Path          string
	Variables     []*stablev1alpha1.Variable
}

// GenerateProvider save provider file to be executed by terraform
func (tf TerraformGenerator) GenerateProvider(file file.UtilFile, environmentComponentDirectory string, componentName string) error {
	terraformDirectory := tf.generateTerraformDirectory(environmentComponentDirectory, componentName)
	err := file.SaveFileFromString(`
provider "aws" {
	region = "us-east-1"
	version = "~> 3.0"
}
	`, terraformDirectory, "provider.tf")
	if err != nil {
		return err
	}
	return nil
}

// GenerateFromTemplate save terraform backend config
func (tf TerraformGenerator) generateTerraformDirectory(environmentComponentDirectory string, environmentComponentName string) string {
	return environmentComponentDirectory + "/" + environmentComponentName + "/terraform"
}

func (tf TerraformGenerator) GenerateFromTemplate(vars interface{}, environmentComponentDirectory string, componentName string, fileUtil file.UtilFile, templatePath string, fileName string) error {
	template, err := template.ParseFiles(templatePath)
	if err != nil {
		return err
	}
	terraformDirectory := tf.generateTerraformDirectory(environmentComponentDirectory, componentName)
	err = fileUtil.SaveFileFromTemplate(template, vars, terraformDirectory, fileName+".tf")
	if err != nil {
		return err
	}
	return nil
}
