package health

import (
	"context"

	"github.com/pkg/errors"
)

func (s *Service) checkDB(ctx context.Context) error {
	if err := s.db.PingContext(ctx); err != nil {
		return errors.Wrap(err, "error executing database ping")
	}

	return nil
}
