package zlstate

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/compuzest/zlifecycle-state-manager/app/util"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=./mock_s3.go -package=zlstate "github.com/compuzest/zlifecycle-state-manager/app/zlstate" S3API
type S3API interface {
	GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
	PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
	HeadObject(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error)
}

var (
	ErrKeyNotExists     = errors.New("key does not exist")
	ErrKeyAlreadyExists = errors.New("object already exists")
)

type S3Backend struct {
	ctx    context.Context
	log    *logrus.Entry
	bucket string
	s3     S3API
}

func NewS3Backend(ctx context.Context, log *logrus.Entry, bucket string, api S3API) *S3Backend {
	return &S3Backend{
		ctx:    ctx,
		log:    log,
		bucket: bucket,
		s3:     api,
	}
}

func NewS3Client(ctx context.Context) (*s3.Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "error loading default s3 config")
	}
	return s3.NewFromConfig(cfg), nil
}

// Get returns the state file whose key is the path in the bucket for which the backend was created
func (s *S3Backend) Get(key string) (*ZLState, error) {
	s.log.WithField("key", key).Info("Getting zLstate from remote backend [s3]")
	exists, err := s.exists(key)
	if err != nil {
		return nil, errors.Wrapf(err, "error checking does object already exist for key: [%s]", key)
	}
	if !exists {
		return nil, errors.WithStack(ErrKeyNotExists)
	}
	input := s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}
	output, err := s.s3.GetObject(s.ctx, &input)
	if err != nil {
		return nil, errors.Wrap(err, "error getting zLstate object from s3")
	}

	defer output.Body.Close()

	obj, err := util.ReadBody(output.Body)
	if err != nil {
		return nil, errors.Wrap(err, "error reading body of zLstate file")
	}

	var zlstate *ZLState
	if err := util.FromJSON(&zlstate, obj); err != nil {
		return nil, errors.Wrap(err, "error unmarshalling zLstate")
	}

	return zlstate, nil
}

func (s *S3Backend) PatchComponent(key string, component, status string) (*ZLState, error) {
	zlState, err := s.Get(key)
	if err != nil {
		return nil, errors.Wrap(err, "error getting zLstate from remote backend")
	}

	var oldStatus string
	updated := false
	for _, c := range zlState.Components {
		if c.Name != component {
			continue
		}
		oldStatus = c.Status
		c.Status = status
		c.UpdatedAt = time.Now().UTC()
		zlState.UpdatedAt = time.Now().UTC()
		updated = true
		break
	}
	if !updated {
		return nil, errors.Errorf("component not found: %s", component)
	}

	if err := s.Put(key, zlState, true); err != nil {
		return nil, errors.Wrap(err, "error persisting zLstate to remote backend")
	}

	msg := fmt.Sprintf("updated environment component [%s] status from [%s] to [%s]", component, oldStatus, status)
	s.log.WithFields(logrus.Fields{
		"company":     zlState.Company,
		"team":        zlState.Team,
		"environment": zlState.Environment,
		"component":   component,
	}).Info(msg)

	return zlState, nil
}

func (s *S3Backend) Put(key string, state *ZLState, force bool) error {
	s.log.WithFields(logrus.Fields{
		"key":     key,
		"zLstate": state,
	}).Info("Putting zLstate to remote backend [s3]")
	if !force {
		exists, err := s.exists(key)
		if err != nil {
			return errors.Wrapf(err, "error checking does object already exist for key: [%s]", key)
		}
		if exists {
			s.log.WithField("key", key).Info("State already exists, returning early")
			return errors.WithStack(ErrKeyAlreadyExists)
		}
	}

	addDefaults(state)

	input := s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(util.ToJSONBytes(state, false)),
	}

	_, err := s.s3.PutObject(s.ctx, &input)
	if err != nil {
		return errors.Wrap(err, "error saving zLstate")
	}

	s.log.WithFields(logrus.Fields{
		"key":     key,
		"zLstate": state,
	}).Info("zLstate persisted successfully to remote backend [s3]")

	return nil
}

func (s *S3Backend) UpsertComponent(key string, component *Component) (*ZLState, error) {
	zlstate, err := s.Get(key)
	if err != nil {
		return nil, errors.Wrapf(err, "error getting zlstate with key [%s]", key)
	}

	patched := false
	for i, c := range zlstate.Components {
		if c.Name == component.Name {
			zlstate.Components[i] = component
			patched = true
		}
	}

	if !patched {
		zlstate.Components = append(zlstate.Components, component)
	}

	if err := s.Put(key, zlstate, true); err != nil {
		return nil, errors.Wrapf(err, "error updating zlstate for key [%s]", key)
	}

	return zlstate, nil
}

func (s *S3Backend) DeleteComponent(key, component string) (*ZLState, error) {
	zlstate, err := s.Get(key)
	if err != nil {
		return nil, errors.Wrapf(err, "error getting zlstate with key [%s]", key)
	}

	for i, c := range zlstate.Components {
		if c.Name == component {
			zlstate.Components = append(zlstate.Components[:i], zlstate.Components[i+1:]...)
		}
	}

	if err := s.Put(key, zlstate, true); err != nil {
		return nil, errors.Wrapf(err, "error updating zlstate for key [%s]", key)
	}

	return zlstate, nil
}

var _ Backend = (*S3Backend)(nil)

func (s *S3Backend) exists(key string) (bool, error) {
	_, err := s.s3.HeadObject(s.ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		var responseError *awshttp.ResponseError
		if errors.As(err, &responseError) && responseError.ResponseError.HTTPStatusCode() == http.StatusNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func addDefaults(state *ZLState) {
	if state.CreatedAt.IsZero() {
		state.CreatedAt = time.Now().UTC()
	}
	if state.UpdatedAt.IsZero() {
		state.UpdatedAt = time.Now().UTC()
	}

	for _, c := range state.Components {
		if c.CreatedAt.IsZero() {
			c.CreatedAt = time.Now().UTC()
		}
		if c.UpdatedAt.IsZero() {
			c.UpdatedAt = time.Now().UTC()
		}
	}
}
