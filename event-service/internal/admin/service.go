package admin

import (
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type API interface {
	WipeDatabase() error
}

type Service struct {
	db *sqlx.DB
}

func NewService(db *sqlx.DB) *Service {
	return &Service{
		db: db,
	}
}

func (s *Service) WipeDatabase() error {
	if _, err := s.db.Exec(sqlDropTableEvent()); err != nil {
		return errors.Wrap(err, "error dropping event table")
	}
	if _, err := s.db.Exec(sqlTruncateSchemaMigrations()); err != nil {
		return errors.Wrap(err, "error truncating schema_migrations table")
	}
	return nil
}
