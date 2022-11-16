package cloudknitservice

import (
	"context"

	"github.com/sirupsen/logrus"
)

type API interface {
	Get(ctx context.Context, organizationName string, log *logrus.Entry) (*GetOrganizationResponse, error)
}
