package k8s

//go:generate mockgen --build_flags=--mod=mod -destination=./mock_k8s_api.go -package=k8s "github.com/compuzest/zlifecycle-il-operator/controllers/external/k8s" API
type API interface {
	DescribeCluster(name string) (*ClusterInfo, error)
}

type ClusterInfo struct {
	Name                 string `json:"name"`
	Version              string `json:"version"`
	CertificateAuthority string `json:"certificateAuthority"`
	Endpoint             string `json:"endpoint"`
	BearerToken          string `json:"bearerToken"`
}
