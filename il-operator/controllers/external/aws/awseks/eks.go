package awseks

import (
	"context"
	awsv1 "github.com/aws/aws-sdk-go/aws"
	credentialsv1 "github.com/aws/aws-sdk-go/aws/credentials"
	sessionv1 "github.com/aws/aws-sdk-go/aws/session"
	"github.com/compuzest/zlifecycle-il-operator/controllers/external/aws/awscfg"
	"github.com/sirupsen/logrus"

	awsv2 "github.com/aws/aws-sdk-go-v2/aws"
	eksv2 "github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/pkg/errors"
	"sigs.k8s.io/aws-iam-authenticator/pkg/token"
)

type EKS struct {
	ctx context.Context
	cl  awscfg.ConfigLoader
	ec  *eksv2.Client
	cfg *awsv2.Config
	log *logrus.Entry
}

func LazyLoadEKS(ctx context.Context, cl awscfg.ConfigLoader, log *logrus.Entry) *EKS {
	return &EKS{
		ctx: ctx,
		cl:  cl,
		log: log,
	}
}

func NewEKS(ctx context.Context, cl awscfg.ConfigLoader, log *logrus.Entry) (*EKS, error) {
	cfg, err := cl.LoadConfig(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "error loading config")
	}

	return &EKS{
		ctx: ctx,
		ec:  eksv2.NewFromConfig(cfg),
		cl:  cl,
		cfg: &cfg,
		log: log,
	}, nil
}

func (e *EKS) init(ctx context.Context) error {
	cfg, err := e.cl.LoadConfig(ctx)
	if err != nil {
		return errors.Wrap(err, "error loading config with static credentials provider")
	}

	e.cfg = &cfg
	e.ec = eksv2.NewFromConfig(cfg)
	return nil
}

func (e *EKS) DescribeCluster(ctx context.Context, name string) (*ClusterInfo, error) {
	if e.ec == nil {
		if err := e.init(ctx); err != nil {
			return nil, errors.Wrap(err, "error initializing eks client")
		}
	}

	input := &eksv2.DescribeClusterInput{
		Name: awsv2.String(name),
	}
	info, err := e.ec.DescribeCluster(e.ctx, input)
	if err != nil {
		return nil, errors.Wrapf(err, "error describing cluster %s", name)
	}

	sess, err := e.newAWSSession(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "error creating new aws session")
	}

	tok, err := e.newEKSToken(*info.Cluster.Name, sess)
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

func (e *EKS) newAWSSession(ctx context.Context) (*sessionv1.Session, error) {
	credsv2, err := e.cfg.Credentials.Retrieve(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "error retrieving credentials from aws v2 config")
	}
	credsv1 := credentialsv1.NewStaticCredentials(credsv2.AccessKeyID, credsv2.SecretAccessKey, credsv2.SessionToken)
	sess, err := sessionv1.NewSession(&awsv1.Config{Credentials: credsv1})
	if err != nil {
		return nil, errors.Wrap(err, "error creating new aws session")
	}

	return sess, nil
}

func (e *EKS) newEKSToken(cluster string, session *sessionv1.Session) (*token.Token, error) {
	gen, err := token.NewGenerator(true, false)
	if err != nil {
		return nil, errors.Wrapf(err, "error creating aws iam authenticator as token generator for cluster %s", cluster)
	}
	opts := &token.GetTokenOptions{
		ClusterID: cluster,
		Session:   session,
	}
	tok, err := gen.GetWithOptions(opts)
	if err != nil {
		return nil, errors.Wrapf(err, "error generating token using aws iam authenticator for cluster %s", cluster)
	}

	return &tok, nil
}

var _ API = (*EKS)(nil)
