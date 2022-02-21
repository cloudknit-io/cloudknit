package terraformgenerator

import (
	stablev1 "github.com/compuzest/zlifecycle-il-operator/api/v1"
	"github.com/compuzest/zlifecycle-il-operator/controllers/codegen/file"
	"github.com/go-logr/logr"
)

// TerraformBackendConfig variables for creating tf backend.
type TerraformBackendConfig struct {
	Region        string
	Key           string
	Bucket        string
	DynamoDBTable string
	Profile       string
	TeamName      string
	EnvName       string
	ComponentName string
}

// TerraformVersionsConfig variables for creating tf versions.
type TerraformVersionsConfig struct {
	Version string
}

// TerraformModuleConfig variables for creating tf module.
type TerraformModuleConfig struct {
	ComponentName string
	Source        string
	Path          string
	Version       string
	Variables     []*Variable
	VariablesFile string
	Secrets       []*stablev1.Secret
}

type Variable struct {
	Name  string
	Value string
}

// TerraformOutputsConfig for creating tf module outputs.
type TerraformOutputsConfig struct {
	ComponentName string
	Outputs       []*stablev1.Output
}

// TerraformSecretsConfig for creating tf secrets.
type TerraformSecretsConfig struct {
	Secrets []Secret
}

type TerraformProviderConfig struct {
	Region     string
	AssumeRole *AssumeRole
	Profile    string
	Alias      string
}

type AssumeRole struct {
	RoleARN     string
	SessionName string
	ExternalID  string
}

type Secret struct {
	Key  string
	Name string
}

// TerraformDataConfig variables for creating tf backend.
type TerraformDataConfig struct {
	Region    string
	Bucket    string
	Profile   string
	TeamName  string
	EnvName   string
	DependsOn []string
}

// UtilTerraformGenerator package interface for generating terraform files.
type UtilTerraformGenerator interface {
	GenerateTerraform(tempILRepoDir string, fileUtil file.API, vars *TemplateVariables, environmentComponentDirectory string) error
	GenerateSharedProvider(file file.API, environmentComponentDirectory string, componentName string) error
	GenerateFromTemplate(
		vars interface{},
		environmentComponentDirectory string,
		componentName string,
		fileUtil file.API,
		templateName string,
		filePath string,
	) error
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
	EnvCompModuleVersion string
	EnvCompOutputs       []*stablev1.Output
	EnvCompDependsOn     []string
	EnvCompAWSConfig     *stablev1.AWS
}
