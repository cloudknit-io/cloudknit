package event

import "github.com/jmoiron/sqlx"

func (s *Service) sqlInsertEvent() (*sqlx.Stmt, error) {
	return s.db.Preparex(
		"INSERT INTO event(id, company, team, environment, created_at, event_type, message, debug) VALUES(:id, :company, :team, :environment, :created_at, :event_type, :message, :debug)",
	)
}

func (s *Service) sqlGetEventsForEnvironment() (*sqlx.NamedStmt, error) {
	return s.db.PrepareNamed("SELECT * from event WHERE event.company = :company AND event.team = :team AND event.environment = :environment")
}

func (s *Service) sqlGetLatestEventForEnvironment() (*sqlx.NamedStmt, error) {
	return s.db.PrepareNamed(
		"SELECT * from event WHERE event.company = :company AND event.team = :team AND event.environment = :environment ORDER BY event.created_at DESC LIMIT 1",
	)
}

func (s *Service) sqlGetEventsForTeam() (*sqlx.NamedStmt, error) {
	return s.db.PrepareNamed("SELECT * from event WHERE event.company = :company AND event.team = :team")
}

func (s *Service) sqlGetLatestEventForTeam() (*sqlx.NamedStmt, error) {
	return s.db.PrepareNamed("SELECT * from event WHERE event.company = :company AND event.team = :team ORDER BY event.created_at DESC LIMIT 1")
}

func (s *Service) sqlGetEventsForCompany() (*sqlx.NamedStmt, error) {
	return s.db.PrepareNamed("SELECT * from event WHERE event.company = :company")
}

func (s *Service) sqlGetLatestEventForCompany() (*sqlx.NamedStmt, error) {
	return s.db.PrepareNamed("SELECT * from event WHERE event.company = :company ORDER BY event.created_at DESC LIMIT 1")
}
