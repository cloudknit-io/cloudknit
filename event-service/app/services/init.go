package services

import (
	"context"
	"github.com/compuzest/zlifecycle-event-service/app/sse"

	"github.com/compuzest/zlifecycle-event-service/app/db"
	"github.com/compuzest/zlifecycle-event-service/app/event"
	"github.com/compuzest/zlifecycle-event-service/app/health"
	"github.com/pkg/errors"
)

type Services struct {
	ES        event.API
	SS        health.API
	SSEBroker sse.API
}

func NewServices() (*Services, error) {
	ctx := context.Background()
	sqldb, err := db.NewDatabase(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "error creating new database connection")
	}
	es := event.NewService(sqldb)
	hs := health.NewService(es, sqldb)
	sseBroker := sse.NewSSE()

	return &Services{ES: es, SS: hs, SSEBroker: sseBroker}, nil
}
