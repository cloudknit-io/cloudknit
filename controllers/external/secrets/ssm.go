package secrets

import (
	"context"

	"github.com/compuzest/zlifecycle-il-operator/controllers/util"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/compuzest/zlifecycle-il-operator/controllers/env"
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	kClient "sigs.k8s.io/controller-runtime/pkg/client"
)

type SSM struct {
	ctx       context.Context
	k8sClient kClient.Client
	ssmClient *ssm.Client
}

func LazyLoadSSM(ctx context.Context, client kClient.Client) *SSM {
	return &SSM{
		ctx:       ctx,
		k8sClient: client,
	}
}

func NewSSM(ctx context.Context, client kClient.Client) (*SSM, error) {
	creds, err := getCreds(ctx, client)
	if err != nil {
		return nil, errors.Wrap(err, "error getting AWS creds from shared aws creds secret")
	}

	p := credentials.NewStaticCredentialsProvider(creds.AccessKeyID, creds.SecretAccessKey, creds.SessionToken)
	cfg, err := config.LoadDefaultConfig(ctx, config.WithCredentialsProvider(p))
	if err != nil {
		return nil, errors.Wrap(err, "error loading default aws config using static credentials provider")
	}

	return &SSM{
		ctx:       ctx,
		k8sClient: client,
		ssmClient: ssm.NewFromConfig(cfg),
	}, nil
}

func (s *SSM) init() error {
	creds, err := getCreds(s.ctx, s.k8sClient)
	if err != nil {
		return errors.Wrap(err, "error getting AWS creds from shared aws creds secret")
	}

	p := credentials.NewStaticCredentialsProvider(creds.AccessKeyID, creds.SecretAccessKey, creds.SessionToken)
	cfg, err := config.LoadDefaultConfig(s.ctx, config.WithCredentialsProvider(p))
	if err != nil {
		return errors.Wrap(err, "error loading default aws config using static credentials provider")
	}

	s.ssmClient = ssm.NewFromConfig(cfg)
	return nil
}

func getCreds(ctx context.Context, client kClient.Client) (*AWSCredentials, error) {
	var credsSecret v1.Secret
	key := kClient.ObjectKey{Name: env.Config.SharedAWSCredsSecret, Namespace: env.ExecutorNamespace()}
	if err := client.Get(ctx, key, &credsSecret); err != nil {
		return nil, errors.Wrap(err, "error getting shared aws secret")
	}

	accessKeyID := string(credsSecret.Data[util.AWSAccessKeyID])
	secretAccessKey := string(credsSecret.Data[util.AWSSecretAccessKey])

	if accessKeyID == "" || secretAccessKey == "" {
		return nil, errors.New("missing AWS Access Key ID and/or AWS Secret Access key in shared aws secret")
	}

	return &AWSCredentials{AccessKeyID: accessKeyID, SecretAccessKey: secretAccessKey}, nil
}

func (s *SSM) GetSecret(key string) (*Secret, error) {
	if s.ssmClient == nil {
		if err := s.init(); err != nil {
			return nil, errors.Wrap(err, "error initializing ssm client")
		}
	}

	input := ssm.GetParameterInput{
		Name:           &key,
		WithDecryption: true,
	}
	output, err := s.ssmClient.GetParameter(s.ctx, &input)
	if err != nil {
		return nil, errors.Wrapf(err, "error getting parameter %s from ssm", key)
	}

	return &Secret{Value: output.Parameter.Value, Key: *output.Parameter.Name, Exists: output.Parameter.Value != nil}, nil
}

func (s *SSM) GetSecrets(keys ...string) ([]*Secret, error) {
	if s.ssmClient == nil {
		if err := s.init(); err != nil {
			return nil, errors.Wrap(err, "error initializing ssm client")
		}
	}

	input := ssm.GetParametersInput{
		Names:          keys,
		WithDecryption: true,
	}
	output, err := s.ssmClient.GetParameters(s.ctx, &input)
	if err != nil {
		return nil, errors.Wrapf(err, "error getting parameters %s from ssm", keys)
	}

	secrets := make([]*Secret, 0, len(keys))

	for _, scrt := range output.Parameters {
		secret := &Secret{Value: scrt.Value, Key: *scrt.Name, Exists: scrt.Value != nil}
		secrets = append(secrets, secret)
	}

	return secrets, nil
}

var _ API = (*SSM)(nil)
