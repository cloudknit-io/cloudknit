package secrets

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/compuzest/zlifecycle-il-operator/controllers/env"
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	kClient "sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	keyAWSAccessKeyID     = "aws_access_key_id"
	keyAWSSecretAccessKey = "aws_secret_access_key"
)

type SSM struct {
	ctx    context.Context
	client *ssm.Client
}

func NewSSM(ctx context.Context, client kClient.Client) (*SSM, error) {
	var credsSecret v1.Secret
	key := kClient.ObjectKey{Name: env.Config.SharedAWSCredsSecret, Namespace: env.ExecutorNamespace()}
	if err := client.Get(ctx, key, &credsSecret); err != nil {
		return nil, errors.Wrap(err, "error getting shared aws secret")
	}

	accessKeyID := string(credsSecret.Data[keyAWSAccessKeyID])
	secretAccessKey := string(credsSecret.Data[keyAWSSecretAccessKey])

	p := credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, "")
	cfg, err := config.LoadDefaultConfig(ctx, config.WithCredentialsProvider(p))
	if err != nil {
		return nil, errors.Wrap(err, "error loading default aws config")
	}
	ssmClient := ssm.NewFromConfig(cfg)

	return &SSM{
		ctx:    ctx,
		client: ssmClient,
	}, nil
}

func (s *SSM) GetSecret(key string) (*Secret, error) {
	input := ssm.GetParameterInput{
		Name:           &key,
		WithDecryption: true,
	}
	output, err := s.client.GetParameter(s.ctx, &input)
	if err != nil {
		return nil, errors.Wrapf(err, "error getting parameter %s from SSM", key)
	}

	return &Secret{Value: output.Parameter.Value, Key: *output.Parameter.Name, Exists: output.Parameter.Value != nil}, nil
}

func (s *SSM) GetSecrets(keys ...string) ([]*Secret, error) {
	input := ssm.GetParametersInput{
		Names:          keys,
		WithDecryption: true,
	}
	output, err := s.client.GetParameters(s.ctx, &input)
	if err != nil {
		return nil, errors.Wrapf(err, "error getting parameters %s from SSM", keys)
	}

	secrets := make([]*Secret, 0, len(keys))

	for _, scrt := range output.Parameters {
		secret := &Secret{Value: scrt.Value, Key: *scrt.Name, Exists: scrt.Value != nil}
		secrets = append(secrets, secret)
	}

	return secrets, nil
}

var _ API = (*SSM)(nil)
