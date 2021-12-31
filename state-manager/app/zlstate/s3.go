package zlstate

import (
	"bytes"
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/compuzest/zlifecycle-state-manager/app/util"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type S3Backend struct {
	ctx    context.Context
	log    *logrus.Entry
	bucket string
	s3     *s3.Client
}

func NewS3Backend(ctx context.Context, log *logrus.Entry, bucket string) (*S3Backend, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "error loading default aws config")
	}
	return &S3Backend{
		ctx:    ctx,
		log:    log,
		bucket: bucket,
		s3:     s3.NewFromConfig(cfg),
	}, nil
}

// Get returns the state file whose key is the path in the bucket for which the backend was created
func (s *S3Backend) Get(key string) (*ZLState, error) {
	s.log.WithField("key", key).Info("Getting zLstate from remote backend [s3]")
	exists, err := s.exists(key)
	if err != nil {
		return nil, errors.Wrapf(err, "error checking does object already exist for key: [%s]", key)
	}
	if !exists {
		return nil, errors.Errorf("key does not exist: [%s]", key)
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

var ErrKeyExists = errors.Errorf("object already exists")

func (s *S3Backend) Put(key string, state *ZLState) error {
	s.log.WithFields(logrus.Fields{
		"key":     key,
		"zLstate": state,
	}).Info("Putting zLstate to remote backend [s3]")
	exists, err := s.exists(key)
	if err != nil {
		return errors.Wrapf(err, "error checking does object already exist for key: [%s]", key)
	}
	if exists {
		s.log.WithField("key", key).Info("State already exists, returning early")
		return errors.WithStack(ErrKeyExists)
	}

	addDefaults(state)

	body, err := util.ToJSON(state)
	if err != nil {
		return errors.Wrap(err, "error serializing zLstate")
	}
	input := s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(body),
	}

	_, err = s.s3.PutObject(s.ctx, &input)
	if err != nil {
		return errors.Wrap(err, "error saving zLstate")
	}

	s.log.WithFields(logrus.Fields{
		"key":     key,
		"zLstate": state,
	}).Info("zLstate persisted successfully to remote backend [s3]")

	return nil
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
