package filereconciler

import (
	"io/ioutil"

	"github.com/compuzest/zlifecycle-il-operator/controllers/util/common"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/github"
)

func downloadFile(api github.RepositoryAPI, repoURL string, ref string, path string) (fileBytes []byte, exists bool, err error) {
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
	return buff, true, nil
}
