package awscfg

import (
	"context"

	awsv2 "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	secrets2 "github.com/compuzest/zlifecycle-il-operator/controller/codegen/secret"
	"github.com/compuzest/zlifecycle-il-operator/controller/env"
	"github.com/compuzest/zlifecycle-il-operator/controller/external/secrets"
	"github.com/compuzest/zlifecycle-il-operator/controller/util"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	kClient "sigs.k8s.io/controller-runtime/pkg/client"
)

type ConfigLoader interface {
	LoadConfig(ctx context.Context) (awsv2.Config, error)
}

type SSMCredentialsLoader struct {
	s secrets.API
	l *logrus.Entry
	i *secrets2.Identifier
}

func (l *SSMCredentialsLoader) LoadConfig(ctx context.Context) (awsv2.Config, error) {
	creds, err := secrets.GetAWSCredentials(ctx, l.s, l.i, l.l)
	if err != nil {
		return awsv2.Config{}, errors.Wrapf(err, "error fetching AWS creds for eks")
	}

	p := credentials.NewStaticCredentialsProvider(creds.AccessKeyID, creds.SecretAccessKey, creds.SessionToken)
	return config.LoadDefaultConfig(ctx, config.WithCredentialsProvider(p))
}

func NewSSMCredentialsLoader(secretsClient secrets.API, id *secrets2.Identifier, logger *logrus.Entry) *SSMCredentialsLoader {
	return &SSMCredentialsLoader{
		s: secretsClient,
		i: id,
		l: logger,
	}
}

type K8sSecretCredentialsLoader struct {
	kc     kClient.Client
	secret string
}

func (l *K8sSecretCredentialsLoader) LoadConfig(ctx context.Context) (awsv2.Config, error) {
	creds, err := getCredentialsFromSecret(ctx, l.secret, l.kc)
	if err != nil {
		return awsv2.Config{}, errors.Wrap(err, "error getting AWS credentials from shared aws credentials secret")
	}
	p := credentials.NewStaticCredentialsProvider(creds.AccessKeyID, creds.SecretAccessKey, creds.SessionToken)
	cfg, err := config.LoadDefaultConfig(ctx, config.WithCredentialsProvider(p))
	if err != nil {
		return awsv2.Config{}, errors.Wrap(err, "error loading default aws config using static credentials provider")
	}

	return cfg, nil
}

func getCredentialsFromSecret(ctx context.Context, secret string, client kClient.Client) (*secrets.AWSCredentials, error) {
	var credsSecret v1.Secret
	key := kClient.ObjectKey{Name: secret, Namespace: env.ExecutorNamespace()}
	if err := client.Get(ctx, key, &credsSecret); err != nil {
		return nil, errors.Wrap(err, "error getting shared aws secret")
	}

	accessKeyID := string(credsSecret.Data[util.AWSAccessKeyID])
	secretAccessKey := string(credsSecret.Data[util.AWSSecretAccessKey])

	if accessKeyID == "" || secretAccessKey == "" {
		return nil, errors.New("missing AWS Access Key ID and/or AWS Secret Access key in shared aws secret")
	}

	return &secrets.AWSCredentials{AccessKeyID: accessKeyID, SecretAccessKey: secretAccessKey}, nil
}

func NewK8sSecretCredentialsLoader(kc kClient.Client, secret string) *K8sSecretCredentialsLoader {
	return &K8sSecretCredentialsLoader{kc: kc, secret: secret}
}

type DefaultCredentialsLoader struct{}

func (l *DefaultCredentialsLoader) LoadConfig(ctx context.Context) (awsv2.Config, error) {
	return config.LoadDefaultConfig(ctx)
}

func NewDefaultCredentialsLoader() *DefaultCredentialsLoader {
	return &DefaultCredentialsLoader{}
}
