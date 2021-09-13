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
	"github.com/compuzest/zlifecycle-il-operator/controllers/secrets"
	"os"
	"path/filepath"
	"text/template"

	"github.com/go-logr/logr"

	stablev1 "github.com/compuzest/zlifecycle-il-operator/api/v1"
	env "github.com/compuzest/zlifecycle-il-operator/controllers/util/env"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/file"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/il"
)

// UtilTerraformGenerator package interface for generating terraform files
type UtilTerraformGenerator interface {
	GenerateTerraform(fileUtil file.Service, environmentComponent *stablev1.EnvironmentComponent, environment *stablev1.Environment, environmentComponentDirectory string) error
	GenerateProvider(file file.Service, environmentComponentDirectory string, componentName string) error
	GenerateSharedProvider(file file.Service, environmentComponentDirectory string, componentName string) error
	GenerateFromTemplate(vars interface{}, environmentComponentDirectory string, componentName string, fileUtil file.Service, templateName string, filePath string) error
}

type TerraformGenerator struct {
	UtilTerraformGenerator
	Log logr.Logger
}

type TemplateVariables struct {
	TeamName             string
	EnvName              string
	EnvCompName          string
	EnvCompVariables     []*stablev1.Variable
	EnvCompVariablesFile string
	EnvCompSecrets       []*stablev1.Secret
	SecretScope          string
	EnvCompModuleSource  string
	EnvCompModulePath    string
	EnvCompModuleName    string
	EnvCompOutputs       []*stablev1.Output
	EnvCompDependsOn     []string
}

var DefaultTerraformVersion = "0.13.2"

func (tf TerraformGenerator) GenerateTerraform(
	fileUtil file.Service,
	vars TemplateVariables,
	environmentComponentDirectory string,
) error {
	componentName := vars.EnvCompName

	backendConfig := TerraformBackendConfig{
		Region:        "us-east-1",
		Profile:       "compuzest-shared",
		Version:       DefaultTerraformVersion,
		Bucket:        "zlifecycle-tfstate-" + env.Config.CompanyName,
		DynamoDBTable: "zlifecycle-tflock-" + env.Config.CompanyName,
		TeamName:      vars.TeamName,
		EnvName:       vars.EnvName,
		ComponentName: componentName,
	}

	moduleConfig := TerraformModuleConfig{
		ComponentName: componentName,
		Source:        il.EnvComponentModuleSource(vars.EnvCompModuleSource, vars.EnvCompModuleName),
		Path:          il.EnvComponentModulePath(vars.EnvCompModulePath),
		Variables:     vars.EnvCompVariables,
		VariablesFile: vars.EnvCompVariablesFile,
		Secrets:       vars.EnvCompSecrets,
	}

	outputsConfig := TerraformOutputsConfig{
		ComponentName: componentName,
		Outputs:       vars.EnvCompOutputs,
	}

	dataConfig := TerraformDataConfig{
		Region:    "us-east-1",
		Profile:   "compuzest-shared",
		Bucket:    "zlifecycle-tfstate-" + env.Config.CompanyName,
		TeamName:  vars.TeamName,
		EnvName:   vars.EnvName,
		DependsOn: vars.EnvCompDependsOn,
	}

	secretsMeta := secrets.SecretMeta{
		Company:              env.Config.CompanyName,
		Team:                 vars.TeamName,
		Environment:          vars.EnvName,
		EnvironmentComponent: vars.EnvCompName,
	}
	secretsConfig, err := createSecretsConfig(vars.EnvCompSecrets, secretsMeta)
	if err != nil {
		return err
	}

	if err := tf.GenerateProvider(fileUtil, environmentComponentDirectory, componentName); err != nil {
		return err
	}

	if err := tf.GenerateSharedProvider(fileUtil, environmentComponentDirectory, componentName); err != nil {
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

	if len(outputsConfig.Outputs) > 0 {
		if err := tf.GenerateFromTemplate(outputsConfig, environmentComponentDirectory, componentName, fileUtil, filepath.Join(workingDir, "templates/terraform_outputs.tmpl"), "outputs"); err != nil {
			return err
		}
	}

	if len(dataConfig.DependsOn) > 0 {
		if err := tf.GenerateFromTemplate(dataConfig, environmentComponentDirectory, componentName, fileUtil, filepath.Join(workingDir, "templates/terraform_data.tmpl"), "data"); err != nil {
			return err
		}
	}

	if vars.EnvCompSecrets != nil && len(vars.EnvCompSecrets) > 0 {
		if err := tf.GenerateFromTemplate(secretsConfig, environmentComponentDirectory, componentName, fileUtil, filepath.Join(workingDir, "templates/terraform_secrets.tmpl"), "secrets"); err != nil {
			return err
		}
	}

	return nil
}

func createSecretsConfig(secretArray []*stablev1.Secret, meta secrets.SecretMeta) (*TerraformSecretsConfig, error) {
	var scopedSecrets []Secret
	for _, s := range secretArray {
		scope := s.Scope
		if scope == "" {
			scope = "component"
		}
		key, err := secrets.CreateKey(s.Key, scope, meta)
		if err != nil {
			return nil, err
		}
		scopedSecrets = append(scopedSecrets, Secret{Key: key, Name: s.Name})
	}
	conf := TerraformSecretsConfig{Secrets: scopedSecrets}

	return &conf, nil
}

// GenerateSharedProvider save provider file to be executed by terraform
func (tf TerraformGenerator) GenerateSharedProvider(file file.Service, environmentComponentDirectory string, componentName string) error {
	terraformDirectory := tf.GenerateTerraformIlPath(environmentComponentDirectory, componentName)
	err := file.SaveFileFromString(`
provider "aws" {
	region  = "us-east-1"
	version = "~> 3.0"
	profile = "compuzest-shared"
	alias   = "shared"
}
	`, terraformDirectory, "provider_shared.tf")
	if err != nil {
		return err
	}
	return nil
}

// GenerateProvider save provider file to be executed by terraform
func (tf TerraformGenerator) GenerateProvider(file file.Service, environmentComponentDirectory string, componentName string) error {
	terraformDirectory := tf.GenerateTerraformIlPath(environmentComponentDirectory, componentName)
	err := file.SaveFileFromString(`
provider "aws" {
	region  = "us-east-1"
	version = "~> 3.0"
}
	`, terraformDirectory, "provider.tf")
	if err != nil {
		return err
	}
	return nil
}

// GenerateFromTemplate save terraform backend config
func (tf TerraformGenerator) GenerateTerraformIlPath(environmentComponentDirectory string, environmentComponentName string) string {
	return environmentComponentDirectory + "/" + environmentComponentName + "/terraform"
}

func (tf TerraformGenerator) GenerateFromTemplate(vars interface{}, environmentComponentDirectory string, componentName string, fileUtil file.Service, templatePath string, fileName string) error {
	template, err := template.ParseFiles(templatePath)
	if err != nil {
		return err
	}
	terraformDirectory := tf.GenerateTerraformIlPath(environmentComponentDirectory, componentName)
	err = fileUtil.SaveFileFromTemplate(template, vars, terraformDirectory, fileName+".tf")
	if err != nil {
		return err
	}
	return nil
}
