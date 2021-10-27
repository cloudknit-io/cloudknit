package gotfvars

import (
	"fmt"
	"io/ioutil"

	"github.com/compuzest/zlifecycle-il-operator/controllers/util/common"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/github"
	"github.com/go-logr/logr"
)

func GetVariablesFromTfvarsFile(log logr.Logger, api github.RepositoryAPI, repoURL string, ref string, path string) (string, error) {
	log.Info("Downloading tfvars file", "repoUrl", repoURL, "ref", ref, "path", path)
	buff, exists, err := downloadTfvarsFile(api, repoURL, ref, path)
	if err != nil {
		return "", err
	}
	if !exists {
		return "", fmt.Errorf("file does not exist: %s/%s?ref=%s", repoURL, path, ref)
	}
	tfvars := string(buff)

	return tfvars, nil
}

func downloadTfvarsFile(api github.RepositoryAPI, repoURL string, ref string, path string) (file []byte, exists bool, err error) {
	rc, exists, err := github.DownloadFile(api, repoURL, ref, path)
	if err != nil {
		return nil, false, err
	}
	if !exists {
		return nil, false, nil
	}
	defer common.CloseBody(rc)
	buff, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, false, err
	}
	return buff, true,nil
}
