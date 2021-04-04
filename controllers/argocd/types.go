package argocd

import (
	"net/http"

	appv1 "github.com/argoproj/argo-cd/pkg/apis/application/v1alpha1"
	"github.com/go-logr/logr"
)

type Api interface {
	GetAuthToken() (*GetTokenResponse, error)
	ListRepositories(bearerToken string) (*RepositoryList, *http.Response, error)
	CreateRepository(body CreateRepoBody, bearerToken string) (*http.Response, error)
	CreateApplication(application *appv1.Application, bearerToken string) (*http.Response, error)
	DoesApplicationExist(name string, bearerToken string) (bool, error)
}

type HttpApi struct {
	Api
	ServerUrl string
	Log       logr.Logger
}

type MockAPI struct{}

type Credentials struct {
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
