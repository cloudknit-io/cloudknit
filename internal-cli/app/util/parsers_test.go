package util_test

import (
	"testing"

	"github.com/compuzest/zlifecycle-internal-cli/app/env"

	"github.com/compuzest/zlifecycle-internal-cli/app/util"

	"github.com/stretchr/testify/assert"
)

func TestParseRepositoryInfo(t *testing.T) {
	t.Parallel()

	if env.TestMode != env.TestModeUnit && env.TestMode != env.TestModeAll {
		t.Skip()
	}

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

func TestReverseParseGitURL(t *testing.T) {
	t.Parallel()

	if env.TestMode != env.TestModeUnit && env.TestMode != env.TestModeAll {
		t.Skip()
	}

	url1 := "https://github.com/CompuZest/zlifecycle-il-operator"

	owner, repo, err := util.ReverseParseGitURL(url1)
	assert.NoError(t, err)
	assert.Equal(t, owner, "CompuZest")
	assert.Equal(t, repo, "zlifecycle-il-operator")

	url2 := "https://github.com/CompuZest/zlifecycle-il-operator.git"
	owner, repo, err = util.ReverseParseGitURL(url2)
	assert.NoError(t, err)
	assert.Equal(t, owner, "CompuZest")
	assert.Equal(t, repo, "zlifecycle-il-operator")

	url3 := "git@github.com:CompuZest/zlifecycle-il-operator"
	owner, repo, err = util.ReverseParseGitURL(url3)
	assert.NoError(t, err)
	assert.Equal(t, owner, "CompuZest")
	assert.Equal(t, repo, "zlifecycle-il-operator")

	url4 := "git@github.com:CompuZest/zlifecycle-il-operator.git"
	owner, repo, err = util.ReverseParseGitURL(url4)
	assert.NoError(t, err)
	assert.Equal(t, owner, "CompuZest")
	assert.Equal(t, repo, "zlifecycle-il-operator")

	owner, repo, err = util.ReverseParseGitURL("")
	assert.Empty(t, owner)
	assert.Empty(t, repo)
	assert.Error(t, err)
}

func TestRewriteGitURLToHTTPS(t *testing.T) {
	t.Parallel()

	if env.TestMode != env.TestModeUnit && env.TestMode != env.TestModeAll {
		t.Skip()
	}

	testRepo1 := "git@github.com:test/test"
	expected1 := "https://github.com/test/test"
	transformed1 := util.RewriteGitURLToHTTPS(testRepo1)
	assert.Equal(t, transformed1, expected1)

	testRepo2 := "https://github.com/hello/world"
	transformed2 := util.RewriteGitURLToHTTPS(testRepo2)
	assert.Equal(t, transformed2, testRepo2)

	testRepo3 := "https://github.com/CompuZest/leet"
	expected3 := "git@github.com:CompuZest/leet"
	transformed3 := util.RewriteGitURLToSSH(testRepo3)
	assert.Equal(t, transformed3, expected3)

	testRepo4 := "git@github.com:CompuZest/rocks"
	transformed4 := util.RewriteGitURLToSSH(testRepo4)
	assert.Equal(t, transformed4, testRepo4)

	testRepo5 := "git@gitlab.com:CompuZest/rocks"
	transformed5 := util.RewriteGitURLToSSH(testRepo5)
	assert.Equal(t, transformed5, testRepo5)

	testRepo6 := "https://gitlab.com/CompuZest/leet"
	expected6 := "git@gitlab.com:CompuZest/leet"
	transformed6 := util.RewriteGitURLToSSH(testRepo6)
	assert.Equal(t, transformed6, expected6)
}
