package terraform_test

import (
	"os"
	"testing"

	"github.com/compuzest/zlifecycle-il-operator/controller/common/git"

	"github.com/sirupsen/logrus"

	"github.com/compuzest/zlifecycle-il-operator/controller/codegen/file"
	"github.com/compuzest/zlifecycle-il-operator/controller/codegen/terraform"
	"github.com/compuzest/zlifecycle-il-operator/controller/env"
	"github.com/compuzest/zlifecycle-il-operator/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGenerateCustomTerraform(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockFileService := file.NewMockAPI(mockCtrl)
	mockGitClient := git.NewMockAPI(mockCtrl)

	testRepo := "git@github.com:test/foo.git"
	testPath := "/some/module"
	testTFDirectory := "/tmp/some/dir"

	mockGitClient.EXPECT().Clone(testRepo, gomock.Any()).Return(nil)
	mockFileService.EXPECT().CopyDirContent(gomock.Any(), testTFDirectory, true).Return(nil)

	log := logrus.NewEntry(logrus.New())
	err := terraform.GenerateCustomTerraform(mockFileService, mockGitClient, testRepo, testPath, testTFDirectory, log)
	assert.Nil(t, err)
}

func TestGenerateTerraform(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockFileService := file.NewMockAPI(mockCtrl)

	expectedBackend := `terraform {
    backend "s3" {
        region         = "us-east-1"
        bucket         = "zlifecycle-dev-tfstate-zbank"
        encrypt        = "true"
        key            = "design/dev/networking/terraform.tfstate"
        profile        = "compuzest-shared"
        dynamodb_table = "zlifecycle-dev-tflock-zbank"
    }
}
`
	first := mockFileService.EXPECT().SaveFileFromString(expectedBackend, gomock.Any(), "terraform.tf")

	expectedModule := `module "networking" {
	source  = "git@github.com:terraform-aws-modules/terraform-aws-vpc.git"
    foo = "bar"
    baz = "test"
}
`
	second := mockFileService.EXPECT().SaveFileFromString(expectedModule, gomock.Any(), "module.tf").After(first)

	expectedProvider := `provider "aws" {
	region  = "us-east-1"
}
`
	third := mockFileService.EXPECT().SaveFileFromString(expectedProvider, gomock.Any(), "provider.tf").After(second)

	expectedSharedProvider := `provider "aws" {
	region  = "us-east-1"
	profile = "compuzest-shared"
	alias   = "shared"
}
`
	fourth := mockFileService.EXPECT().SaveFileFromString(expectedSharedProvider, gomock.Any(), "provider_shared.tf").After(third)

	expectedVersions := `terraform {
	required_version = "1.0.9"

	required_providers {
		aws = {
			version = "~> 4.0"
		}
	}
}
`
	mockFileService.EXPECT().SaveFileFromString(expectedVersions, gomock.Any(), "versions.tf").After(fourth)

	env.Config.ZLEnvironment = "dev"
	env.Config.CompanyName = "zbank"

	mockEnv := mocks.GetMockEnv1(false)
	mockTFVars := `foo = "bar"
baz = "test"`

	tfDirectory := os.TempDir()
	vars := terraform.NewTemplateVariablesFromEnvironment(&mockEnv, mockEnv.Spec.Components[0], mockTFVars, nil)

	err := terraform.GenerateTerraform(mockFileService, vars, tfDirectory)
	assert.Nil(t, err)
}
