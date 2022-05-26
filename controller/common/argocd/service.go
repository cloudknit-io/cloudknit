package argocd

import (
	appv1 "github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	"net/http"
)

//nolint
//go:generate mockgen --build_flags=--mod=mod -destination=./mock_argocd_api.go -package=argocd "github.com/compuzest/zlifecycle-il-operator/controller/common/argocd" API

type API interface {
	GetAuthToken() (*GetTokenResponse, error)
	ListRepositories(bearerToken string) (*RepositoryList, *http.Response, error)
	CreateRepository(body interface{}, bearerToken string) (*http.Response, error)
	CreateApplication(application *appv1.Application, bearerToken string) (*http.Response, error)
	DeleteApplication(name string, bearerToken string) error
	DoesApplicationExist(name string, bearerToken string) (exists bool, err error)
	CreateProject(project *CreateProjectBody, bearerToken string) (*http.Response, error)
	DoesProjectExist(name string, bearerToken string) (exists bool, response *http.Response, err error)
	ListClusters(name *string, bearerToken string) (*ListClustersResponse, error)
	RegisterCluster(body *RegisterClusterBody, bearerToken string) (*http.Response, error)
	UpdateCluster(clusterURL string, body *UpdateClusterBody, updatedFields []string, bearerToken string) (*http.Response, error)
}
