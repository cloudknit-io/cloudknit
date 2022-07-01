package watcher_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/compuzest/zlifecycle-il-operator/controller/services/watcher"

	argocdapi "github.com/compuzest/zlifecycle-il-operator/controller/common/argocd"
	"github.com/compuzest/zlifecycle-il-operator/controller/util"
	"github.com/golang/mock/gomock"
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

	mockArgocdResponse := util.CreateMockResponse(200)

	mockArgocdClient := argocdapi.NewMockAPI(mockCtrl)
	mockToken := &argocdapi.GetTokenResponse{Token: "test_token"}
	mockArgocdClient.EXPECT().GetAuthToken().Return(mockToken, nil)
	list := argocdapi.RepositoryList{Items: []argocdapi.Repository{}}
	token := fmt.Sprintf("Bearer %s", mockToken.Token)
	mockArgocdClient.EXPECT().ListRepositories(token).Return(&list, util.CreateMockResponse(200), nil)
	mockArgocdClient.EXPECT().CreateRepository(gomock.Any(), token).Return(mockArgocdResponse, nil)
	logger := logrus.New().WithField("name", "TestLogger")
	r, err := watcher.NewGitHubAppWatcher(ctx, 1, 2, mockArgocdClient, []byte("test"), logger)
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

	mockArgocdClient := argocdapi.NewMockAPI(mockCtrl)
	mockToken := &argocdapi.GetTokenResponse{Token: "test_token"}
	mockArgocdClient.EXPECT().GetAuthToken().Return(mockToken, nil)
	list := argocdapi.RepositoryList{Items: []argocdapi.Repository{}}
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
