package tftmpl_test

import (
	"strings"
	"testing"

	v1 "github.com/compuzest/zlifecycle-il-operator/api/v1"

	"github.com/compuzest/zlifecycle-il-operator/controllers/codegen/tfgen"

	"github.com/compuzest/zlifecycle-il-operator/controllers/codegen/tfgen/tftmpl"
	"github.com/stretchr/testify/assert"
)

func TestNewTerraformTemplates(t *testing.T) {
	t.Parallel()

	_, err := tftmpl.NewTerraformTemplates()
	assert.Nil(t, err)
}

func TestTerraformTemplates_Templates(t *testing.T) {
	t.Parallel()

	tpl, err := tftmpl.NewTerraformTemplates()
	assert.Nil(t, err)

	templates := tpl.Templates()
	assert.Len(t, templates, 7)

	expectedTemplates := []tftmpl.TemplateName{
		tftmpl.TmplTFVersions,
		tftmpl.TmplTFBackend,
		tftmpl.TmplTFModule,
		tftmpl.TmplTFData,
		tftmpl.TmplTFOutputs,
		tftmpl.TmplTFProvider,
		tftmpl.TmplTFSecrets,
	}

	assert.ElementsMatch(t, templates, expectedTemplates)
}

func TestTerraformTemplates_ExecuteVersions(t *testing.T) {
	t.Parallel()

	tpl, err := tftmpl.NewTerraformTemplates()
	assert.Nil(t, err)

	vars := tfgen.TerraformVersionsConfig{
		TerraformVersion: "1.2.3",
		AWSVersion:       "4.1",
	}
	output, err := tpl.Execute(vars, tftmpl.TmplTFVersions)
	assert.Nil(t, err)

	expected := `
terraform {
	required_version = "1.2.3"

	required_providers {
		aws = {
			version = "~> 4.1"
		}
	}
}
`
	f1 := strings.Fields(output)
	f2 := strings.Fields(expected)
	assert.ElementsMatch(t, f1, f2)
}

func TestTerraformTemplates_ExecuteBackend(t *testing.T) {
	t.Parallel()

	tpl, err := tftmpl.NewTerraformTemplates()
	assert.Nil(t, err)

	vars := tfgen.TerraformBackendConfig{
		Region:        "us-east-1",
		Key:           "some/test/key.tfstate",
		Bucket:        "test-state",
		DynamoDBTable: "test-lt",
		Profile:       "zlc",
		Encrypt:       true,
	}
	output, err := tpl.Execute(vars, tftmpl.TmplTFBackend)
	assert.Nil(t, err)

	expected := `
terraform {
    backend "s3" {
        region         = "us-east-1"
        bucket         = "test-state"
        encrypt        = "true"
        key            = "some/test/key.tfstate"
        profile        = "zlc"
        dynamodb_table = "test-lt"
    }
}
`
	f1 := strings.Fields(output)
	f2 := strings.Fields(expected)
	assert.ElementsMatch(t, f1, f2)
}

func TestTerraformTemplates_ExecuteData(t *testing.T) {
	t.Parallel()

	tpl, err := tftmpl.NewTerraformTemplates()
	assert.Nil(t, err)

	vars := tfgen.TerraformDataConfig{
		Region:      "us-east-2",
		Bucket:      "tfstate",
		Profile:     "test-profile",
		Team:        "checkout",
		Environment: "dev",
		DependsOn:   []string{"payment"},
	}
	output, err := tpl.Execute(vars, tftmpl.TmplTFData)
	assert.Nil(t, err)

	expected := `data "terraform_remote_state" "payment" {
  backend = "s3"
  config = {
    region  = "us-east-2"
    profile = "test-profile"
    bucket  = "tfstate"
    key     = "checkout/dev/payment/terraform.tfstate"
  }
}`
	f1 := strings.Fields(output)
	f2 := strings.Fields(expected)
	assert.ElementsMatch(t, f1, f2)
}

func TestTerraformTemplates_ExecuteModule(t *testing.T) {
	t.Parallel()

	tpl, err := tftmpl.NewTerraformTemplates()
	assert.Nil(t, err)

	vars := tfgen.TerraformModuleConfig{
		Component: "k8s",
		Source:    "git@github.com:CompuZest/test",
		Path:      "modules/k8s",
		Variables: []*tfgen.Variable{{
			Name:  "test",
			Value: "yes",
		}},
		VariablesFile: `foo = "bar"
baz = 123
ok = true`,
		Secrets: []*v1.Secret{{
			Key:  "/app/database/admin_password",
			Name: "db_password",
		}},
	}
	output, err := tpl.Execute(vars, tftmpl.TmplTFModule)
	assert.Nil(t, err)

	expected := `module "k8s" {
	source = "git@github.com:CompuZest/test//modules/k8s"
	test = yes
	db_password = data.aws_ssm_parameter.db_password.value
	foo = "bar"
	baz = 123
	ok = true
}`
	f1 := strings.Fields(output)
	f2 := strings.Fields(expected)
	assert.ElementsMatch(t, f1, f2)
}

func TestTerraformTemplates_ExecuteOutputs(t *testing.T) {
	t.Parallel()

	tpl, err := tftmpl.NewTerraformTemplates()
	assert.Nil(t, err)

	vars := tfgen.TerraformOutputsConfig{
		Component: "k8s",
		Outputs: []*v1.Output{{
			Name:      "cluster_name",
			Sensitive: true,
		}},
	}
	output, err := tpl.Execute(vars, tftmpl.TmplTFOutputs)
	assert.Nil(t, err)

	expected := `output "cluster_name" {
  value = module.k8s.cluster_name
  sensitive = true
}`
	f1 := strings.Fields(output)
	f2 := strings.Fields(expected)
	assert.ElementsMatch(t, f1, f2)
}

func TestTerraformTemplates_ExecuteProviderNoAssumeRole(t *testing.T) {
	t.Parallel()

	tpl, err := tftmpl.NewTerraformTemplates()
	assert.Nil(t, err)

	vars := tfgen.TerraformProviderConfig{
		Region:     "eu-west-1",
		AssumeRole: nil,
		Profile:    "test-profile",
		Alias:      "shared2",
	}
	output, err := tpl.Execute(vars, tftmpl.TmplTFProvider)
	assert.Nil(t, err)

	expected := `provider "aws" {
  region  = "eu-west-1"
	profile = "test-profile"
	alias   = "shared2"
}`
	f1 := strings.Fields(output)
	f2 := strings.Fields(expected)
	assert.ElementsMatch(t, f1, f2)
}

func TestTerraformTemplates_ExecuteProviderAssumeRole(t *testing.T) {
	t.Parallel()

	tpl, err := tftmpl.NewTerraformTemplates()
	assert.Nil(t, err)

	vars := tfgen.TerraformProviderConfig{
		Region: "eu-west-1",
		AssumeRole: &tfgen.AssumeRole{
			RoleARN:     "some-role-arn",
			SessionName: "test-session",
			ExternalID:  "ext-1234",
		},
		Profile: "test-profile",
		Alias:   "shared2",
	}
	output, err := tpl.Execute(vars, tftmpl.TmplTFProvider)
	assert.Nil(t, err)

	expected := `provider "aws" {
  region  = "eu-west-1"
	assume_role {
		role_arn     = "some-role-arn"
		session_name = "test-session"
		external_id  = "ext-1234"
	}
	profile = "test-profile"
	alias   = "shared2"
}`
	f1 := strings.Fields(output)
	f2 := strings.Fields(expected)
	assert.ElementsMatch(t, f1, f2)
}

func TestTerraformTemplates_ExecuteSecrets(t *testing.T) {
	t.Parallel()

	tpl, err := tftmpl.NewTerraformTemplates()
	assert.Nil(t, err)

	vars := tfgen.TerraformSecretsConfig{
		Secrets: []*tfgen.Secret{{
			Key:  "/app/database/admin_password",
			Name: "db_password",
		}},
	}
	output, err := tpl.Execute(vars, tftmpl.TmplTFSecrets)
	assert.Nil(t, err)

	expected := `data "aws_ssm_parameter" "db_password" {
		name     = "/app/database/admin_password"
		provider = aws.shared
}
`
	f1 := strings.Fields(output)
	f2 := strings.Fields(expected)
	assert.ElementsMatch(t, f1, f2)
}
