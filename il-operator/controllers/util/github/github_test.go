package github_test

import (
	"fmt"
	"github.com/compuzest/zlifecycle-il-operator/mocks"
	"testing"

	"github.com/compuzest/zlifecycle-il-operator/controllers/util/common"
	github2 "github.com/compuzest/zlifecycle-il-operator/controllers/util/github"
	"github.com/golang/mock/gomock"
	"github.com/google/go-github/v32/github"
	"github.com/stretchr/testify/assert"
	ctrl "sigs.k8s.io/controller-runtime"
)

func TestTryCreateRepositoryExisting(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRepositoryAPI := mocks.NewMockRepositoryApi(mockCtrl)

	testOwner := "compuzest"
	testRepo := "test_repo"
	testURL := fmt.Sprintf("git@github.com:%s/%s", testOwner, testRepo)
	testGitHubRepo := github.Repository{Name: &testRepo, URL: &testURL}
	mockResponse := github.Response{Response: common.CreateMockResponse(200)}
	mockRepositoryAPI.EXPECT().GetRepository(testOwner, testRepo).Return(&testGitHubRepo, &mockResponse, nil)

	log := ctrl.Log.WithName("TestTryCreateRepositoryExisting")

	created, err := github2.TryCreateRepository(log, mockRepositoryAPI, testOwner, testRepo)
	assert.False(t, created)
	assert.NoError(t, err)
}

func TestTryCreateRepositoryNew(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRepositoryAPI := mocks.NewMockRepositoryApi(mockCtrl)

	testOwner := "compuzest"
	testRepo := "test_repo"
	testURL := fmt.Sprintf("git@github.com:%s/%s", testOwner, testRepo)
	mockResponse1 := github.Response{Response: common.CreateMockResponse(200)}
	mockRepositoryAPI.EXPECT().GetRepository(testOwner, testRepo).Return(nil, &mockResponse1, nil)
	testGitHubRepo := github.Repository{Name: &testRepo, URL: &testURL}
	mockResponse2 := github.Response{Response: common.CreateMockResponse(200)}
	mockRepositoryAPI.EXPECT().CreateRepository(testOwner, testRepo).Return(&testGitHubRepo, &mockResponse2, nil)

	log := ctrl.Log.WithName("TestTryCreateRepositoryNew")

	created, err := github2.TryCreateRepository(log, mockRepositoryAPI, testOwner, testRepo)
	assert.True(t, created)
	assert.NoError(t, err)
}

func TestCreateRepoWebhookNew(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRepositoryAPI := mocks.NewMockRepositoryApi(mockCtrl)

	active := true
	events := []string{"push"}

	testOwner, testRepo := "CompuZest", "test_repo"
	testPayloadURL1 := "https://test1.webhook.com"
	testCfg1 := map[string]interface{}{"content_type": "json", "url": testPayloadURL1}
	testHook1 := github.Hook{Active: &active, Config: testCfg1, Events: events}
	mockResponse1 := github.Response{Response: common.CreateMockResponse(200)}
	mockRepositoryAPI.EXPECT().ListHooks(
		gomock.Eq(testOwner),
		gomock.Eq(testRepo),
		gomock.Nil(),
	).Return([]*github.Hook{&testHook1}, &mockResponse1, nil)
	testPayloadURL2 := "https://test2.webhook.com"
	webHookSecret2 := "secret"
	testCfg2 := map[string]interface{}{"content_type": "json", "secret": "secret", "url": testPayloadURL2}
	testHook2 := github.Hook{Active: &active, Config: testCfg2, Events: events}
	expectedURL := "https://test2.webhook.com"
	expectedID := int64(1)
	expectedHook := github.Hook{Active: &active, Config: testCfg2, Events: events, URL: &expectedURL, ID: &expectedID}
	mockResponse2 := github.Response{Response: common.CreateMockResponse(200)}
	mockRepositoryAPI.EXPECT().CreateHook(
		gomock.Eq(testOwner),
		gomock.Eq(testRepo),
		gomock.Eq(&testHook2),
	).Return(&expectedHook, &mockResponse2, nil)

	log := ctrl.Log.WithName("TestCreateRepoWebhookNew")

	testRepoURL := fmt.Sprintf("git@github.com:%s/%s", testOwner, testRepo)
	created, err := github2.CreateRepoWebhook(log, mockRepositoryAPI, testRepoURL, testPayloadURL2, webHookSecret2)
	assert.False(t, created)
	assert.NoError(t, err)
}

func TestCreateRepoWebhookExisting(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRepositoryAPI := mocks.NewMockRepositoryApi(mockCtrl)

	active := true
	events := []string{"push"}

	testOwner, testRepo := "CompuZest", "test_repo"
	testPayloadURL1 := "https://test1.webhook.com"
	webHookSecret1 := "secret"
	testCfg1 := map[string]interface{}{"content_type": "json", "url": testPayloadURL1}
	testHook1 := github.Hook{Active: &active, Config: testCfg1, Events: events}
	mockResponse1 := github.Response{Response: common.CreateMockResponse(200)}
	mockRepositoryAPI.EXPECT().ListHooks(
		gomock.Eq(testOwner),
		gomock.Eq(testRepo),
		gomock.Nil(),
	).Return([]*github.Hook{&testHook1}, &mockResponse1, nil)

	log := ctrl.Log.WithName("TestCreateRepoWebhookExisting")

	testRepoURL := fmt.Sprintf("git@github.com:%s/%s", testOwner, testRepo)
	created, err := github2.CreateRepoWebhook(log, mockRepositoryAPI, testRepoURL, testPayloadURL1, webHookSecret1)
	assert.True(t, created)
	assert.NoError(t, err)
}
