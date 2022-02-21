package watcher_test

import (
	"context"
	"fmt"
	"testing"

	github2 "github.com/compuzest/zlifecycle-il-operator/controllers/external/github"

	"github.com/compuzest/zlifecycle-il-operator/controllers/external/argocd"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util"
	"github.com/compuzest/zlifecycle-il-operator/controllers/watcher"

	"github.com/golang/mock/gomock"
	"github.com/google/go-github/v42/github"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
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
	mockArgocdResponse := util.CreateMockResponse(200)
	mockGitHubResponse := util.CreateMockGithubResponse(200)

	mockGitClient := github2.NewMockAPI(mockCtrl)
	mockGitClient.EXPECT().FindRepositoryInstallation(owner, repository).Return(&mockInstallation, mockGitHubResponse, nil)
	mockArgocdClient := argocd.NewMockAPI(mockCtrl)
	mockToken := &argocd.GetTokenResponse{Token: "test_token"}
	mockArgocdClient.EXPECT().GetAuthToken().Return(mockToken, nil)
	list := argocd.RepositoryList{Items: []argocd.Repository{}}
	token := fmt.Sprintf("Bearer %s", mockToken.Token)
	mockArgocdClient.EXPECT().ListRepositories(token).Return(&list, util.CreateMockResponse(200), nil)
	mockArgocdClient.EXPECT().CreateRepository(gomock.Any(), token).Return(mockArgocdResponse, nil)
	logger := logrus.New().WithField("name", "TestLogger")
	r, err := watcher.NewGitHubAppWatcher(ctx, mockGitClient, mockArgocdClient, []byte("test"), logger)
	assert.NoError(t, err)
	assert.IsType(t, r, &watcher.GitHubAppWatcher{})

	err = r.Watch(repoURL)
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

	mockArgocdResponse := util.CreateMockResponse(200)

	mockArgocdClient := argocd.NewMockAPI(mockCtrl)
	mockToken := &argocd.GetTokenResponse{Token: "test_token"}
	mockArgocdClient.EXPECT().GetAuthToken().Return(mockToken, nil)
	list := argocd.RepositoryList{Items: []argocd.Repository{}}
	token := fmt.Sprintf("Bearer %s", mockToken.Token)
	mockArgocdClient.EXPECT().ListRepositories(token).Return(&list, util.CreateMockResponse(200), nil)
	mockArgocdClient.EXPECT().CreateRepository(gomock.Any(), token).Return(mockArgocdResponse, nil)
	logger := logrus.New().WithField("name", "TestLogger")
	testPrivateKey := []byte("test")
	r := watcher.NewSSHWatcher(ctx, mockArgocdClient, testPrivateKey, logger)
	assert.IsType(t, r, &watcher.SSHWatcher{})

	err := r.Watch(repoURL)
	assert.NoError(t, err)
}
