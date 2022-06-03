package event

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type (
	Scope  string
	Filter string
)

const (
	ScopeCompany     Scope  = "company"
	ScopeTeam        Scope  = "team"
	ScopeEnvironment Scope  = "environment"
	FilterAll        Filter = "all"
	FilterLatest     Filter = "latest"
)

type API interface {
	Record(ctx context.Context, e *RecordPayload, log *logrus.Entry) (*Event, error)
	List(ctx context.Context, payload *ListPayload, scope Scope, filter Filter, log *logrus.Entry) (events []*Event, err error)
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

	event := NewEvent(p.Company, p.Team, p.Environment, p.Message, Type(p.EventType), p.Debug)

	log.Infof(
		"Recording new %s event with ID [%s] for company [%s], team [%s] and environment [%s]",
		event.EventType, event.ID, event.Company, event.Team, event.Environment,
	)

	if err := s.insertEvent(ctx, event); err != nil {
		return nil, errors.Wrapf(
			err,
			"error persisting event [%s] for company [%s], team [%s] and environment [%s]",
			event.ID, event.Company, event.Team, event.Environment,
		)
	}

	return event, nil
}

func (s *Service) List(ctx context.Context, payload *ListPayload, scope Scope, filter Filter, log *logrus.Entry) (events []*Event, err error) {
	if err = validateListPayload(payload, scope); err != nil {
		return nil, errors.Wrap(err, "error validating list events payload")
	}

	if filter == "" {
		filter = FilterAll
	}

	log.WithFields(logrus.Fields{
		"company": payload.Company, "team": payload.Team, "environment": payload.Environment,
	}).Infof("Listing events with scope [%s] and filter [%s]", scope, filter)

	events, err = s.selectEvents(ctx, payload, scope, filter)
	if err != nil {
		return nil, errors.Wrapf(err, "error selecting events for scope %s and filter %s", scope, filter)
	}

	return events, err
}
