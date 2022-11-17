package cloudknitservice

type GetOrganizationBody struct {
	OrganizationName string `json:"organizationName"`
}

type Organization struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	GitHubRepo    string `json:"githubRepo"`
	GitHubOrgName string `json:"githubOrgName"`
	Provisioned   bool   `json:"provisioned"`
}
