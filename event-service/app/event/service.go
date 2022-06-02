package event

import (
	"context"

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
	Record(ctx context.Context, e *Event) error
	List(ctx context.Context, payload *ListPayload, scope Scope, filter Filter) (events []*Event, err error)
}

type Service struct {
	db *sqlx.DB
}

func NewService(db *sqlx.DB) *Service {
	return &Service{db: db}
}

func (s *Service) Record(ctx context.Context, e *Event) error {
	stmt, err := s.sqlInsertEvent()
	if err != nil {
		return errors.Wrap(err, "error preparing statement: insert-event")
	}

	result, err := stmt.ExecContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "error executing prepared statement: insert-event")
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "error getting affected rows for executed statement: insert-event ")
	}

	if rows != 1 {
		return errors.Errorf("invalid affected rows count, must be 1: %d", rows)
	}

	return nil
}

func (s *Service) List(ctx context.Context, payload *ListPayload, scope Scope, filter Filter) (events []*Event, err error) {
	if err := validatePayload(payload, scope); err != nil {
		return nil, errors.Wrap(err, "error validating list events payload")
	}

	stmt, err := s.getListStmt(scope, filter)
	if err != nil {
		return nil, err
	}

	if err = stmt.SelectContext(ctx, &events, payload); err != nil {
		return nil, errors.Wrapf(err, "error listing events for %s scope and %s filter", scope, filter)
	}
	return
}

func validatePayload(p *ListPayload, scope Scope) error {
	switch scope {
	case ScopeCompany:
		if p.Company == "" {
			return errors.New("company must be defined when scope is set to company")
		}
	case ScopeTeam:
		if p.Company == "" {
			return errors.New("company must be defined when scope is set to team")
		}
		if p.Team == "" {
			return errors.New("company must be defined when scope is set to team")
		}
	case ScopeEnvironment:
		if p.Company == "" {
			return errors.New("company must be defined when scope is set to environment")
		}
		if p.Team == "" {
			return errors.New("company must be defined when scope is set to environment")
		}
		if p.Environment == "" {
			return errors.New("company must be defined when scope is set to environment")
		}
	default:
		return errors.Errorf("invalid scope: %s", scope)
	}
	return nil
}

func (s *Service) getListStmt(scope Scope, filter Filter) (stmt *sqlx.NamedStmt, err error) {
	switch scope {
	case ScopeCompany:
		stmt, err = s.sqlGetEventsForCompany()
		if err != nil {
			return nil, errors.Wrap(err, "error preparing get-events-for-company statement")
		}
		if filter == FilterLatest {
			stmt, err = s.sqlGetLatestEventForCompany()
			if err != nil {
				return nil, errors.Wrap(err, "error preparing get-latest-events-for-company statement")
			}
		}
	case ScopeTeam:
		stmt, err = s.sqlGetEventsForTeam()
		if err != nil {
			return nil, errors.Wrap(err, "error preparing get-events-for-team statement")
		}
		if filter == FilterLatest {
			stmt, err = s.sqlGetLatestEventForTeam()
			if err != nil {
				return nil, errors.Wrap(err, "error preparing get-latest-events-for-team statement")
			}
		}
	case ScopeEnvironment:
		stmt, err = s.sqlGetEventsForEnvironment()
		if err != nil {
			return nil, errors.Wrap(err, "error preparing get-events-for-environment statement")
		}
		if filter == FilterLatest {
			stmt, err = s.sqlGetLatestEventForEnvironment()
			if err != nil {
				return nil, errors.Wrap(err, "error preparing get-latest-events-for-environment statement")
			}
		}
	default:
		return nil, errors.Errorf("invalid scope: %s", scope)
	}
	return stmt, err
}
