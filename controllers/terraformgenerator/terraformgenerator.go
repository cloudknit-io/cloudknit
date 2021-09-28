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
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/compuzest/zlifecycle-il-operator/controllers/secrets"

	stablev1 "github.com/compuzest/zlifecycle-il-operator/api/v1"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/env"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/file"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/il"
)

const (
	DefaultTerraformVersion = "0.13.2"
	DefaultRegion           = "us-east-1"
	DefaultBackendRegion    = "us-east-1"
	DefaultSharedRegion     = "us-east-1"
	DefaultSharedProfile    = "compuzest-shared"
	DefaultSharedAlias      = "shared"
)

func (tf TerraformGenerator) GenerateTerraform(
	fileUtil file.Service,
	vars *TemplateVariables,
	environmentComponentDirectory string,
) error {
	componentName := vars.EnvCompName

	backendConfig := TerraformBackendConfig{
		Region:        DefaultRegion,
		Profile:       "compuzest-shared",
		Version:       DefaultTerraformVersion,
		Bucket:        fmt.Sprintf("zlifecycle-tfstate-%s", env.Config.CompanyName),
		DynamoDBTable: fmt.Sprintf("zlifecycle-tflock-%s", env.Config.CompanyName),
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
		Region:    DefaultSharedRegion,
		Profile:   DefaultSharedProfile,
		Bucket:    fmt.Sprintf("zlifecycle-tfstate-%s", env.Config.CompanyName),
		TeamName:  vars.TeamName,
		EnvName:   vars.EnvName,
		DependsOn: vars.EnvCompDependsOn,
	}

	var assumeRole *AssumeRole
	if vars.EnvCompAWSConfig != nil && vars.EnvCompAWSConfig.AssumeRole != nil {
		assumeRole = &AssumeRole{
			RoleARN:     vars.EnvCompAWSConfig.AssumeRole.RoleARN,
			SessionName: vars.EnvCompAWSConfig.AssumeRole.SessionName,
			ExternalID:  vars.EnvCompAWSConfig.AssumeRole.ExternalID,
		}
	}

	region := DefaultRegion
	if vars.EnvCompAWSConfig != nil && vars.EnvCompAWSConfig.Region != "" {
		region = vars.EnvCompAWSConfig.Region
	}

	providerConfig := TerraformProviderConfig{
		Region:     region,
		AssumeRole: assumeRole,
	}

	sharedProviderConfig := TerraformProviderConfig{
		Region:  DefaultSharedRegion,
		Profile: DefaultSharedProfile,
		Alias:   DefaultSharedAlias,
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

	workingDir, err := os.Getwd()
	if err != nil {
		return err
	}

	if err := tf.GenerateFromTemplate(&providerConfig, environmentComponentDirectory, componentName, fileUtil, getTemplatePath(workingDir, "terraform_provider.tmpl"), "provider.tf"); err != nil {
		return err
	}

	if err := tf.GenerateFromTemplate(&sharedProviderConfig, environmentComponentDirectory, componentName, fileUtil, getTemplatePath(workingDir, "terraform_provider.tmpl"), "provider_shared.tf"); err != nil {
		return err
	}

	if err := tf.GenerateFromTemplate(&backendConfig, environmentComponentDirectory, componentName, fileUtil, getTemplatePath(workingDir, "terraform_backend.tmpl"), "terraform.tf"); err != nil {
		return err
	}

	if err := tf.GenerateFromTemplate(&moduleConfig, environmentComponentDirectory, componentName, fileUtil, getTemplatePath(workingDir, "terraform_module.tmpl"), "module.tf"); err != nil {
		return err
	}

	if len(outputsConfig.Outputs) > 0 {
		if err := tf.GenerateFromTemplate(&outputsConfig, environmentComponentDirectory, componentName, fileUtil, getTemplatePath(workingDir, "terraform_outputs.tmpl"), "outputs.tf"); err != nil {
			return err
		}
	}

	if len(dataConfig.DependsOn) > 0 {
		if err := tf.GenerateFromTemplate(&dataConfig, environmentComponentDirectory, componentName, fileUtil, getTemplatePath(workingDir, "terraform_data.tmpl"), "data.tf"); err != nil {
			return err
		}
	}

	if vars.EnvCompSecrets != nil && len(vars.EnvCompSecrets) > 0 {
		if err := tf.GenerateFromTemplate(&secretsConfig, environmentComponentDirectory, componentName, fileUtil, getTemplatePath(workingDir, "terraform_secrets.tmpl"), "secrets.tf"); err != nil {
			return err
		}
	}

	return nil
}

func createSecretsConfig(secretArray []*stablev1.Secret, meta secrets.SecretMeta) (*TerraformSecretsConfig, error) {
	scopedSecrets := make([]Secret, 0, len(secretArray))
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

func getTemplatePath(rootDir string, tpl string) string {
	return filepath.Join(rootDir, "templates", tpl)
}

// GenerateFromTemplate save terraform backend config
func (tf TerraformGenerator) GenerateTerraformIlPath(environmentComponentDirectory string, environmentComponentName string) string {
	return filepath.Join(environmentComponentDirectory, environmentComponentName, "terraform")
}

func (tf TerraformGenerator) GenerateFromTemplate(vars interface{}, environmentComponentDirectory string, componentName string, fileUtil file.Service, templatePath string, fileName string) error {
	tpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return err
	}
	terraformDirectory := tf.GenerateTerraformIlPath(environmentComponentDirectory, componentName)
	err = fileUtil.SaveFileFromTemplate(tpl, vars, terraformDirectory, fileName)
	if err != nil {
		return err
	}
	return nil
}
