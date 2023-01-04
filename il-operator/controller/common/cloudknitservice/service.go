package cloudknitservice

import (
	"context"
	"net/http"

	"github.com/compuzest/zlifecycle-il-operator/controller/util"
	"github.com/sirupsen/logrus"
)

type API interface {
	GetOrganization(ctx context.Context, organizationName string, log *logrus.Entry) (*Organization, error)
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
