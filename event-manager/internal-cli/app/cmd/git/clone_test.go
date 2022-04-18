package git_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/compuzest/zlifecycle-internal-cli/app/cmd/git"
	"github.com/compuzest/zlifecycle-internal-cli/app/env"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCloneCmdGitHubAppInternal(t *testing.T) {
	t.Parallel()

	if env.TestMode != "integration" {
		t.Skip()
	}

	repo := "https://github.com/zlifecycle-il/app-zmart-zl-il.git"
	dir := filepath.Join(env.TestDir, uuid.New().String())

	env.GitAuth = "github-app-internal"
	env.Verbose = true

	cmd := git.NewCloneCmd()
	cmd.SetArgs([]string{repo})

	err := cmd.Flags().Set("dir", dir)
	assert.Nil(t, err)
	err = cmd.Execute()
	assert.Nil(t, err)

	assert.DirExists(t, dir)

	t.Cleanup(func() {
		_ = os.RemoveAll(dir)
	})
}
