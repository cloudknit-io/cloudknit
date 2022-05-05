package aws_test

import (
	"fmt"
	"github.com/compuzest/zlifecycle-internal-cli/app/util"
	"os"
	"path/filepath"
	"testing"

	"github.com/compuzest/zlifecycle-internal-cli/app/cmd/aws"
	"github.com/compuzest/zlifecycle-internal-cli/app/env"
	"github.com/stretchr/testify/assert"
)

var td = util.NewTestDirName()

func TestMain(m *testing.M) {
	_ = os.Mkdir(td, 0o744)
	exitVal := m.Run()

	_ = os.RemoveAll(td)
	os.Exit(exitVal)
}

func TestConfigureCmd(t *testing.T) {
	t.Parallel()

	if env.TestMode != env.TestModeIntegration && env.TestMode != env.TestModeAll {
		t.Skip()
	}

	env.AWSConfigFile = filepath.Join(td, "credentials")
	env.AWSProfile = "compuzest-bootstrap"
	env.AWSGeneratedProfile = "customer-state"
	env.Verbose = true

	cmd := aws.NewConfigureCmd()
	cmd.SetArgs([]string{})
	err := cmd.Flags().Set("auth-mode", "profile")
	assert.Nil(t, err)
	err = cmd.Flags().Set("profile", "compuzest-shared")
	assert.Nil(t, err)
	err = cmd.Flags().Set("generated-profile", "customer-state")
	assert.Nil(t, err)
	err = cmd.Flags().Set("company", "zbank")
	assert.Nil(t, err)
	err = cmd.Flags().Set("team", "checkout")
	assert.Nil(t, err)

	err = cmd.Execute()
	assert.Nil(t, err)

	b, err := os.ReadFile(env.AWSConfigFile)
	assert.Nil(t, err)

	expected := fmt.Sprintf("[%s]", env.AWSGeneratedProfile)
	assert.Contains(t, string(b), expected)
}
