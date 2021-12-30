package zlstate

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/compuzest/zlifecycle-state-manager/app/env"
	"github.com/compuzest/zlifecycle-state-manager/app/util"
	"github.com/pkg/errors"
	"time"
)

type S3Backend struct {
	bucket string
	s3     *s3.S3
}

func NewS3Backend(bucket string) *S3Backend {
	mySession := session.Must(session.NewSession())
	return &S3Backend{
		bucket: bucket,
		s3:     s3.New(mySession, aws.NewConfig().WithRegion(env.Config().AWSRegion)),
	}
}

// Get returns the state file whose key is the path in the bucket for which the backend was created
func (s *S3Backend) Get(key string) (*ZLState, error) {
	input := s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}
	output, err := s.s3.GetObject(&input)
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

func (s *S3Backend) Put(key string, state *ZLState) error {
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

	_, err = s.s3.PutObject(&input)
	if err != nil {
		return errors.Wrap(err, "error saving zLstate")
	}

	return nil
}

var _ Backend = (*S3Backend)(nil)

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
