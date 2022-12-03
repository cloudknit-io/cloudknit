package event

import (
	"context"

	"github.com/cloudknit-io/cloudknit/event-service/internal/util"

	"github.com/sirupsen/logrus"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type (
	Scope  string
	Filter string
)

type API interface {
	Record(ctx context.Context, e *RecordPayload, log *logrus.Entry) (*Event, error)
	List(ctx context.Context, payload *ListPayload, log *logrus.Entry) (events []*Event, err error)
}

type Service struct {
	db *sqlx.DB
}

func NewService(db *sqlx.DB) *Service {
	return &Service{db: db}
}

func (s *Service) Record(ctx context.Context, p *RecordPayload, log *logrus.Entry) (*Event, error) {
	if err := validateRecordPayload(p); err != nil {
		return nil, errors.Wrap(err, "error validating record payload")
	}

	event, err := NewEvent(p.Scope, p.Object, p.Meta, p.Payload, Type(p.EventType), p.Debug)
	if err != nil {
		return nil, errors.Wrapf(err, "error creating new event for object [%s] with scope [%s]", p.Object, p.Scope)
	}

	log.WithField("meta", p.Meta).Infof(
		"Recording new %s event with ID [%s] for object [%s] and scope [%s]",
		event.EventType, event.ID, event.Object, event.Scope,
	)

	if err := s.insertEvent(ctx, event); err != nil {
		return nil, errors.Wrapf(err, "error persisting event [%s] for object [%s]", event.ID, event.Object)
	}

	return event, nil
}

func (s *Service) List(ctx context.Context, p *ListPayload, log *logrus.Entry) (events []*Event, err error) {
	if err = validateListPayload(p); err != nil {
		return nil, errors.Wrap(err, "error validating list events payload")
	}

	if p.Filter == "" {
		p.Filter = FilterAll
	}

	log.WithField("payload", p).Infof("Listing events with scope [%s] and filter [%s]", p.Scope, p.Filter)

	events, err = s.selectEvents(ctx, p)
	if err != nil {
		return nil, errors.Wrapf(err, "error selecting events for scope %s and filter %s", p.Scope, p.Filter)
	}

	return events, err
}

func GetFamilyForType(eventType Type) (Family, error) {
	switch {
	case isValidationEvent(eventType):
		return FamilyValidation, nil
	case isReconcileEvent(eventType):
		return FamilyReconcile, nil
	default:
		return "", errors.Errorf("invalid event type: %s", eventType)
	}
}

func isValidationEvent(eventType Type) bool {
	return util.IsInSlice(
		eventType,
		[]Type{
			EnvironmentSpecValidationSuccess,
			EnvironmentSchemaValidationSuccess,
			EnvironmentSpecValidationError,
			EnvironmentSchemaValidationError,
			TeamSpecValidationSuccess,
			TeamSchemaValidationSuccess,
			TeamSpecValidationError,
			TeamSchemaValidationError,
		},
	)
}

func isReconcileEvent(eventType Type) bool {
	return util.IsInSlice(
		eventType,
		[]Type{
			EnvironmentReconcileSuccess,
			EnvironmentReconcileError,
			TeamReconcileSuccess,
			TeamReconcileError,
		},
	)
}

func IsErrorEvent(eventType Type) bool {
	return util.IsInSlice(
		eventType,
		[]Type{
			EnvironmentSpecValidationError,
			EnvironmentSchemaValidationError,
			EnvironmentReconcileError,
			TeamSpecValidationError,
			TeamSchemaValidationError,
			TeamReconcileError,
		},
	)
}

func isSupportedEvent(eventType Type) bool {
	return util.IsInSlice(
		eventType,
		[]Type{
			EnvironmentSpecValidationError,
			EnvironmentSpecValidationSuccess,
			EnvironmentSchemaValidationError,
			EnvironmentSchemaValidationSuccess,
			EnvironmentReconcileSuccess,
			EnvironmentReconcileError,
			TeamSpecValidationError,
			TeamSpecValidationSuccess,
			TeamSchemaValidationError,
			TeamSchemaValidationSuccess,
			TeamReconcileSuccess,
			TeamReconcileError,
		},
	)
}
