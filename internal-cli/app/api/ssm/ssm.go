package secrets

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/credentials"

	aws2 "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/compuzest/zlifecycle-internal-cli/app/api/aws"

	"github.com/compuzest/zlifecycle-internal-cli/app/lib/secret"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/pkg/errors"
)

type SSM struct {
	ctx       context.Context
	ssmClient *ssm.Client
}

func NewSSM(ctx context.Context, auth *aws.Auth) (*SSM, error) {
	cfg, err := newConfig(ctx, auth)
	if err != nil {
		return nil, errors.Wrap(err, "error loading default aws config using static credentials provider")
	}

	return &SSM{
		ctx:       ctx,
		ssmClient: ssm.NewFromConfig(cfg),
	}, nil
}

func newConfig(ctx context.Context, auth *aws.Auth) (aws2.Config, error) {
	if auth == nil {
		return aws2.Config{}, errors.New("auth not provided")
	}
	switch auth.Mode {
	case aws.AuthModeProfile:
		if auth.Profile == "" {
			return aws2.Config{}, errors.New("profile not provided")
		}
		return config.LoadDefaultConfig(ctx, config.WithSharedConfigProfile(auth.Profile), config.WithRegion(auth.Region))
	case aws.AuthModeStatic:
		if auth.AccessKeyID == "" {
			return aws2.Config{}, errors.New("aws access key id not provided")
		}
		if auth.SecretAccessKey == "" {
			return aws2.Config{}, errors.New("aws secret access key not provided")
		}
		loader := config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(auth.AccessKeyID, auth.SecretAccessKey, ""))
		return config.LoadDefaultConfig(ctx, loader)
	default:
		return config.LoadDefaultConfig(ctx, config.WithRegion(auth.Region))
	}
}

func (s *SSM) GetSecret(key string) (*secret.Secret, error) {
	input := ssm.GetParameterInput{
		Name:           &key,
		WithDecryption: true,
	}
	output, err := s.ssmClient.GetParameter(s.ctx, &input)
	if err != nil {
		return nil, errors.Wrapf(err, "error getting parameter %s from ssm", key)
	}

	return &secret.Secret{Value: output.Parameter.Value, Key: *output.Parameter.Name, Exists: output.Parameter.Value != nil}, nil
}

func (s *SSM) GetSecrets(keys ...string) ([]*secret.Secret, error) {
	input := ssm.GetParametersInput{
		Names:          keys,
		WithDecryption: true,
	}
	output, err := s.ssmClient.GetParameters(s.ctx, &input)
	if err != nil {
		return nil, errors.Wrapf(err, "error getting parameters %s from ssm", keys)
	}

	scrts := make([]*secret.Secret, 0, len(keys))

	for _, scrt := range output.Parameters {
		s := &secret.Secret{Value: scrt.Value, Key: *scrt.Name, Exists: scrt.Value != nil}
		scrts = append(scrts, s)
	}

	return scrts, nil
}

var _ secret.API = (*SSM)(nil)
