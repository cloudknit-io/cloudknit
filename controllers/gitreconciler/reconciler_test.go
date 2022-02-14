package gitreconciler_test

import (
	"context"
	"fmt"
	"github.com/compuzest/zlifecycle-il-operator/controllers/gitreconciler"
	"github.com/compuzest/zlifecycle-il-operator/mocks"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"testing"
)

var ctx = context.Background()

func TestReconciler_Subscribe(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockClient := mocks.NewMockClient(mockCtrl)

	owner := "CompuZest"
	repository1 := "test1"
	repoURL1 := fmt.Sprintf("git@github.com:%s/%s", owner, repository1)
	repository2 := "test2"
	repoURL2 := fmt.Sprintf("git@github.com:%s/%s", owner, repository2)

	r := gitreconciler.NewReconciler(ctx, logrus.New().WithField("name", "TestLogger"), mockClient)

	testSubscriber1 := client.ObjectKey{Name: "test", Namespace: "test"}
	result := r.Subscribe(repoURL1, testSubscriber1)
	assert.False(t, result)

	result = r.Subscribe(repoURL1, testSubscriber1)
	assert.True(t, result)

	testSubscriber2 := client.ObjectKey{Name: "test2", Namespace: "test2"}
	result = r.Subscribe(repoURL1, testSubscriber2)
	assert.False(t, result)

	result = r.Subscribe(repoURL2, testSubscriber2)
	assert.False(t, result)

	repos := r.Repositories()
	assert.Len(t, repos, 2)
}

func TestReconciler_UnsubscribeAll(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockClient := mocks.NewMockClient(mockCtrl)

	owner := "CompuZest"
	repository1 := "test1"
	repoURL1 := fmt.Sprintf("git@github.com:%s/%s", owner, repository1)

	r := gitreconciler.NewReconciler(ctx, logrus.New().WithField("name", "TestLogger"), mockClient)

	testSubscriber1 := client.ObjectKey{Name: "test", Namespace: "test"}
	result := r.Subscribe(repoURL1, testSubscriber1)
	assert.False(t, result)

	testSubscriber2 := client.ObjectKey{Name: "test2", Namespace: "test2"}
	result = r.Subscribe(repoURL1, testSubscriber2)
	assert.False(t, result)

	err := r.UnsubscribeAll(testSubscriber1)
	assert.NoError(t, err)

	repos := r.Repositories()
	assert.Len(t, repos, 1)
	assert.NotContains(t, repos[0].Subscribers, testSubscriber1)
}
