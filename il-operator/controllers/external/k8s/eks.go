package k8s

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/credentials"
	secrets2 "github.com/compuzest/zlifecycle-il-operator/controllers/codegen/secrets"
	"github.com/compuzest/zlifecycle-il-operator/controllers/external/secrets"
	"github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/pkg/errors"
	"sigs.k8s.io/aws-iam-authenticator/pkg/token"
)

type EKS struct {
	ctx          context.Context
	secretClient secrets.API
	eksClient    *eks.Client
	secretMeta   *secrets2.Meta
	creds        *secrets.AWSCreds
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

	p := credentials.NewStaticCredentialsProvider(creds.AccessKeyID, creds.SecretAccessKey, creds.SessionToken)
	cfg, err := config.LoadDefaultConfig(ctx, config.WithCredentialsProvider(p))
	if err != nil {
		return nil, errors.Wrap(err, "error loading config with static credentials provider")
	}

	return &EKS{
		ctx:        ctx,
		eksClient:  eks.NewFromConfig(cfg),
		creds:      creds,
		secretMeta: secretMeta,
		log:        log,
	}, nil
}

func (e *EKS) init() error {
	creds, err := secrets.GetAWSCreds(e.secretClient, e.secretMeta, e.log)
	if err != nil {
		return errors.Wrapf(err, "error fetching AWS creds for eks")
	}

	p := credentials.NewStaticCredentialsProvider(creds.AccessKeyID, creds.SecretAccessKey, creds.SessionToken)
	cfg, err := config.LoadDefaultConfig(e.ctx, config.WithCredentialsProvider(p))
	if err != nil {
		return errors.Wrap(err, "error loading config with static credentials provider")
	}

	e.eksClient = eks.NewFromConfig(cfg)
	return nil
}

func (e *EKS) DescribeCluster(name string) (*ClusterInfo, error) {
	if e.eksClient == nil {
		if err := e.init(); err != nil {
			return nil, errors.Wrap(err, "error initializing eks client")
		}
	}

	input := &eks.DescribeClusterInput{
		Name: aws.String(name),
	}
	info, err := e.eksClient.DescribeCluster(e.ctx, input)
	if err != nil {
		return nil, errors.Wrapf(err, "error describing cluster %s", name)
	}

	gen, err := token.NewGenerator(true, false)
	if err != nil {
		return nil, errors.Wrapf(err, "error creating aws iam authenticator as token generator for cluster %s", name)
	}
	opts := &token.GetTokenOptions{
		ClusterID: *info.Cluster.Name,
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
