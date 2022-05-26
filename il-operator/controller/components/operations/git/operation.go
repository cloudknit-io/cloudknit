package git

import (
	"fmt"
	gitapi "github.com/compuzest/zlifecycle-il-operator/controller/common/git"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"github.com/sirupsen/logrus"

	"github.com/compuzest/zlifecycle-il-operator/controller/util"

	gogit "github.com/go-git/go-git/v5"
)

type CleanupFunc func()

func PullOrClone(gitAPI gitapi.API, repoURL string) error {
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

func CloneTemp(gitAPI gitapi.API, repo string, log *logrus.Entry) (dir string, cleanup CleanupFunc, err error) {
	repo = util.RewriteGitHubURLToHTTPS(repo, false)

	log.Infof("Cloning repository %s into a temp folder", repo)

	tempDir, err := createTempDir(repo)
	if err != nil {
		return "", nil, errors.Wrapf(err, "error generating temp dir using system tempdir")
	}
	if tempDir == "" {
		return "", nil, errors.Errorf("invalid tempdir using system tempdir: it is empty")
	}

	if err := gitAPI.Clone(repo, tempDir); err != nil {
		return "", nil, err
	}

	cleanupFunc := func() {
		_ = os.RemoveAll(tempDir)
	}

	return tempDir, cleanupFunc, nil
}

func createTempDir(repo string) (string, error) {
	gitPrefix := "git@"
	httpsPrefix := "https://"

	dirName := strings.Trim(repo, ".git")
	dirName = strings.TrimPrefix(dirName, gitPrefix)
	dirName = strings.TrimPrefix(dirName, httpsPrefix)
	dirName = strings.ReplaceAll(dirName, ":", "-")
	dirName = strings.ReplaceAll(dirName, "/", "-")

	pattern := fmt.Sprintf("repo-%s-", dirName)
	tempDir, err := ioutil.TempDir("", pattern)
	return tempDir, err
}
