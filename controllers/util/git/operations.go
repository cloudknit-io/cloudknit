package git

import (
	"io/ioutil"
	"os"
)

type CleanupFunc func() error

func CloneTemp(gitAPI Git, repo string) (dir string, cleanup CleanupFunc, err error) {
	tempDir, err := ioutil.TempDir("", "il-")
	if err != nil {
		return "", nil, err
	}

	if err := gitAPI.Clone(repo, tempDir); err != nil {
		return "", nil, err
	}

	cleanupFunc := func() error {
		return os.RemoveAll(tempDir)
	}

	return tempDir, cleanupFunc, nil
}
