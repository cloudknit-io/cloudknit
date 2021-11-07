package git

import (
	"io/ioutil"
	"os"
	"strings"
)

type CleanupFunc func() error

func CloneTemp(gitAPI Git, repo string) (dir string, cleanup CleanupFunc, err error) {
	httpsRepo := repo
	if strings.HasPrefix(httpsRepo, "git@github.com:") {
		httpsRepo = strings.ReplaceAll(httpsRepo, "git@github.com:", "https://github.com/")
	}
	tempDir, err := ioutil.TempDir("", "il-")
	if err != nil {
		return "", nil, err
	}

	if err := gitAPI.Clone(httpsRepo, tempDir); err != nil {
		return "", nil, err
	}

	cleanupFunc := func() error {
		return os.RemoveAll(tempDir)
	}

	return tempDir, cleanupFunc, nil
}
