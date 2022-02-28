package k8s

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/pkg/errors"
	"sigs.k8s.io/aws-iam-authenticator/pkg/token"
)

type EKS struct {
	ctx    context.Context
	client *eks.Client
}

func NewEKS(ctx context.Context) (*EKS, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "error loading default config")
	}

	return &EKS{
		ctx:    ctx,
		client: eks.NewFromConfig(cfg),
	}, nil
}

func (e *EKS) DescribeCluster(name string) (*ClusterInfo, error) {
	input := &eks.DescribeClusterInput{
		Name: aws.String(name),
	}
	info, err := e.client.DescribeCluster(e.ctx, input)
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
