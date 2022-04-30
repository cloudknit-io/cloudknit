package tfgen

import (
	stablev1 "github.com/compuzest/zlifecycle-il-operator/api/v1"
	"github.com/compuzest/zlifecycle-il-operator/controllers/env"
	"github.com/compuzest/zlifecycle-il-operator/controllers/external/secrets"
)

// TerraformBackendConfig variables for creating tf backend.
type TerraformBackendConfig struct {
	Region        string
	Key           string
	Bucket        string
	DynamoDBTable string
	Profile       string
	Encrypt       bool
}

// TerraformVersionsConfig variables for creating tf versions.
type TerraformVersionsConfig struct {
	TerraformVersion string
	AWSVersion       string
}

// TerraformModuleConfig variables for creating tf module.
type TerraformModuleConfig struct {
	Component     string
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
	Component string
	Outputs   []*stablev1.Output
}

// TerraformSecretsConfig for creating tf secrets.
type TerraformSecretsConfig struct {
	Secrets []*Secret
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
	Region      string
	Bucket      string
	Profile     string
	Team        string
	Environment string
	DependsOn   []string
}

type TemplateVariables struct {
	ZLEnvironment          string
	Company                string
	Team                   string
	Environment            string
	EnvironmentComponent   string
	EnvCompVariables       []*stablev1.Variable
	EnvCompVariablesFile   string
	EnvCompSecrets         []*stablev1.Secret
	EnvCompModuleSource    string
	EnvCompModulePath      string
	EnvCompModuleName      string
	EnvCompModuleVersion   string
	EnvCompOutputs         []*stablev1.Output
	EnvCompDependsOn       []string
	EnvCompAWSConfig       *stablev1.AWS
	TerraformVersion       string
	AWSProviderVersion     string
	AWSRegion              string
	AWSSharedRegion        string
	AWSSharedProviderAlias string
	AWSSharedProfile       string
	AWSProfile             string
	AWSStateBucket         string
	AWSStateLockTable      string
	AWSProviderAlias       string
	AWSStateKey            string
	AWSStateProfile        string
}

func NewTemplateVariablesFromEnvironment(
	e *stablev1.Environment,
	ec *stablev1.EnvironmentComponent,
	tfvars string,
	tfcfg *secrets.TerraformStateConfig,
) *TemplateVariables {
	vars := &TemplateVariables{
		ZLEnvironment:          env.Config.ZLEnvironment,
		Company:                env.Config.CompanyName,
		Team:                   e.Spec.TeamName,
		Environment:            e.Spec.EnvName,
		EnvironmentComponent:   ec.Name,
		EnvCompVariables:       ec.Variables,
		EnvCompVariablesFile:   tfvars,
		EnvCompSecrets:         ec.Secrets,
		EnvCompModuleSource:    ec.Module.Source,
		EnvCompModulePath:      ec.Module.Path,
		EnvCompModuleName:      ec.Module.Name,
		EnvCompModuleVersion:   ec.Module.Version,
		EnvCompOutputs:         ec.Outputs,
		EnvCompDependsOn:       ec.DependsOn,
		EnvCompAWSConfig:       ec.AWS,
		TerraformVersion:       env.Config.TerraformDefaultVersion,
		AWSProviderVersion:     env.Config.TerraformDefaultAWSProviderVersion,
		AWSRegion:              env.Config.AWSRegion,
		AWSSharedRegion:        env.Config.TerraformDefaultSharedAWSRegion,
		AWSSharedProviderAlias: env.Config.TerraformDefaultSharedAWSAlias,
		AWSSharedProfile:       env.Config.TerraformDefaultSharedAWSProfile,
		AWSProfile:             env.Config.TerraformDefaultSharedAWSProfile,
		AWSStateProfile:        env.Config.TerraformDefaultSharedAWSProfile,
	}

	if tfcfg != nil {
		vars.AWSStateBucket = tfcfg.Bucket
		vars.AWSStateLockTable = tfcfg.LockTable
		vars.AWSStateProfile = env.Config.TerraformCustomerStateAWSProfile
	}

	return vars
}
