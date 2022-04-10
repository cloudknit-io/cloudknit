package argocd

import (
	"context"
	"net/http"

	"github.com/compuzest/zlifecycle-internal-cli/app/common"
	"github.com/compuzest/zlifecycle-internal-cli/app/env"
	"github.com/sirupsen/logrus"
)

type API interface{}

type Service struct {
	ctx        context.Context
	log        *logrus.Entry
	host       string
	httpClient *http.Client
}

func NewService(ctx context.Context, log *logrus.Entry) *Service {
	return &Service{
		ctx:        ctx,
		log:        log,
		host:       env.ArgoCDURL,
		httpClient: common.NewHTTPClient(),
	}
}
