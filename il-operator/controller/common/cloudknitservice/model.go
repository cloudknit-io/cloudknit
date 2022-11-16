package cloudknitservice

import (
	"time"
)

type GetOrganizationBody struct {
	OrganizationName string `json:"organizationName"`
}

type GetOrganizationResponse struct {
	Organization *Organization `json:"organization"`
}

type Organization struct {
	ID            int       `json:"id"`
	Name          string    `json:"name"`
	GitHubRepo    string    `json:"githubRepo"`
	GitHubOrgName string    `json:"githubOrgName"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	Provisioned   bool      `json:"provisioned"`
}
