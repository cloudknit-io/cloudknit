package terraformgenerator

import stablev1 "github.com/compuzest/zlifecycle-il-operator/api/v1"

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
