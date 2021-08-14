package argocd

import (
	appv1 "github.com/argoproj/argo-cd/pkg/apis/application/v1alpha1"
	"github.com/go-logr/logr"
	_ "github.com/golang/mock/mockgen/model"
	"net/http"
)

//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=../../mocks/mock_argocd_api.go -package=mocks "github.com/compuzest/zlifecycle-il-operator/controllers/argocd" Api

type Api interface {
	GetAuthToken() (*GetTokenResponse, error)
	ListRepositories(bearerToken string) (*RepositoryList, *http.Response, error)
	CreateRepository(body CreateRepoBody, bearerToken string) (*http.Response, error)
	CreateApplication(application *appv1.Application, bearerToken string) (*http.Response, error)
	DeleteApplication(name string, bearerToken string) error
	DoesApplicationExist(name string, bearerToken string) (bool, error)
}

type HttpApi struct {
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
