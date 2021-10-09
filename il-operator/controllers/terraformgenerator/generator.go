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
	"io"
	"os"
	"path/filepath"
	"text/template"

	"github.com/compuzest/zlifecycle-il-operator/controllers/secrets"

	"github.com/compuzest/zlifecycle-il-operator/controllers/util/env"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/file"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/il"
)

const (
	DefaultTerraformVersion = "0.13.2"
	DefaultRegion           = "us-east-1"
	DefaultSharedRegion     = "us-east-1"
	DefaultSharedProfile    = "compuzest-shared"
	DefaultSharedAlias      = "shared"
)

func GenerateTerraform(
	fileService file.Service,
	vars *TemplateVariables,
	componentDirectory string,
) error {
	componentName := vars.EnvCompName

	backendConfig := TerraformBackendConfig{
		Region:        DefaultSharedRegion,
		Profile:       DefaultSharedProfile,
		Version:       DefaultTerraformVersion,
		Bucket:        fmt.Sprintf("zlifecycle-tfstate-%s", env.Config.CompanyName),
		DynamoDBTable: fmt.Sprintf("zlifecycle-tflock-%s", env.Config.CompanyName),
		TeamName:      vars.TeamName,
		EnvName:       vars.EnvName,
		ComponentName: componentName,
	}

	standardizedVariables, err := standardizeVariables(vars.EnvCompVariables)
	if err != nil {
		return err
	}
	moduleConfig := TerraformModuleConfig{
		ComponentName: componentName,
		Source:        il.EnvComponentModuleSource(vars.EnvCompModuleSource, vars.EnvCompModuleName),
		Path:          vars.EnvCompModulePath,
		Version:       vars.EnvCompModuleVersion,
		Variables:     standardizedVariables,
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

	tfPath := terraformPath{componentDirectory: componentDirectory, componentName: componentName}

	terraformDirectory := il.TerraformIlPath(tfPath.componentDirectory, tfPath.componentName)
	if err := generateFile(fileService, &providerConfig, terraformDirectory, "provider.tf", "terraform_provider"); err != nil {
		return err
	}

	if err := generateFile(fileService, &sharedProviderConfig, terraformDirectory, "provider_shared.tf", "terraform_provider"); err != nil {
		return err
	}

	if err := generateFile(fileService, &backendConfig, terraformDirectory, "terraform.tf", "terraform_backend"); err != nil {
		return err
	}

	if err := generateFile(fileService, &moduleConfig, terraformDirectory, "module.tf", "terraform_module"); err != nil {
		return err
	}

	if len(outputsConfig.Outputs) > 0 {
		if err := generateFile(fileService, &outputsConfig, terraformDirectory, "outputs.tf", "terraform_outputs"); err != nil {
			return err
		}
	}

	if len(dataConfig.DependsOn) > 0 {
		if err := generateFile(fileService, &dataConfig, terraformDirectory, "data.tf", "terraform_data"); err != nil {
			return err
		}
	}

	if vars.EnvCompSecrets != nil && len(vars.EnvCompSecrets) > 0 {
		if err := generateFile(fileService, &secretsConfig, terraformDirectory, "secrets.tf", "terraform_secrets"); err != nil {
			return err
		}
	}

	return nil
}

func generateFile(service file.Service, templateVars interface{}, terraformDirectory string, fileName string, templateName string) error {
	f, err := service.NewFile(terraformDirectory, fileName)
	if err != nil {
		return err
	}
	return generateFromTemplate(templateVars, templateName, f)
}

// generateFromTemplate save terraform backend config
func generateFromTemplate(vars interface{}, templateName string, writer io.Writer) error {
	workingDir, err := os.Getwd()
	if err != nil {
		return err
	}

	templatePath := getTemplatePath(workingDir, templateName)

	tpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return err
	}

	return tpl.Execute(writer, vars)
}

func getTemplatePath(rootDir string, tmpl string) string {
	return filepath.Join(rootDir, "templates", fmt.Sprintf("%s.tmpl", tmpl))
}
