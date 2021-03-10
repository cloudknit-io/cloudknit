package argocd

import (
	"bytes"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	ctrl "sigs.k8s.io/controller-runtime"
	"testing"
)

func TestRegisterRepoNewRepo(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockArgocdAPI := NewMockArgocdAPI(mockCtrl)

	mockArgocdAPI.EXPECT().GetAuthToken(gomock.Any()).Return(&GetTokenResponse{Token: "test_token"}, nil)
	repo := Repository{Repo: "git@github.com:CompuZest/test_repo.git", Name: "test_repo"}
	list := RepositoryList{Items: []Repository{repo}}
	r1    := ioutil.NopCloser(bytes.NewReader([]byte{}))
	mockArgocdAPI.EXPECT().ListRepositories(gomock.Any(), gomock.Any()).Return(&list, &http.Response{Body: r1}, nil)
	r2    := ioutil.NopCloser(bytes.NewReader([]byte{}))
	mockArgocdAPI.EXPECT().CreateRepository(gomock.Any(), gomock.Any(), gomock.Any()).Return(&http.Response{Body: r2}, nil)

	log := ctrl.Log.WithName("TestRegisterRepoNewRepo")
	registered, err := RegisterRepo(log, mockArgocdAPI, RepoOpts{})
	assert.True(t, registered)
	assert.NoError(t, err)
}

func TestRegisterRepoExistingRepo(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockArgocdAPI := NewMockArgocdAPI(mockCtrl)

	mockArgocdAPI.EXPECT().GetAuthToken(gomock.Any()).Return(&GetTokenResponse{Token: "test_token"}, nil)
	repo := Repository{Repo: "git@github.com:CompuZest/test_repo.git", Name: "test_repo"}
	list := RepositoryList{Items: []Repository{repo}}
	r1    := ioutil.NopCloser(bytes.NewReader([]byte{}))
	mockArgocdAPI.EXPECT().ListRepositories(gomock.Any(), gomock.Any()).Return(&list, &http.Response{Body: r1}, nil)

	log := ctrl.Log.WithName("TestRegisterRepoExistingRepo")
	registered, err := RegisterRepo(log, mockArgocdAPI, RepoOpts{RepoUrl: repo.Repo})
	assert.False(t, registered)
	assert.NoError(t, err)
}

