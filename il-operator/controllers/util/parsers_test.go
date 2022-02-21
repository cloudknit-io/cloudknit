package util_test

import (
	"testing"

	"github.com/compuzest/zlifecycle-il-operator/controllers/util"

	"github.com/stretchr/testify/assert"
)

func TestParseRepositoryInfo(t *testing.T) {
	t.Parallel()

	url1 := "https://github.com/CompuZest/zlifecycle-il-operator"

	owner, repo, err := util.ParseRepositoryInfo(url1)
	assert.NoError(t, err)
	assert.Equal(t, owner, "CompuZest")
	assert.Equal(t, repo, "zlifecycle-il-operator")

	url2 := "https://github.com/CompuZest/zlifecycle-il-operator.git"
	owner, repo, err = util.ParseRepositoryInfo(url2)
	assert.NoError(t, err)
	assert.Equal(t, owner, "CompuZest")
	assert.Equal(t, repo, "zlifecycle-il-operator")

	url3 := "git@github.com:CompuZest/zlifecycle-il-operator"
	owner, repo, err = util.ParseRepositoryInfo(url3)
	assert.NoError(t, err)
	assert.Equal(t, owner, "CompuZest")
	assert.Equal(t, repo, "zlifecycle-il-operator")

	url4 := "git@github.com:CompuZest/zlifecycle-il-operator.git"
	owner, repo, err = util.ParseRepositoryInfo(url4)
	assert.NoError(t, err)
	assert.Equal(t, owner, "CompuZest")
	assert.Equal(t, repo, "zlifecycle-il-operator")

	owner, repo, err = util.ParseRepositoryInfo("")
	assert.Empty(t, owner)
	assert.Empty(t, repo)
	assert.Error(t, err)
}
