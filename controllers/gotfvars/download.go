package gotfvars

import (
	"io/ioutil"

	"github.com/compuzest/zlifecycle-il-operator/controllers/util/common"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/github"
	"github.com/go-logr/logr"
)

func GetVariablesFromTfvarsFile(log logr.Logger, api github.RepositoryAPI, repoURL string, ref string, path string) (string, error) {
	log.Info("Downloading tfvars file", "repoUrl", repoURL, "ref", ref, "path", path)
	buff, err := downloadTfvarsFile(api, repoURL, ref, path)
	if err != nil {
		return "", err
	}
	tfvars := string(buff)

	return tfvars, nil
}

func downloadTfvarsFile(api github.RepositoryAPI, repoURL string, ref string, path string) ([]byte, error) {
	rc, err := github.DownloadFile(api, repoURL, ref, path)
	if err != nil {
		return nil, err
	}
	defer common.CloseBody(rc)
	buff, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, err
	}
	return buff, nil
}
