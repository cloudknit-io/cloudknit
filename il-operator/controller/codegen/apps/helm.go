package apps

type HelmChart struct {
	APIVersion string `json:"apiVersion"`
	Name       string `json:"name"`
	Version    string `json:"version"`
}

func NewHelmChart(name string) *HelmChart {
	return &HelmChart{
		APIVersion: "v2",
		Name:       name,
		Version:    "1.0.0",
	}
}
