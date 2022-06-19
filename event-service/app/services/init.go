package services

import (
	"context"
	"github.com/compuzest/zlifecycle-event-service/app/status"

	"github.com/compuzest/zlifecycle-event-service/app/stream"

	"github.com/compuzest/zlifecycle-event-service/app/db"
	"github.com/compuzest/zlifecycle-event-service/app/event"
	"github.com/compuzest/zlifecycle-event-service/app/health"
	"github.com/pkg/errors"
)

type Services struct {
	ES        event.API
	HS        health.API
	SS        status.API
	SSEBroker stream.API
}

func NewServices() (*Services, error) {
	ctx := context.Background()
	sqldb, err := db.NewDatabase(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "error creating new database connection")
	}
	es := event.NewService(sqldb)
	hs := health.NewService(sqldb)
	ss := status.NewService(es, sqldb)
	sseBroker := stream.NewService()

	return &Services{ES: es, HS: hs, SS: ss, SSEBroker: sseBroker}, nil
}
