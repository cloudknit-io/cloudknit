package github_test

import (
	"fmt"
	"testing"

	github2 "github.com/compuzest/zlifecycle-il-operator/controller/external/github"
	"github.com/compuzest/zlifecycle-il-operator/controller/util"

	"github.com/sirupsen/logrus"

	"github.com/golang/mock/gomock"
	"github.com/google/go-github/v42/github"
	"github.com/stretchr/testify/assert"
)

func TestTryCreateRepositoryExisting(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockClient := github2.NewMockAPI(mockCtrl)

	testOwner, testRepo := "compuzest1", "test_repo1"
	testURL := fmt.Sprintf("git@github.com:%s/%s", testOwner, testRepo)
	testGitHubRepo := github.Repository{Name: &testRepo, URL: &testURL}
	mockResponse := github.Response{Response: util.CreateMockResponse(200)}
	mockClient.EXPECT().GetRepository(testOwner, testRepo).Return(&testGitHubRepo, &mockResponse, nil)

	log := logrus.NewEntry(logrus.New())

	created, err := github2.CreateRepository(log, mockClient, testOwner, testRepo)
	assert.False(t, created)
	assert.NoError(t, err)
}

func TestTryCreateRepositoryNew(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockClient := github2.NewMockAPI(mockCtrl)

	testOwner, testRepo := "compuzest3", "test_repo3"
	testURL := fmt.Sprintf("git@github.com:%s/%s", testOwner, testRepo)
	mockResponse1 := github.Response{Response: util.CreateMockResponse(200)}
	mockClient.EXPECT().GetRepository(testOwner, testRepo).Return(nil, &mockResponse1, nil)
	testGitHubRepo := github.Repository{Name: &testRepo, URL: &testURL}
	mockResponse2 := github.Response{Response: util.CreateMockResponse(200)}
	mockClient.EXPECT().CreateRepository(testOwner, testRepo).Return(&testGitHubRepo, &mockResponse2, nil)

	log := logrus.NewEntry(logrus.New())

	created, err := github2.CreateRepository(log, mockClient, testOwner, testRepo)
	assert.True(t, created)
	assert.NoError(t, err)
}

func TestCreateRepoWebhookNew(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockClient := github2.NewMockAPI(mockCtrl)

	active := true
	events := []string{"push"}

	testOwner, testRepo := "CompuZest2", "test_repo2"
	testPayloadURL1 := "https://test1.webhook.com"
	testCfg1 := map[string]interface{}{"content_type": "json", "url": testPayloadURL1}
	testHook1 := github.Hook{Active: &active, Config: testCfg1, Events: events}
	mockResponse1 := github.Response{Response: util.CreateMockResponse(200)}
	mockClient.EXPECT().ListHooks(
		gomock.Eq(testOwner),
		gomock.Eq(testRepo),
		gomock.Nil(),
	).Return([]*github.Hook{&testHook1}, &mockResponse1, nil)
	testPayloadURL2 := "https://test2.webhook.com"
	webHookSecret2 := "secret"
	testCfg2 := map[string]interface{}{"content_type": "json", "secret": "secret", "url": testPayloadURL2}
	testHook2 := github.Hook{Active: &active, Config: testCfg2, Events: events}
	expectedID := int64(1)
	expectedHook := github.Hook{Active: &active, Config: testCfg2, Events: events, URL: &testPayloadURL2, ID: &expectedID}
	mockResponse2 := github.Response{Response: util.CreateMockResponse(200)}
	mockClient.EXPECT().CreateHook(
		gomock.Eq(testOwner),
		gomock.Eq(testRepo),
		gomock.Eq(&testHook2),
	).Return(&expectedHook, &mockResponse2, nil)

	log := logrus.NewEntry(logrus.New())

	testRepoURL := fmt.Sprintf("git@github.com:%s/%s", testOwner, testRepo)
	created, err := github2.CreateRepoWebhook(log, mockClient, testRepoURL, testPayloadURL2, webHookSecret2)
	assert.False(t, created)
	assert.NoError(t, err)
}

func TestCreateRepoWebhookExisting(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockClient := github2.NewMockAPI(mockCtrl)

	active := true
	events := []string{"push"}

	testOwner, testRepo := "CompuZest", "test_repo"
	testPayloadURL1 := "https://test11.webhook.com"
	webHookSecret1 := "secret"
	testCfg1 := map[string]interface{}{"content_type": "json", "url": testPayloadURL1}
	testHook1 := github.Hook{Active: &active, Config: testCfg1, Events: events}
	mockResponse1 := github.Response{Response: util.CreateMockResponse(200)}
	mockClient.EXPECT().ListHooks(
		gomock.Eq(testOwner),
		gomock.Eq(testRepo),
		gomock.Nil(),
	).Return([]*github.Hook{&testHook1}, &mockResponse1, nil)

	log := logrus.NewEntry(logrus.New())

	testRepoURL := fmt.Sprintf("git@github.com:%s/%s", testOwner, testRepo)
	created, err := github2.CreateRepoWebhook(log, mockClient, testRepoURL, testPayloadURL1, webHookSecret1)
	assert.True(t, created)
	assert.NoError(t, err)
}
