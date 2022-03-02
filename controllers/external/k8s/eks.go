package k8s

import (
	"context"

	credentialsv2 "github.com/aws/aws-sdk-go-v2/credentials"
	awsv1 "github.com/aws/aws-sdk-go/aws"
	credentialsv1 "github.com/aws/aws-sdk-go/aws/credentials"
	sessionv1 "github.com/aws/aws-sdk-go/aws/session"
	secrets2 "github.com/compuzest/zlifecycle-il-operator/controllers/codegen/secrets"
	"github.com/compuzest/zlifecycle-il-operator/controllers/external/secrets"
	"github.com/sirupsen/logrus"

	awsv2 "github.com/aws/aws-sdk-go-v2/aws"
	configv2 "github.com/aws/aws-sdk-go-v2/config"
	eksv2 "github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/pkg/errors"
	"sigs.k8s.io/aws-iam-authenticator/pkg/token"
)

type EKS struct {
	ctx          context.Context
	secretClient secrets.API
	eksClient    *eksv2.Client
	secretMeta   *secrets2.Meta
	creds        *secrets.AWSCreds
	cfg          *awsv2.Config
	log          *logrus.Entry
}

func LazyLoadEKS(ctx context.Context, secretClient secrets.API, secretMeta *secrets2.Meta, log *logrus.Entry) *EKS {
	return &EKS{
		ctx:          ctx,
		secretMeta:   secretMeta,
		secretClient: secretClient,
		log:          log,
	}
}

func NewEKS(ctx context.Context, secretClient secrets.API, secretMeta *secrets2.Meta, log *logrus.Entry) (*EKS, error) {
	creds, err := secrets.GetAWSCreds(secretClient, secretMeta, log)
	if err != nil {
		return nil, errors.Wrapf(err, "error fetching AWS creds for eks")
	}

	p := credentialsv2.NewStaticCredentialsProvider(creds.AccessKeyID, creds.SecretAccessKey, creds.SessionToken)
	cfg, err := configv2.LoadDefaultConfig(ctx, configv2.WithCredentialsProvider(p))
	if err != nil {
		return nil, errors.Wrap(err, "error loading config with static credentials provider")
	}

	return &EKS{
		ctx:        ctx,
		eksClient:  eksv2.NewFromConfig(cfg),
		creds:      creds,
		cfg:        &cfg,
		secretMeta: secretMeta,
		log:        log,
	}, nil
}

func (e *EKS) init() error {
	creds, err := secrets.GetAWSCreds(e.secretClient, e.secretMeta, e.log)
	if err != nil {
		return errors.Wrapf(err, "error fetching AWS creds for eks")
	}

	p := credentialsv2.NewStaticCredentialsProvider(creds.AccessKeyID, creds.SecretAccessKey, creds.SessionToken)
	cfg, err := configv2.LoadDefaultConfig(e.ctx, configv2.WithCredentialsProvider(p))
	if err != nil {
		return errors.Wrap(err, "error loading config with static credentials provider")
	}

	e.creds = creds
	e.cfg = &cfg
	e.eksClient = eksv2.NewFromConfig(cfg)
	return nil
}

func (e *EKS) DescribeCluster(name string) (*ClusterInfo, error) {
	if e.eksClient == nil {
		if err := e.init(); err != nil {
			return nil, errors.Wrap(err, "error initializing eks client")
		}
	}

	input := &eksv2.DescribeClusterInput{
		Name: awsv2.String(name),
	}
	info, err := e.eksClient.DescribeCluster(e.ctx, input)
	if err != nil {
		return nil, errors.Wrapf(err, "error describing cluster %s", name)
	}

	gen, err := token.NewGenerator(true, false)
	if err != nil {
		return nil, errors.Wrapf(err, "error creating aws iam authenticator as token generator for cluster %s", name)
	}

	creds := credentialsv1.NewStaticCredentials(e.creds.AccessKeyID, e.creds.SecretAccessKey, e.creds.SessionToken)
	sess, err := sessionv1.NewSession(&awsv1.Config{Credentials: creds})
	if err != nil {
		return nil, errors.Wrap(err, "error creating new aws session")
	}
	opts := &token.GetTokenOptions{
		ClusterID: *info.Cluster.Name,
		Session:   sess,
	}
	tok, err := gen.GetWithOptions(opts)
	if err != nil {
		return nil, errors.Wrapf(err, "error generating token using aws iam authenticator for cluster %s", name)
	}
	return &ClusterInfo{
		Name:                 *info.Cluster.Name,
		Version:              *info.Cluster.Version,
		CertificateAuthority: *info.Cluster.CertificateAuthority.Data,
		Endpoint:             *info.Cluster.Endpoint,
		BearerToken:          tok.Token,
	}, nil
}

var _ API = (*EKS)(nil)
