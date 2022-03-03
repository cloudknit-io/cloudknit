package git

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

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
	dirName := strings.ReplaceAll(strings.TrimPrefix(strings.Trim(httpsRepo, ".git"), httpsPrefix), "/", "-")
	log.Infof("Cloning repository %s into a temp folder", httpsRepo)
	pattern := fmt.Sprintf("repo-%s-", dirName)
	tempDir, err := ioutil.TempDir("", pattern)
	if err != nil {
		return "", nil, errors.Wrapf(err, "error generating temp dir using system tempdir and pattern %s", pattern)
	}
	if tempDir == "" {
		return "", nil, errors.Errorf("invalid tempdir using system tempdir and pattern %s: tempdir is empty", pattern)
	}

	if err := gitAPI.Clone(httpsRepo, tempDir); err != nil {
		return "", nil, err
	}

	cleanupFunc := func() {
		_ = os.RemoveAll(tempDir)
	}

	return tempDir, cleanupFunc, nil
}

func CloneTempSSH(gitAPI API, repo string, log *logrus.Entry) (dir string, cleanup CleanupFunc, err error) {
	log.Infof("Cloning repository %s into a temp folder", repo)
	tempDir, err := ioutil.TempDir("", "il-")
	if err != nil {
		return "", nil, errors.Wrapf(err, "error generating temp dir using system tempdir")
	}
	if tempDir == "" {
		return "", nil, errors.Errorf("invalid tempdir using system tempdir. tempdir is empty")
	}

	if err := gitAPI.Clone(repo, tempDir); err != nil {
		return "", nil, err
	}

	cleanupFunc := func() {
		_ = os.RemoveAll(tempDir)
	}

	return tempDir, cleanupFunc, nil
}
