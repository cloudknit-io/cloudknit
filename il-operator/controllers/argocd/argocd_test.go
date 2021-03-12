package argocd

import (
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/common"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	ctrl "sigs.k8s.io/controller-runtime"
	"testing"
)

func TestRegisterRepoNewRepo(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockApi := NewMockApi(mockCtrl)

	repoOpts := RepoOpts{RepoUrl: "git@github.com:CompuZest/test_repo.git", SshPrivateKey: "test_key"}

	mockApi.EXPECT().GetAuthToken().Return(&GetTokenResponse{Token: "test_token"}, nil)
	repo := Repository{Repo: "git@github.com:CompuZest/test_repo2.git", Name: "test_repo2"}
	list := RepositoryList{Items: []Repository{repo}}
	mockApi.EXPECT().ListRepositories(gomock.Any()).Return(&list, common.CreateMockResponse(200), nil)
	createRepoBody := CreateRepoBody{Name: "test_repo", Repo: "git@github.com:CompuZest/test_repo.git", SshPrivateKey: "test_key"}
	mockApi.EXPECT().CreateRepository(createRepoBody, gomock.Any()).Return(common.CreateMockResponse(200), nil)

	log := ctrl.Log.WithName("TestRegisterRepoNewRepo")
	registered, err := RegisterRepo(log, mockApi, repoOpts)
	assert.True(t, registered)
	assert.NoError(t, err)
}

func TestRegisterRepoExistingRepo(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockArgocdAPI := NewMockApi(mockCtrl)

	mockArgocdAPI.EXPECT().GetAuthToken().Return(&GetTokenResponse{Token: "test_token"}, nil)
	repo := Repository{Repo: "git@github.com:CompuZest/test_repo.git", Name: "test_repo"}
	list := RepositoryList{Items: []Repository{repo}}
	mockArgocdAPI.EXPECT().ListRepositories(gomock.Any()).Return(&list, common.CreateMockResponse(200), nil)

	log := ctrl.Log.WithName("TestRegisterRepoExistingRepo")

	repoOpts := RepoOpts{RepoUrl: "git@github.com:CompuZest/test_repo.git", SshPrivateKey: "test_key"}
	registered, err := RegisterRepo(log, mockArgocdAPI, repoOpts)
	assert.False(t, registered)
	assert.NoError(t, err)
}

