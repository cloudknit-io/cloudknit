package cloudknitservice

import (
	"context"
	"net/http"

	stablev1 "github.com/compuzest/zlifecycle-il-operator/api/v1"
	"github.com/compuzest/zlifecycle-il-operator/controller/util"
	"github.com/sirupsen/logrus"
)

type API interface {
	GetOrganization(ctx context.Context, organizationName string, log *logrus.Entry) (*Organization, error)
	PostError(ctx context.Context, organizationName string, environment *stablev1.Environment, allErrs []string, log *logrus.Entry) error
}

type Service struct {
	host       string
	httpClient *http.Client
}

func NewService(host string) *Service {
	return &Service{
		host:       host,
		httpClient: util.GetHTTPClient(),
	}
}
