package k8s

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
