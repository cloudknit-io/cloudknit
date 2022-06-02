package services

import (
	"context"

	"github.com/compuzest/zlifecycle-event-service/app/db"
	"github.com/compuzest/zlifecycle-event-service/app/event"
	"github.com/pkg/errors"
)

type Services struct {
	ES event.API
}

func NewServices() (*Services, error) {
	ctx := context.Background()
	db, err := db.NewDatabase(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "error creating new database")
	}
	es := event.NewService(db)

	return &Services{ES: es}, nil
}
