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

	"github.com/compuzest/zlifecycle-il-operator/controllers/util"

	"github.com/compuzest/zlifecycle-il-operator/controllers/codegen/file"
	"github.com/compuzest/zlifecycle-il-operator/controllers/codegen/il"
	"github.com/compuzest/zlifecycle-il-operator/controllers/codegen/secrets"

	"github.com/compuzest/zlifecycle-il-operator/controllers/env"
)

func GenerateTerraform(
	fileAPI file.API,
	vars *TemplateVariables,
	tfDirectory string,
) error {
	componentName := vars.EnvCompName

	backendConfig := TerraformBackendConfig{
		Region:        env.Config.TerraformDefaultSharedAWSRegion,
		Profile:       env.Config.TerraformDefaultSharedAWSProfile,
		Bucket:        fmt.Sprintf("zlifecycle-%s-tfstate-%s", env.Config.Environment, env.Config.CompanyName),
		DynamoDBTable: fmt.Sprintf("zlifecycle-%s-tflock-%s", env.Config.Environment, env.Config.CompanyName),
		TeamName:      vars.TeamName,
		EnvName:       vars.EnvName,
		ComponentName: componentName,
	}

	versionsConfig := TerraformVersionsConfig{
		Version: env.Config.TerraformDefaultVersion,
	}

	standardizedVariables, err := standardizeVariables(vars.EnvCompVariables)
	if err != nil {
		return err
	}

	moduleSource := il.EnvironmentComponentModuleSource(vars.EnvCompModuleSource, vars.EnvCompModuleName)
	moduleConfig := TerraformModuleConfig{
		ComponentName: componentName,
		Source:        util.RewriteGitHubURLToHTTPS(moduleSource, true),
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
		Region:    env.Config.TerraformDefaultSharedAWSRegion,
		Profile:   env.Config.TerraformDefaultSharedAWSProfile,
		Bucket:    fmt.Sprintf("zlifecycle-%s-tfstate-%s", env.Config.Environment, env.Config.CompanyName),
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

	region := env.Config.TerraformDefaultAWSRegion
	if vars.EnvCompAWSConfig != nil && vars.EnvCompAWSConfig.Region != "" {
		region = vars.EnvCompAWSConfig.Region
	}

	providerConfig := TerraformProviderConfig{
		Region:     region,
		AssumeRole: assumeRole,
	}

	sharedProviderConfig := TerraformProviderConfig{
		Region:  env.Config.TerraformDefaultSharedAWSRegion,
		Profile: env.Config.TerraformDefaultSharedAWSProfile,
		Alias:   env.Config.TerraformDefaultSharedAWSAlias,
	}

	secretsMeta := secrets.Meta{
		Company:              env.Config.CompanyName,
		Team:                 vars.TeamName,
		Environment:          vars.EnvName,
		EnvironmentComponent: vars.EnvCompName,
	}
	secretsConfig, err := createSecretsConfig(vars.EnvCompSecrets, secretsMeta)
	if err != nil {
		return err
	}

	if err := generateFile(fileAPI, &providerConfig, tfDirectory, "provider.tf", "terraform_provider"); err != nil {
		return err
	}

	if err := generateFile(fileAPI, &sharedProviderConfig, tfDirectory, "provider_shared.tf", "terraform_provider"); err != nil {
		return err
	}

	if err := generateFile(fileAPI, &backendConfig, tfDirectory, "terraform.tf", "terraform_backend"); err != nil {
		return err
	}

	if err := generateFile(fileAPI, &versionsConfig, tfDirectory, "versions.tf", "terraform_versions"); err != nil {
		return err
	}

	if err := generateFile(fileAPI, &moduleConfig, tfDirectory, "module.tf", "terraform_module"); err != nil {
		return err
	}

	if len(outputsConfig.Outputs) > 0 {
		if err := generateFile(fileAPI, &outputsConfig, tfDirectory, "outputs.tf", "terraform_outputs"); err != nil {
			return err
		}
	}

	if len(dataConfig.DependsOn) > 0 {
		if err := generateFile(fileAPI, &dataConfig, tfDirectory, "data.tf", "terraform_data"); err != nil {
			return err
		}
	}

	if vars.EnvCompSecrets != nil && len(vars.EnvCompSecrets) > 0 {
		if err := generateFile(fileAPI, &secretsConfig, tfDirectory, "secrets.tf", "terraform_secrets"); err != nil {
			return err
		}
	}

	return nil
}

func generateFile(service file.API, templateVars interface{}, terraformDirectory string, fileName string, templateName string) error {
	f, err := service.NewFile(terraformDirectory, fileName)
	if err != nil {
		return err
	}
	return generateFromTemplate(templateVars, templateName, f)
}

// generateFromTemplate save terraform backend config.
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
