package services

import (
	"context"
	"github.com/sirupsen/logrus"

	"github.com/compuzest/zlifecycle-event-service/internal/status"

	"github.com/compuzest/zlifecycle-event-service/internal/stream"

	"github.com/compuzest/zlifecycle-event-service/internal/db"
	"github.com/compuzest/zlifecycle-event-service/internal/event"
	"github.com/compuzest/zlifecycle-event-service/internal/health"
	"github.com/pkg/errors"
)

type Services struct {
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
	sseBroker := stream.NewService(l)

	return &Services{ES: es, HS: hs, SS: ss, SSEBroker: sseBroker}, nil
}
