package awsssm

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/compuzest/zlifecycle-il-operator/controller/external/aws/awscfg"
	"github.com/compuzest/zlifecycle-il-operator/controller/external/secrets"
	"github.com/pkg/errors"
)

type SSM struct {
	cl        awscfg.ConfigLoader
	ssmClient *ssm.Client
}

func LazyLoadSSM(cl awscfg.ConfigLoader) *SSM {
	return &SSM{
		cl: cl,
	}
}

func NewSSM(ctx context.Context, cl awscfg.ConfigLoader) (*SSM, error) {
	cfg, err := cl.LoadConfig(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "error loading aws config")
	}

	return &SSM{
		cl:        cl,
		ssmClient: ssm.NewFromConfig(cfg),
	}, nil
}

func (s *SSM) init(ctx context.Context) error {
	cfg, err := s.cl.LoadConfig(ctx)
	if err != nil {
		return errors.Wrap(err, "error loading default aws config using static credentials provider")
	}

	s.ssmClient = ssm.NewFromConfig(cfg)
	return nil
}

func (s *SSM) GetSecret(ctx context.Context, key string) (*secrets.Secret, error) {
	if s.ssmClient == nil {
		if err := s.init(ctx); err != nil {
			return nil, errors.Wrap(err, "error initializing ssm client")
		}
	}

	input := ssm.GetParameterInput{
		Name:           &key,
		WithDecryption: true,
	}
	output, err := s.ssmClient.GetParameter(ctx, &input)
	if err != nil {
		return nil, errors.Wrapf(err, "error getting parameter %s from ssm", key)
	}

	return &secrets.Secret{Value: output.Parameter.Value, Key: *output.Parameter.Name, Exists: output.Parameter.Value != nil}, nil
}

func (s *SSM) GetSecrets(ctx context.Context, keys ...string) ([]*secrets.Secret, error) {
	if s.ssmClient == nil {
		if err := s.init(ctx); err != nil {
			return nil, errors.Wrap(err, "error initializing ssm client")
		}
	}

	input := ssm.GetParametersInput{
		Names:          keys,
		WithDecryption: true,
	}
	output, err := s.ssmClient.GetParameters(ctx, &input)
	if err != nil {
		return nil, errors.Wrapf(err, "error getting parameters %s from ssm", keys)
	}

	scrts := make([]*secrets.Secret, 0, len(keys))

	for _, s := range output.Parameters {
		scrt := &secrets.Secret{Value: s.Value, Key: *s.Name, Exists: s.Value != nil}
		scrts = append(scrts, scrt)
	}

	return scrts, nil
}

var _ secrets.API = (*SSM)(nil)
