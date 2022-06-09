package health

import (
	"context"
	"github.com/pkg/errors"
)

func (s *Service) checkDB(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, sqlDBHealthcheck())
	if err != nil {
		return errors.Wrap(err, "error executing healthcheck sql query")
	}

	return nil
}

func sqlDBHealthcheck() string {
	return "SELECT 1"
}
