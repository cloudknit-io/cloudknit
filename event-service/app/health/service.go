package health

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/sirupsen/logrus"
)

type API interface {
	Healthcheck(ctx context.Context, fullCheck bool, log *logrus.Entry) *Healthcheck
}

type Service struct {
	db *sqlx.DB
}

func NewService(db *sqlx.DB) *Service {
	return &Service{db: db}
}

func (s *Service) Healthcheck(ctx context.Context, fullCheck bool, log *logrus.Entry) *Healthcheck {
	hc := Healthcheck{
		Timestamp: time.Now(),
	}

	if fullCheck {
		hc.Components = s.fullCheck(ctx, log)
	}

	checkComponentHealth(&hc, log)

	return &hc
}

func (s *Service) fullCheck(ctx context.Context, log *logrus.Entry) []*Component {
	var components []*Component

	dbComponent := Component{
		Name:     "db",
		Critical: true,
	}
	if err := s.checkDB(ctx); err != nil {
		log.Errorf("error performing db healthcheck: %v", err)
		dbComponent.Status = HealthcheckError
	} else {
		dbComponent.Status = HealthcheckOK
	}
	components = append(components, &dbComponent)

	return components
}

func checkComponentHealth(hc *Healthcheck, log *logrus.Entry) {
	hc.Status = HealthcheckOK
	hc.Code = 200
	for _, c := range hc.Components {
		if c.Status == HealthcheckDegraded && hc.Status != HealthcheckError {
			log.Warnf("Component [%s] status is in degraded state and critial status is set to %t", c.Name, c.Critical)
			log.Warnf("Marking healthcheck status as degraded")
			hc.Status = HealthcheckDegraded
		}
		if c.Status == HealthcheckError {
			log.Warnf("Component [%s] status is in error state and critial status is set to %t", c.Name, c.Critical)
			if c.Critical {
				log.Warnf("Marking healthcheck status as failed")
				hc.Status = HealthcheckError
				hc.Code = 500
			} else {
				log.Errorf("Marking healthcheck status as degraded")
				hc.Status = HealthcheckDegraded
				hc.Code = 200
			}

			break
		}
	}
}
