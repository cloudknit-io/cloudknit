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

package tfgen

import (
	"fmt"
	"path/filepath"

	"github.com/compuzest/zlifecycle-il-operator/controller/external/git"
	"github.com/sirupsen/logrus"

	"github.com/compuzest/zlifecycle-il-operator/controller/codegen/file"
	"github.com/compuzest/zlifecycle-il-operator/controller/codegen/il"
	"github.com/compuzest/zlifecycle-il-operator/controller/codegen/secret"
	"github.com/compuzest/zlifecycle-il-operator/controller/codegen/tfgen/tftmpl"
	"github.com/compuzest/zlifecycle-il-operator/controller/util"
	"github.com/pkg/errors"
)

func GenerateCustomTerraform(
	fs file.API,
	gitService git.API,
	sourceRepoURL string,
	sourceTFModulePath string,
	tfDirectory string,
	log *logrus.Entry,
) error {
	dir, cleanup, err := git.CloneTemp(gitService, sourceRepoURL, log)
	if err != nil {
		return errors.Wrapf(err, "error cloning repo [%s]", sourceRepoURL)
	}
	defer cleanup()

	sourceTFModuleAbsPath := filepath.Join(dir, sourceTFModulePath)
	return fs.CopyDirContent(sourceTFModuleAbsPath, tfDirectory)
}

func GenerateTerraform(
	fileService file.API,
	vars *TemplateVariables,
	tfDirectory string,
) error {
	tpl, err := tftmpl.NewTerraformTemplates()
	if err != nil {
		return errors.Wrap(err, "error instantiating terraform templates service")
	}

	// GENERATE TEMPLATES

	backend, err := generateBackend(tpl, vars)
	if err != nil {
		return errors.Wrap(err, "error generating backend from template")
	}

	data, err := generateData(tpl, vars)
	if err != nil {
		return errors.Wrap(err, "error generating data from template")
	}

	module, err := generateModule(tpl, vars)
	if err != nil {
		return errors.Wrap(err, "error generating module from template")
	}

	outputs, err := generateOutputs(tpl, vars)
	if err != nil {
		return errors.Wrap(err, "error generating outputs from template")
	}

	provider, err := generateProviders(tpl, vars)
	if err != nil {
		return errors.Wrap(err, "error generating provider from template")
	}

	sharedProvider, err := generateSharedProviders(tpl, vars)
	if err != nil {
		return errors.Wrap(err, "error generating shared provider from template")
	}

	scrts, err := generateSecrets(tpl, vars)
	if err != nil {
		return errors.Wrap(err, "error generating secrets from template")
	}

	versions, err := generateVersions(tpl, vars)
	if err != nil {
		return errors.Wrap(err, "error generating versions from template")
	}

	// SAVE TEMPLATES

	if err := fileService.SaveFileFromString(backend, tfDirectory, "terraform.tf"); err != nil {
		return err
	}

	if len(vars.EnvCompDependsOn) > 0 {
		if err := fileService.SaveFileFromString(data, tfDirectory, "data.tf"); err != nil {
			return err
		}
	}

	if err := fileService.SaveFileFromString(module, tfDirectory, "module.tf"); err != nil {
		return err
	}

	if len(vars.EnvCompOutputs) > 0 {
		if err := fileService.SaveFileFromString(outputs, tfDirectory, "outputs.tf"); err != nil {
			return err
		}
	}

	if err := fileService.SaveFileFromString(provider, tfDirectory, "provider.tf"); err != nil {
		return err
	}

	if err := fileService.SaveFileFromString(sharedProvider, tfDirectory, "provider_shared.tf"); err != nil {
		return err
	}

	if len(vars.EnvCompSecrets) > 0 {
		if err := fileService.SaveFileFromString(scrts, tfDirectory, "secrets.tf"); err != nil {
			return err
		}
	}

	if err := fileService.SaveFileFromString(versions, tfDirectory, "versions.tf"); err != nil {
		return err
	}

	return nil
}

func generateBackend(tpl *tftmpl.TerraformTemplates, vars *TemplateVariables) (string, error) {
	bucket := vars.AWSStateBucket
	if bucket == "" {
		bucket = defaultStateBucket(vars.ZLEnvironment, vars.Company)
	}
	lockTable := vars.AWSStateLockTable
	if lockTable == "" {
		lockTable = defaultStateLockTable(vars.ZLEnvironment, vars.Company)
	}
	backendConfig := TerraformBackendConfig{
		Region:        vars.AWSSharedRegion,
		Profile:       vars.AWSStateProfile,
		Bucket:        bucket,
		DynamoDBTable: lockTable,
		Key:           defaultStateKey(vars.Team, vars.Environment, vars.EnvironmentComponent),
		Encrypt:       true,
	}

	return tpl.Execute(backendConfig, tftmpl.TmplTFBackend)
}

func defaultStateBucket(zlEnvironment, company string) string {
	return fmt.Sprintf("zlifecycle-%s-tfstate-%s", zlEnvironment, company)
}

func defaultStateLockTable(zlEnvironment, company string) string {
	return fmt.Sprintf("zlifecycle-%s-tflock-%s", zlEnvironment, company)
}

func defaultStateKey(team, environment, component string) string {
	return fmt.Sprintf("%s/%s/%s/terraform.tfstate", team, environment, component)
}

func generateData(tpl *tftmpl.TerraformTemplates, vars *TemplateVariables) (string, error) {
	dataConfig := TerraformDataConfig{
		Region:      vars.AWSRegion,
		Profile:     vars.AWSProfile,
		Bucket:      defaultStateBucket(vars.ZLEnvironment, vars.Company),
		Team:        vars.Team,
		Environment: vars.Environment,
		DependsOn:   vars.EnvCompDependsOn,
	}

	return tpl.Execute(dataConfig, tftmpl.TmplTFData)
}

func generateModule(tpl *tftmpl.TerraformTemplates, vars *TemplateVariables) (string, error) {
	standardizedVariables, err := standardizeVariables(vars.EnvCompVariables)
	if err != nil {
		return "", err
	}

	moduleSource := il.EnvironmentComponentModuleSource(vars.EnvCompModuleSource, vars.EnvCompModuleName)
	moduleConfig := TerraformModuleConfig{
		Component:     vars.EnvironmentComponent,
		Source:        util.RewriteGitHubURLToHTTPS(moduleSource, true),
		Path:          vars.EnvCompModulePath,
		Version:       vars.EnvCompModuleVersion,
		Variables:     standardizedVariables,
		VariablesFile: vars.EnvCompVariablesFile,
		Secrets:       vars.EnvCompSecrets,
	}

	return tpl.Execute(moduleConfig, tftmpl.TmplTFModule)
}

func generateOutputs(tpl *tftmpl.TerraformTemplates, vars *TemplateVariables) (string, error) {
	if len(vars.EnvCompOutputs) == 0 {
		return "", nil
	}
	outputsConfig := TerraformOutputsConfig{
		Component: vars.EnvironmentComponent,
		Outputs:   vars.EnvCompOutputs,
	}

	return tpl.Execute(outputsConfig, tftmpl.TmplTFOutputs)
}

func generateProviders(tpl *tftmpl.TerraformTemplates, vars *TemplateVariables) (string, error) {
	var assumeRole *AssumeRole
	if vars.EnvCompAWSConfig != nil && vars.EnvCompAWSConfig.AssumeRole != nil {
		assumeRole = &AssumeRole{
			RoleARN:     vars.EnvCompAWSConfig.AssumeRole.RoleARN,
			SessionName: vars.EnvCompAWSConfig.AssumeRole.SessionName,
			ExternalID:  vars.EnvCompAWSConfig.AssumeRole.ExternalID,
		}
	}

	region := vars.AWSRegion
	if vars.EnvCompAWSConfig != nil && vars.EnvCompAWSConfig.Region != "" {
		region = vars.EnvCompAWSConfig.Region
	}

	providerConfig := TerraformProviderConfig{
		Region:     region,
		AssumeRole: assumeRole,
	}

	return tpl.Execute(providerConfig, tftmpl.TmplTFProvider)
}

func generateSharedProviders(tpl *tftmpl.TerraformTemplates, vars *TemplateVariables) (string, error) {
	sharedProviderConfig := TerraformProviderConfig{
		Region:  vars.AWSRegion,
		Profile: vars.AWSSharedProfile,
		Alias:   vars.AWSSharedProviderAlias,
	}

	return tpl.Execute(sharedProviderConfig, tftmpl.TmplTFProvider)
}

func generateSecrets(tpl *tftmpl.TerraformTemplates, vars *TemplateVariables) (string, error) {
	secretsMeta := secret.Identifier{
		Company:              vars.Company,
		Team:                 vars.Team,
		Environment:          vars.Environment,
		EnvironmentComponent: vars.EnvironmentComponent,
	}
	secretsConfig, err := createSecretsConfig(vars.EnvCompSecrets, secretsMeta)
	if err != nil {
		return "", err
	}

	return tpl.Execute(secretsConfig, tftmpl.TmplTFSecrets)
}

func generateVersions(tpl *tftmpl.TerraformTemplates, vars *TemplateVariables) (string, error) {
	versionsConfig := TerraformVersionsConfig{
		TerraformVersion: vars.TerraformVersion,
		AWSVersion:       vars.AWSProviderVersion,
	}

	return tpl.Execute(versionsConfig, tftmpl.TmplTFVersions)
}
