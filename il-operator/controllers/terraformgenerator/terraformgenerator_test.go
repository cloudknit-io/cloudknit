package terraformgenerator_test

import (
	"testing"

	v1 "github.com/compuzest/zlifecycle-il-operator/api/v1"
	"github.com/stretchr/testify/assert"

	"github.com/compuzest/zlifecycle-il-operator/controllers/terraformgenerator"
	"github.com/compuzest/zlifecycle-il-operator/mocks"
	"github.com/golang/mock/gomock"
)

func TestGenerateProvider(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tf := terraformgenerator.TerraformGenerator{}

	m := mocks.NewMockService(ctrl)

	m.
		EXPECT().
		SaveFileFromString(gomock.Any(), gomock.Eq("dev-environment-components/component-name/terraform"), gomock.Eq("provider.tf"))
	err := tf.GenerateProvider(m, "dev-environment-components", "component-name")
	assert.NoError(t, err)
}

func TestGenerateBackendTemplate(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tf := terraformgenerator.TerraformGenerator{}
	dummyConfig := terraformgenerator.TerraformBackendConfig{
		Region:        "test-region",
		Profile:       "test-profile",
		Version:       "1.0",
		Bucket:        "test-bucket",
		DynamoDBTable: "test-tflock",
		TeamName:      "test-team",
		EnvName:       "test-env",
		ComponentName: "test-component",
	}

	m := mocks.NewMockService(ctrl)

	m.
		EXPECT().
		SaveFileFromTemplate(gomock.Any(), dummyConfig, gomock.Eq("env-dir/comp-name/terraform"), gomock.Eq("file-name.tf"))

	err := tf.GenerateFromTemplate(dummyConfig, "env-dir", "comp-name", m, "../../templates/terraform_backend.tmpl", "file-name")

	assert.NoError(t, err)
}

func TestGenerateModuleTemplate(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tf := terraformgenerator.TerraformGenerator{}
	testVariables := []*v1.Variable{
		{Name: "foo", Value: "bar"},
		{Name: "bazz", Value: "fun"},
	}
	moduleConfig := terraformgenerator.TerraformModuleConfig{
		ComponentName: "test-env-component",
		Source:        "git@github.com:CompuZest/zlifecycle-il-operator.git",
		Path:          "test/path.txt",
		Variables:     testVariables,
	}

	m := mocks.NewMockService(ctrl)

	m.
		EXPECT().
		SaveFileFromTemplate(gomock.Any(), moduleConfig, gomock.Eq("env-dir/comp-name/terraform"), gomock.Eq("file-name.tf"))

	err := tf.GenerateFromTemplate(moduleConfig, "env-dir", "comp-name", m, "../../templates/terraform_module.tmpl", "file-name")

	assert.NoError(t, err)
}
