package argocd

import (
	"github.com/go-logr/logr"
	"net/http"
)

type ArgocdAPI interface {
	GetAuthToken() (*GetTokenResponse, error)
	ListRepositories(bearerToken string) (*RepositoryList, *http.Response, error)
	CreateRepository(body CreateRepoBody, bearerToken string) (*http.Response, error)
}

type ArgocdHttpAPI struct {
	ServerUrl string
	Log logr.Logger
}

type ArgocdMockAPI struct {}

type ArgocdCredentials struct {
	Username string
	Password string
}

type GetTokenBody struct {
	Username string
	Password string
}

type GetTokenResponse struct {
	Token string
}

type RepoOpts struct {
	RepoUrl       string
	SshPrivateKey string
}

type CreateRepoBody struct {
	Repo          string `json:"repo"`
	Name          string `json:"name"`
	SshPrivateKey string `json:"sshPrivateKey"`
}

type RepositoryList struct {
	Items []Repository `json:"items"`
}

type Repository struct {
	Repo string `json:"repo"`
	Name string `json:"name"`
}
