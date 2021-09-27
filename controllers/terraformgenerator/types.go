package terraformgenerator

import (
	stablev1 "github.com/compuzest/zlifecycle-il-operator/api/v1"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/file"
	"github.com/go-logr/logr"
)

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
	Variables     []*stablev1.Variable
	VariablesFile string
	Secrets       []*stablev1.Secret
}

// TerraformOutputsConfig for creating tf module outputs
type TerraformOutputsConfig struct {
	ComponentName string
	Outputs       []*stablev1.Output
}

// TerraformSecretsConfig for creating tf secrets
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

// TerraformDataConfig variables for creating tf backend
type TerraformDataConfig struct {
	Region    string
	Bucket    string
	Profile   string
	TeamName  string
	EnvName   string
	DependsOn []string
}

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
	EnvCompAWSConfig     *stablev1.AWS
}