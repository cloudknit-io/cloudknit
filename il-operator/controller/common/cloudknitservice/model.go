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
	ID                int       `json:"id"`
	Name              string    `json:"name"`
	GitHubRepo        string    `json:"githubRepo"`
	GitHubOrgName     string    `json:"githubOrgName"`
	TermsAgreedUserId int       `json:"termsAgreedUserId"`
	Provisioned       bool      `json:"provisioned"`
	Created           time.Time `json:"created"`
	Updated           time.Time `json:"updated"`
}
