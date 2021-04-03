package terraformgenerator_test

import (
	"testing"

	terraformgenerator "github.com/compuzest/zlifecycle-il-operator/controllers/terraformgenerator"
	"github.com/compuzest/zlifecycle-il-operator/mocks"
	"github.com/golang/mock/gomock"
)

func TestGenerateProvider(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tf := terraformgenerator.TerraformGenerator{}

	m := mocks.NewMockUtilFile(ctrl)

	m.
		EXPECT().
		SaveFileFromString(gomock.Any(), gomock.Eq("dev-environment-components/terraform"), gomock.Eq("provider.tf"))
	tf.GenerateProvider(m, "dev-environment-components")
}

func TestGenerateTemplate(t *testing.T) {
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

	m := mocks.NewMockUtilFile(ctrl)

	m.
		EXPECT().
		SaveFileFromTemplate(gomock.Any(), dummyConfig, gomock.Eq("env-dir/terraform"), gomock.Eq("file-name.tf"))

	tf.GenerateFromTemplate(dummyConfig, "env-dir", m, "../../templates/terraform_backend.tmpl", "file-name")
}
