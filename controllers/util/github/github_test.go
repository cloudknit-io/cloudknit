package github

import (
	"fmt"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/common"
	"github.com/golang/mock/gomock"
	"github.com/google/go-github/v32/github"
	"github.com/stretchr/testify/assert"
	ctrl "sigs.k8s.io/controller-runtime"
	"testing"
)

func TestTryCreateRepositoryExisting(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRepositoryApi := NewMockRepositoryApi(mockCtrl)

	testOwner := "compuzest"
	testRepo  := "test_repo"
	testURL   := fmt.Sprintf("git@github.com:%s/%s", testOwner, testRepo)
	testGitHubRepo := github.Repository{Name: &testRepo, URL: &testURL}
	mockResponse := github.Response{Response: common.CreateMockResponse(200)}
	mockRepositoryApi.EXPECT().GetRepository(testOwner, testRepo).Return(&testGitHubRepo, &mockResponse, nil)

	log := ctrl.Log.WithName("TestTryCreateRepositoryExisting")

	created, err := TryCreateRepository(log, mockRepositoryApi, testOwner, testRepo)
	assert.False(t, created)
	assert.NoError(t, err)
}

func TestTryCreateRepositoryNew(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRepositoryApi := NewMockRepositoryApi(mockCtrl)

	testOwner := "compuzest"
	testRepo  := "test_repo"
	testURL   := fmt.Sprintf("git@github.com:%s/%s", testOwner, testRepo)
	mockResponse1 := github.Response{Response: common.CreateMockResponse(200)}
	mockRepositoryApi.EXPECT().GetRepository(testOwner, testRepo).Return(nil, &mockResponse1, nil)
	testGitHubRepo := github.Repository{Name: &testRepo, URL: &testURL}
	mockResponse2 := github.Response{Response: common.CreateMockResponse(200)}
	mockRepositoryApi.EXPECT().CreateRepository(testOwner, testRepo).Return(&testGitHubRepo, &mockResponse2, nil)

	log := ctrl.Log.WithName("TestTryCreateRepositoryNew")

	created, err := TryCreateRepository(log, mockRepositoryApi, testOwner, testRepo)
	assert.True(t, created)
	assert.NoError(t, err)
}

func TestCreateRepoWebhookNew(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRepositoryApi := NewMockRepositoryApi(mockCtrl)

	active := true
	events := []string{"push"}

	testOwner, testRepo := "CompuZest", "test_repo"
	testPayloadUrl1 := "https://test1.webhook.com"
	testCfg1 := map[string]interface{}{"content_type": "json", "url": testPayloadUrl1}
	testHook1 := github.Hook{Active: &active, Config: testCfg1, Events: events}
	mockResponse1 := github.Response{Response: common.CreateMockResponse(200)}
	mockRepositoryApi.EXPECT().ListHooks(
		gomock.Eq(testOwner),
		gomock.Eq(testRepo),
		gomock.Nil(),
	).Return([]*github.Hook{&testHook1}, &mockResponse1, nil)
	testPayloadUrl2 := "https://test2.webhook.com"
	testCfg2 := map[string]interface{}{"content_type": "json", "url": testPayloadUrl2}
	testHook2 := github.Hook{Active: &active, Config: testCfg2, Events: events}
	expectedUrl := "https://test2.webhook.com"
	expectedId := int64(1)
	expectedHook := github.Hook{Active: &active, Config: testCfg2, Events: events, URL: &expectedUrl, ID: &expectedId}
	mockResponse2 := github.Response{Response: common.CreateMockResponse(200)}
	mockRepositoryApi.EXPECT().CreateHook(
		gomock.Eq(testOwner),
		gomock.Eq(testRepo),
		gomock.Eq(&testHook2),
	).Return(&expectedHook, &mockResponse2, nil)

	log := ctrl.Log.WithName("TestCreateRepoWebhookNew")

	testRepoUrl := fmt.Sprintf("git@github.com:%s/%s", testOwner, testRepo)
	created, err := CreateRepoWebhook(log, mockRepositoryApi, testRepoUrl, testPayloadUrl2)
	assert.False(t, created)
	assert.NoError(t, err)
}

func TestCreateRepoWebhookExisting(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRepositoryApi := NewMockRepositoryApi(mockCtrl)

	active := true
	events := []string{"push"}

	testOwner, testRepo := "CompuZest", "test_repo"
	testPayloadUrl1 := "https://test1.webhook.com"
	testCfg1 := map[string]interface{}{"content_type": "json", "url": testPayloadUrl1}
	testHook1 := github.Hook{Active: &active, Config: testCfg1, Events: events}
	mockResponse1 := github.Response{Response: common.CreateMockResponse(200)}
	mockRepositoryApi.EXPECT().ListHooks(
		gomock.Eq(testOwner),
		gomock.Eq(testRepo),
		gomock.Nil(),
	).Return([]*github.Hook{&testHook1}, &mockResponse1, nil)

	log := ctrl.Log.WithName("TestCreateRepoWebhookExisting")

	testRepoUrl := fmt.Sprintf("git@github.com:%s/%s", testOwner, testRepo)
	created, err := CreateRepoWebhook(log, mockRepositoryApi, testRepoUrl, testPayloadUrl1)
	assert.True(t, created)
	assert.NoError(t, err)
}
