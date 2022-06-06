package health

import "github.com/jmoiron/sqlx"

func (s *Service) sqlDBHealthcheck() (*sqlx.Stmt, error) {
	return s.db.Preparex("SELECT 1")
}
