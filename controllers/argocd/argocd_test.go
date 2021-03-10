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

	repoOpts := RepoOpts{RepoUrl: "git@github.com:CompuZest/test_repo.git", SshPrivateKey: "test_key"}

	mockArgocdAPI.EXPECT().GetAuthToken().Return(&GetTokenResponse{Token: "test_token"}, nil)
	repo := Repository{Repo: "git@github.com:CompuZest/test_repo2.git", Name: "test_repo2"}
	list := RepositoryList{Items: []Repository{repo}}
	r1    := ioutil.NopCloser(bytes.NewReader([]byte{}))
	mockArgocdAPI.EXPECT().ListRepositories(gomock.Any()).Return(&list, &http.Response{Body: r1}, nil)
	r2    := ioutil.NopCloser(bytes.NewReader([]byte{}))
	createRepoBody := CreateRepoBody{Name: "test_repo", Repo: "git@github.com:CompuZest/test_repo.git", SshPrivateKey: "test_key"}
	mockArgocdAPI.EXPECT().CreateRepository(createRepoBody, gomock.Any()).Return(&http.Response{Body: r2}, nil)

	log := ctrl.Log.WithName("TestRegisterRepoNewRepo")
	registered, err := RegisterRepo(log, mockArgocdAPI, repoOpts)
	assert.True(t, registered)
	assert.NoError(t, err)
}

func TestRegisterRepoExistingRepo(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockArgocdAPI := NewMockArgocdAPI(mockCtrl)

	mockArgocdAPI.EXPECT().GetAuthToken().Return(&GetTokenResponse{Token: "test_token"}, nil)
	repo := Repository{Repo: "git@github.com:CompuZest/test_repo.git", Name: "test_repo"}
	list := RepositoryList{Items: []Repository{repo}}
	r1    := ioutil.NopCloser(bytes.NewReader([]byte{}))
	mockArgocdAPI.EXPECT().ListRepositories(gomock.Any()).Return(&list, &http.Response{Body: r1}, nil)

	log := ctrl.Log.WithName("TestRegisterRepoExistingRepo")

	repoOpts := RepoOpts{RepoUrl: "git@github.com:CompuZest/test_repo.git", SshPrivateKey: "test_key"}
	registered, err := RegisterRepo(log, mockArgocdAPI, repoOpts)
	assert.False(t, registered)
	assert.NoError(t, err)
}

