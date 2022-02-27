package git

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/compuzest/zlifecycle-il-operator/controllers/util"

	gogit "github.com/go-git/go-git/v5"
)

type CleanupFunc func()

func PullOrClone(gitAPI API, repoURL string) error {
	tempPath := os.TempDir()
	repoPath := filepath.Join(tempPath, repoURL)
	err := gitAPI.Open(repoURL)
	if errors.Is(err, gogit.ErrRepositoryNotExists) {
		if err := gitAPI.Clone(repoURL, repoPath); err != nil {
			return err
		}
		return nil
	}
	if err != nil {
		return err
	}

	return nil
}

func CloneTemp(gitAPI API, repo string, log *logrus.Entry) (dir string, cleanup CleanupFunc, err error) {
	httpsPrefix := "https://github.com/"
	httpsRepo := util.RewriteGitURLToHTTPS(repo)
	dirName := strings.ReplaceAll(strings.TrimPrefix(strings.TrimSuffix(httpsRepo, ".git"), httpsPrefix), "/", "-")
	log.Infof("Cloning repository %s into a temp folder", repo)
	tempDir, err := ioutil.TempDir("", fmt.Sprintf("repo-%s-", dirName))
	if err != nil {
		return "", nil, err
	}

	if err := gitAPI.Clone(httpsRepo, tempDir); err != nil {
		return "", nil, err
	}

	cleanupFunc := func() {
		_ = os.RemoveAll(tempDir)
	}

	return tempDir, cleanupFunc, nil
}
