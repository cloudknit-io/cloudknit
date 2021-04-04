package argocd_test

import (
	"testing"

	"github.com/compuzest/zlifecycle-il-operator/controllers/argocd"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/common"
	"github.com/compuzest/zlifecycle-il-operator/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	ctrl "sigs.k8s.io/controller-runtime"
)

func TestRegisterRepoNewRepo(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockApi := mocks.NewMockApi(mockCtrl)

	repoOpts := argocd.RepoOpts{RepoUrl: "git@github.com:CompuZest/test_repo.git", SshPrivateKey: "test_key"}

	mockApi.EXPECT().GetAuthToken().Return(&argocd.GetTokenResponse{Token: "test_token"}, nil)
	repo := argocd.Repository{Repo: "git@github.com:CompuZest/test_repo2.git", Name: "test_repo2"}
	list := argocd.RepositoryList{Items: []argocd.Repository{repo}}
	mockApi.EXPECT().ListRepositories(gomock.Any()).Return(&list, common.CreateMockResponse(200), nil)
	createRepoBody := argocd.CreateRepoBody{Name: "test_repo", Repo: "git@github.com:CompuZest/test_repo.git", SshPrivateKey: "test_key"}
	mockApi.EXPECT().CreateRepository(createRepoBody, gomock.Any()).Return(common.CreateMockResponse(200), nil)

	log := ctrl.Log.WithName("TestRegisterRepoNewRepo")
	registered, err := argocd.RegisterRepo(log, mockApi, repoOpts)
	assert.True(t, registered)
	assert.NoError(t, err)
}

func TestRegisterRepoExistingRepo(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockArgocdAPI := mocks.NewMockApi(mockCtrl)

	mockArgocdAPI.EXPECT().GetAuthToken().Return(&argocd.GetTokenResponse{Token: "test_token"}, nil)
	repo := argocd.Repository{Repo: "git@github.com:CompuZest/test_repo.git", Name: "test_repo"}
	list := argocd.RepositoryList{Items: []argocd.Repository{repo}}
	mockArgocdAPI.EXPECT().ListRepositories(gomock.Any()).Return(&list, common.CreateMockResponse(200), nil)

	log := ctrl.Log.WithName("TestRegisterRepoExistingRepo")

	repoOpts := argocd.RepoOpts{RepoUrl: "git@github.com:CompuZest/test_repo.git", SshPrivateKey: "test_key"}
	registered, err := argocd.RegisterRepo(log, mockArgocdAPI, repoOpts)
	assert.False(t, registered)
	assert.NoError(t, err)
}
