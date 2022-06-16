package event

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

func (s *Service) insertEvent(ctx context.Context, event *Event) error {
	stmt, err := s.sqlInsertEvent()
	if err != nil {
		return errors.Wrap(err, "error preparing statement: insert-event")
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, event)
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

func (s *Service) selectEvents(ctx context.Context, payload *ListPayload, scope Scope, filter Filter) (events []*Event, err error) {
	stmt, err := s.getListStmt(scope, filter)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	if err = stmt.SelectContext(ctx, &events, payload); err != nil {
		return nil, errors.Wrapf(err, "error listing events for %s scope and %s filter", scope, filter)
	}

	if events == nil {
		events = []*Event{}
	}

	return events, nil
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
		stmt, err = s.sqlGetEventsForEnvironmentObject()
		if err != nil {
			return nil, errors.Wrap(err, "error preparing get-events-for-environment statement")
		}
		if filter == FilterLatest {
			stmt, err = s.sqlGetLatestEventForEnvironmentObject()
			if err != nil {
				return nil, errors.Wrap(err, "error preparing get-latest-events-for-environment statement")
			}
		}
	default:
		return nil, errors.Errorf("invalid scope: %s", scope)
	}
	return stmt, err
}

// QUERIES

func (s *Service) sqlInsertEvent() (*sqlx.NamedStmt, error) {
	return s.db.PrepareNamed(
		"INSERT INTO event(id, object, company, team, environment, created_at, event_type, payload, debug) VALUES(:id, :object, :company, :team, :environment, :created_at, :event_type, :payload, :debug)",
	)
}

func (s *Service) sqlGetEventsForEnvironmentObject() (*sqlx.NamedStmt, error) {
	return s.db.PrepareNamed("SELECT * FROM event WHERE event.object = :object ORDER BY event.created_at DESC")
}

func (s *Service) sqlGetLatestEventForEnvironmentObject() (*sqlx.NamedStmt, error) {
	return s.db.PrepareNamed("SELECT * FROM event WHERE event.object = :object ORDER BY event.created_at DESC LIMIT 1")
}

func (s *Service) sqlGetEventsForEnvironment() (*sqlx.NamedStmt, error) {
	return s.db.PrepareNamed("SELECT * from event WHERE event.company = :company AND event.team = :team AND event.environment = :environment ORDER BY event.company, event.team, event.environment, event.created_at DESC")
}

func (s *Service) sqlGetLatestEventForEnvironment() (*sqlx.NamedStmt, error) {
	return s.db.PrepareNamed(
		"SELECT * from event WHERE event.company = :company AND event.team = :team AND event.environment = :environment ORDER BY event.company, event.team, event.environment, event.created_at DESC LIMIT 1",
	)
}

func (s *Service) sqlGetEventsForTeam() (*sqlx.NamedStmt, error) {
	return s.db.PrepareNamed("SELECT * FROM event WHERE event.company = :company AND event.team = :team ORDER BY event.company, event.team, event.created_at DESC")
}

func (s *Service) sqlGetLatestEventForTeam() (*sqlx.NamedStmt, error) {
	return s.db.PrepareNamed("SELECT * FROM event WHERE event.company = :company AND event.team = :team ORDER BY event.company, event.team, event.created_at DESC LIMIT 1")
}

func (s *Service) sqlGetEventsForCompany() (*sqlx.NamedStmt, error) {
	return s.db.PrepareNamed("SELECT * FROM event WHERE event.company = :company ORDER BY event.company, event.created_at DESC")
}

func (s *Service) sqlGetLatestEventForCompany() (*sqlx.NamedStmt, error) {
	return s.db.PrepareNamed("SELECT * FROM event WHERE event.company = :company ORDER BY event.company, event.created_at DESC LIMIT 1")
}
