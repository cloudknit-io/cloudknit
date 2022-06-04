package services

import (
	"context"

	"github.com/compuzest/zlifecycle-event-service/app/health"
	"github.com/compuzest/zlifecycle-event-service/app/web/sse"

	"github.com/compuzest/zlifecycle-event-service/app/db"
	"github.com/compuzest/zlifecycle-event-service/app/event"
	"github.com/pkg/errors"
)

type Services struct {
	ES        event.API
	HS        health.API
	SSEBroker sse.API
}

func NewServices() (*Services, error) {
	ctx := context.Background()
	sqldb, err := db.NewDatabase(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "error creating new database")
	}
	es := event.NewService(sqldb)
	hs := health.NewService(es)
	sseBroker := sse.NewSSE()

	return &Services{ES: es, HS: hs, SSEBroker: sseBroker}, nil
}
