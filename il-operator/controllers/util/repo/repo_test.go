package repo_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/compuzest/zlifecycle-il-operator/controllers/argocd"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/common"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/repo"
	"github.com/compuzest/zlifecycle-il-operator/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/go-github/v42/github"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
)

func TestGitHubAppTryRegisterRepo(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	ctx := context.Background()
	owner := "CompuZest"
	repository := "test1"
	repoURL := fmt.Sprintf("git@github.com:%s/%s", owner, repository)

	var mockAppID int64 = 1
	var mockInstallationID int64 = 2
	mockInstallation := github.Installation{AppID: &mockAppID, ID: &mockInstallationID}
	mockArgocdResponse := common.CreateMockResponse(200)
	mockGitHubResponse := common.CreateMockGithubResponse(200)
	mockClient := mocks.NewMockClient(mockCtrl)
	mockClient.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Do(
		func(ctx interface{}, key interface{}, s *v1.Secret) interface{} {
			data := map[string][]byte{}
			data["sshPrivateKey"] = []byte("testKey")
			s.Data = data
			return nil
		},
	)

	mockGitAppAPI := mocks.NewMockAppAPI(mockCtrl)
	mockGitAppAPI.EXPECT().FindRepositoryInstallation(owner, repository).Return(&mockInstallation, mockGitHubResponse, nil)
	mockArgocdAPI := mocks.NewMockAPI(mockCtrl)
	mockToken := &argocd.GetTokenResponse{Token: "test_token"}
	mockArgocdAPI.EXPECT().GetAuthToken().Return(mockToken, nil)
	list := argocd.RepositoryList{Items: []argocd.Repository{}}
	token := fmt.Sprintf("Bearer %s", mockToken.Token)
	mockArgocdAPI.EXPECT().ListRepositories(token).Return(&list, common.CreateMockResponse(200), nil)
	mockArgocdAPI.EXPECT().CreateRepository(gomock.Any(), token).Return(mockArgocdResponse, nil)
	logger := logrus.New().WithField("name", "TestLogger")
	r, err := repo.NewGitHubAppService(ctx, mockClient, mockGitAppAPI, mockArgocdAPI, logger)
	assert.NoError(t, err)
	assert.IsType(t, r, &repo.GitHubAppService{})

	err = r.TryRegisterRepo(repoURL)
	assert.NoError(t, err)
}

func TestSSHTryRegisterRepo(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	ctx := context.Background()
	owner := "CompuZest"
	repository := "test2"
	repoURL := fmt.Sprintf("git@github.com:%s/%s", owner, repository)

	mockArgocdResponse := common.CreateMockResponse(200)
	mockClient := mocks.NewMockClient(mockCtrl)
	mockClient.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Do(
		func(ctx interface{}, key interface{}, s *v1.Secret) interface{} {
			data := map[string][]byte{}
			data["sshPrivateKey"] = []byte("testKey")
			s.Data = data
			return nil
		},
	)

	mockArgocdAPI := mocks.NewMockAPI(mockCtrl)
	mockToken := &argocd.GetTokenResponse{Token: "test_token"}
	mockArgocdAPI.EXPECT().GetAuthToken().Return(mockToken, nil)
	list := argocd.RepositoryList{Items: []argocd.Repository{}}
	token := fmt.Sprintf("Bearer %s", mockToken.Token)
	mockArgocdAPI.EXPECT().ListRepositories(token).Return(&list, common.CreateMockResponse(200), nil)
	mockArgocdAPI.EXPECT().CreateRepository(gomock.Any(), token).Return(mockArgocdResponse, nil)
	logger := logrus.New().WithField("name", "TestLogger")
	r := repo.NewSSHService(ctx, mockClient, mockArgocdAPI, logger)
	assert.IsType(t, r, &repo.SSHService{})

	err := r.TryRegisterRepo(repoURL)
	assert.NoError(t, err)
}
