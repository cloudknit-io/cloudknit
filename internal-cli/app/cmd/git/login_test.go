package git_test

import (
	"github.com/compuzest/zlifecycle-internal-cli/app/cmd/git"
	"github.com/compuzest/zlifecycle-internal-cli/app/env"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestMain(m *testing.M) {
	_ = os.Mkdir(env.TestDir, 0o744)
	exitVal := m.Run()

	_ = os.RemoveAll(env.TestDir)
	os.Exit(exitVal)
}

func TestLoginCmdGitHubAppInternal(t *testing.T) {
	t.Parallel()

	if env.TestMode != env.TestModeIntegration && env.TestMode != env.TestModeAll {
		t.Skip()
	}

	org := "zlifecycle-il"

	env.GitAuth = "github-app-internal"
	env.GitConfigDir = env.TestDir
	env.Verbose = true

	cmd := git.NewLoginCmd()
	cmd.SetArgs([]string{org})

	err := cmd.Execute()
	assert.Nil(t, err)

	gitconfigPath := filepath.Join(env.TestDir, ".gitconfig")
	assert.FileExists(t, gitconfigPath)

	t.Cleanup(func() {
		_ = os.Remove(gitconfigPath)
	})
}
