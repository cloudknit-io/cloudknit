package gotfvars

import (
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/common"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/github"
	"github.com/go-logr/logr"
	"io/ioutil"
)

func GetVariablesFromTfvarsFile(log logr.Logger, api github.RepositoryApi, repoUrl string, ref string, path string) (string, error) {
	log.Info("Downloading tfvars file", "repoUrl", repoUrl, "ref", ref, "path", path)
	buff, err := downloadTfvarsFile(api, repoUrl, ref, path)
	if err != nil {
		return "", err
	}
	tfvars := string(buff)

	return tfvars, nil
}

func downloadTfvarsFile(api github.RepositoryApi, repoUrl string, ref string, path string) ([]byte, error) {
	rc, err := github.DownloadFile(api, repoUrl, ref, path)
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
