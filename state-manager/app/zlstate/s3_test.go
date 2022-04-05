package zlstate_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	smithyhttp "github.com/aws/smithy-go/transport/http"
	"github.com/compuzest/zlifecycle-state-manager/app/util"
	"github.com/compuzest/zlifecycle-state-manager/app/zlog"
	"github.com/compuzest/zlifecycle-state-manager/app/zlstate"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var (
	ctx = context.Background()
	log = zlog.PlainLogger().WithField("name", "TestLogger")
)

func TestS3Backend_GetExisting(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockS3API := zlstate.NewMockS3API(mockCtrl)

	bucket := "testBucket"
	key := "testKey"

	hoi := &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}
	mockS3API.EXPECT().HeadObject(gomock.Any(), hoi).Return(nil, nil)
	goi := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}
	testState := zlstate.ZLState{
		Company:     "compuzest",
		Team:        "test",
		Environment: "testEnv",
	}
	response := s3.GetObjectOutput{Body: util.CreateMockBody(&testState)}
	mockS3API.EXPECT().GetObject(gomock.Any(), goi).Return(&response, nil)

	s3Backend := zlstate.NewS3Backend(ctx, log, bucket, mockS3API)
	zlState, err := s3Backend.Get(key)
	assert.NoError(t, err)
	assert.Equal(t, testState.Company, zlState.Company)
	assert.Equal(t, testState.Team, zlState.Team)
	assert.Equal(t, testState.Environment, zlState.Environment)
}

func TestS3Backend_GetNonExisting(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockS3API := zlstate.NewMockS3API(mockCtrl)

	bucket := "testBucket"
	key := "testKey"

	hoi := &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}
	awsErr := &awshttp.ResponseError{
		ResponseError: &smithyhttp.ResponseError{
			Response: &smithyhttp.Response{Response: util.CreateMockResponse(http.StatusNotFound)},
			Err:      nil,
		},
		RequestID: "test",
	}
	mockS3API.EXPECT().HeadObject(gomock.Any(), hoi).Return(nil, awsErr)

	s3Backend := zlstate.NewS3Backend(ctx, log, bucket, mockS3API)
	zlState, err := s3Backend.Get(key)
	assert.Nil(t, zlState)
	assert.ErrorIs(t, err, zlstate.ErrKeyNotExists)
}

func TestS3Backend_PutNonExisting(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockS3API := zlstate.NewMockS3API(mockCtrl)

	bucket := "testBucket"
	key := "testKey"

	hoi := &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}
	awsErr := &awshttp.ResponseError{
		ResponseError: &smithyhttp.ResponseError{
			Response: &smithyhttp.Response{Response: util.CreateMockResponse(http.StatusNotFound)},
			Err:      nil,
		},
		RequestID: "test",
	}
	mockS3API.EXPECT().HeadObject(gomock.Any(), hoi).Return(nil, awsErr)

	testState := zlstate.ZLState{
		Company:     "compuzest",
		Team:        "test",
		Environment: "testEnv",
	}
	mockS3API.EXPECT().PutObject(gomock.Any(), gomock.Any()).Return(nil, nil)

	s3Backend := zlstate.NewS3Backend(ctx, log, bucket, mockS3API)

	err := s3Backend.Put(key, &testState, false)
	assert.NoError(t, err)
}

func TestS3Backend_PutExisting(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockS3API := zlstate.NewMockS3API(mockCtrl)

	bucket := "testBucket"
	key := "testKey"

	hoi := &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}
	mockS3API.EXPECT().HeadObject(gomock.Any(), hoi).Return(nil, nil)

	testState := zlstate.ZLState{
		Company:     "compuzest",
		Team:        "test",
		Environment: "testEnv",
	}

	s3Backend := zlstate.NewS3Backend(ctx, log, bucket, mockS3API)

	err := s3Backend.Put(key, &testState, false)
	assert.ErrorIs(t, err, zlstate.ErrKeyAlreadyExists)
}

func TestS3Backend_PutExistingForce(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockS3API := zlstate.NewMockS3API(mockCtrl)

	bucket := "testBucket"
	key := "testKey"

	testState := zlstate.ZLState{
		Company:     "compuzest",
		Team:        "test",
		Environment: "testEnv",
	}
	mockS3API.EXPECT().PutObject(gomock.Any(), gomock.Any()).Return(nil, nil)

	s3Backend := zlstate.NewS3Backend(ctx, log, bucket, mockS3API)

	err := s3Backend.Put(key, &testState, true)
	assert.NoError(t, err)
}

func TestS3Backend_UpsertComponentNew(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockS3API := zlstate.NewMockS3API(mockCtrl)

	bucket := "testBucket"
	key := "testKey"

	hoi := &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}
	mockS3API.EXPECT().HeadObject(gomock.Any(), hoi).Return(nil, nil)
	goi := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	comp1 := zlstate.Component{Name: "comp1", Type: "terraform"}
	testState := zlstate.ZLState{
		Company:     "compuzest",
		Team:        "test",
		Environment: "testEnv",
		Components:  []*zlstate.Component{&comp1},
	}
	response := s3.GetObjectOutput{Body: util.CreateMockBody(&testState)}
	mockS3API.EXPECT().GetObject(gomock.Any(), goi).Return(&response, nil)
	mockS3API.EXPECT().PutObject(gomock.Any(), gomock.Any()).Return(nil, nil)

	comp2 := zlstate.Component{Name: "comp2", Type: "argocd"}

	s3Backend := zlstate.NewS3Backend(ctx, log, bucket, mockS3API)
	zlst, err := s3Backend.UpsertComponent(key, &comp2)
	assert.NoError(t, err)
	assert.Len(t, zlst.Components, 2)
	assert.Equal(t, zlst.Components[0].Name, "comp1")
	assert.Equal(t, zlst.Components[0].Type, "terraform")
	assert.Equal(t, zlst.Components[1].Name, "comp2")
	assert.Equal(t, zlst.Components[1].Type, "argocd")
}

func TestS3Backend_UpsertComponentExisting(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockS3API := zlstate.NewMockS3API(mockCtrl)

	bucket := "testBucket"
	key := "testKey"

	hoi := &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}
	mockS3API.EXPECT().HeadObject(gomock.Any(), hoi).Return(nil, nil)
	goi := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	comp := zlstate.Component{Name: "comp1", Type: "terraform"}
	testState := zlstate.ZLState{
		Company:     "compuzest",
		Team:        "test",
		Environment: "testEnv",
		Components:  []*zlstate.Component{&comp},
	}
	response := s3.GetObjectOutput{Body: util.CreateMockBody(&testState)}
	mockS3API.EXPECT().GetObject(gomock.Any(), goi).Return(&response, nil)
	mockS3API.EXPECT().PutObject(gomock.Any(), gomock.Any()).Return(nil, nil)

	newComp := zlstate.Component{Name: "comp1", Type: "argocd"}

	s3Backend := zlstate.NewS3Backend(ctx, log, bucket, mockS3API)
	zlst, err := s3Backend.UpsertComponent(key, &newComp)
	assert.NoError(t, err)
	assert.Len(t, zlst.Components, 1)
	assert.Equal(t, zlst.Components[0].Name, "comp1")
	assert.Equal(t, zlst.Components[0].Type, "argocd")
}

func TestS3Backend_PatchComponent(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockS3API := zlstate.NewMockS3API(mockCtrl)

	bucket := "testBucket"
	key := "testKey"

	hoi := &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}
	mockS3API.EXPECT().HeadObject(gomock.Any(), hoi).Return(nil, nil)
	goi := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	comp1 := zlstate.Component{Name: "comp1", Type: "terraform", Status: "not_provisioned"}
	comp2 := zlstate.Component{Name: "comp2", Type: "argocd", Status: "not_provisioned"}

	testState := zlstate.ZLState{
		Company:     "compuzest",
		Team:        "test",
		Environment: "testEnv",
		Components:  []*zlstate.Component{&comp1, &comp2},
	}
	response := s3.GetObjectOutput{Body: util.CreateMockBody(&testState)}
	mockS3API.EXPECT().GetObject(gomock.Any(), goi).Return(&response, nil)
	mockS3API.EXPECT().PutObject(gomock.Any(), gomock.Any()).Return(nil, nil)

	s3Backend := zlstate.NewS3Backend(ctx, log, bucket, mockS3API)
	zlst, err := s3Backend.PatchComponent(key, "comp1", "provisioned")
	assert.NoError(t, err)
	assert.Equal(t, zlst.Components[0].Name, "comp1")
	assert.Equal(t, zlst.Components[0].Status, "provisioned")
	assert.Equal(t, zlst.Components[0].Type, "terraform")
	assert.Equal(t, zlst.Components[1].Name, "comp2")
	assert.Equal(t, zlst.Components[1].Status, "not_provisioned")
	assert.Equal(t, zlst.Components[1].Type, "argocd")

}
