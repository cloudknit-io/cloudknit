package health

import (
	"context"
	"sort"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/compuzest/zlifecycle-event-service/app/util"

	"github.com/compuzest/zlifecycle-event-service/app/event"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type API interface {
	CompanyStatus(ctx context.Context, company string, eventEntries int, log *logrus.Entry) (TeamStatus, error)
	Healthcheck(ctx context.Context, fullCheck bool, log *logrus.Entry) *Healthcheck
}

type Service struct {
	es event.API
	db *sqlx.DB
}

func NewService(es event.API, db *sqlx.DB) *Service {
	return &Service{es: es, db: db}
}

func (s *Service) Healthcheck(ctx context.Context, fullCheck bool, log *logrus.Entry) *Healthcheck {
	log.Info("Performing app healthcheck")

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

func (s *Service) CompanyStatus(ctx context.Context, company string, eventEntries int, log *logrus.Entry) (TeamStatus, error) {
	log.Infof("Performing status check for company [%s]", company)

	payload := event.ListPayload{Company: company}
	events, err := s.es.List(ctx, &payload, event.ScopeCompany, event.FilterAll, log)
	if err != nil {
		return nil, errors.Wrapf(err, "error listing events for company [%s]", company)
	}

	teamEvents := buildTeamEvents(events)
	teamStatus, err := buildTeamStatus(teamEvents, eventEntries)
	if err != nil {
		return nil, errors.Wrapf(err, "error running healthcheck for company [%s]", company)
	}

	return teamStatus, nil
}

func buildTeamStatus(teamEvents TeamEvents, eventEntries int) (TeamStatus, error) {
	teamStatus := make(TeamStatus)

	for team, ee := range teamEvents {
		environmentStatus, err := buildEnvironmentStatus(ee, eventEntries)
		if err != nil {
			return nil, errors.Wrapf(err, "error creating team status for team [%s]", team)
		}
		teamStatus[team] = environmentStatus
	}

	return teamStatus, nil
}

func buildEnvironmentStatus(ee EnvironmentEvents, eventEntries int) (EnvironmentStatus, error) {
	environmentStatus := make(EnvironmentStatus)

	for env, events := range ee {
		status, err := NewEnvironmentStatus(events, eventEntries)
		if err != nil {
			return nil, errors.Wrapf(err, "error creating environment status for environment [%s]", env)
		}
		environmentStatus[env] = append(environmentStatus[env], status)
	}

	return environmentStatus, nil
}

func NewEnvironmentStatus(events []*event.Event, entries int) (*Status, error) {
	var errorMessages []string
	status := StatusOK
	latest := events[0]
	if latest.EventType == event.ValidationError {
		status = StatusError
		if err := util.CycleJSON(latest.Payload, &errorMessages); err != nil {
			return nil, errors.Wrapf(err, "error unmarshalling event [%s] payload", events[0].ID)
		}
	}
	return &Status{
		Events:      util.Truncate(events, entries),
		Company:     latest.Company,
		Team:        latest.Team,
		Environment: latest.Environment,
		Status:      status,
		Errors:      errorMessages,
	}, nil
}

func buildTeamEvents(events []*event.Event) TeamEvents {
	teamEvents := make(TeamEvents)

	for _, e := range events {
		if teamEvents[e.Team] == nil {
			teamEvents[e.Team] = buildEnvironmentEvents(events, e.Team)
		}
	}

	return teamEvents
}

func buildEnvironmentEvents(events []*event.Event, team string) EnvironmentEvents {
	environmentEvents := make(EnvironmentEvents)

	for _, e := range events {
		if !IsValidationEvent(e.EventType) || e.Team != team {
			continue
		}
		if environmentEvents[e.Environment] == nil {
			environmentEvents[e.Environment] = make([]*event.Event, 0)
		}
		if e.Team == team {
			environmentEvents[e.Environment] = append(environmentEvents[e.Environment], e)
		}
	}

	for _, evts := range environmentEvents {
		sort.Slice(evts, func(i, j int) bool {
			return evts[i].CreatedAt.After(evts[j].CreatedAt)
		})
	}

	return environmentEvents
}

func IsValidationEvent(eventType event.Type) bool {
	return util.StringInSlice(string(eventType), []string{string(event.ValidationError), string(event.ValidationSuccess)})
}
