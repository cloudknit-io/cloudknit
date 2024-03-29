package services

import (
	"context"

	"github.com/cloudknit-io/cloudknit/event-service/internal/admin"
	"github.com/sirupsen/logrus"

	"github.com/cloudknit-io/cloudknit/event-service/internal/status"

	"github.com/cloudknit-io/cloudknit/event-service/internal/stream"

	"github.com/cloudknit-io/cloudknit/event-service/internal/db"
	"github.com/cloudknit-io/cloudknit/event-service/internal/event"
	"github.com/cloudknit-io/cloudknit/event-service/internal/health"
	"github.com/pkg/errors"
)

type Services struct {
	AS        admin.API
	ES        event.API
	HS        health.API
	SS        status.API
	SSEBroker stream.API
}

func NewServices(l *logrus.Entry) (*Services, error) {
	ctx := context.Background()
	sqldb, err := db.NewDatabase(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "error creating new database connection")
	}
	es := event.NewService(sqldb)
	hs := health.NewService(sqldb)
	ss := status.NewService(es, sqldb)
	as := admin.NewService(sqldb)
	sseBroker := stream.NewService(l)

	return &Services{AS: as, ES: es, HS: hs, SS: ss, SSEBroker: sseBroker}, nil
}
