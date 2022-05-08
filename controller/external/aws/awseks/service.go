package awseks

import "context"

//go:generate mockgen --build_flags=--mod=mod -destination=./mock_k8s_api.go -package=k8s "github.com/compuzest/zlifecycle-il-operator/controller/external/aws/awseks" API
type API interface {
	DescribeCluster(ctx context.Context, name string) (*ClusterInfo, error)
}

type ClusterInfo struct {
	Name                 string `json:"name"`
	Version              string `json:"version"`
	CertificateAuthority string `json:"certificateAuthority"`
	Endpoint             string `json:"endpoint"`
	BearerToken          string `json:"bearerToken"`
}
