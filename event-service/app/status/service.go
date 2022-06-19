package status

import (
	"context"
	"github.com/compuzest/zlifecycle-event-service/app/event"
	"github.com/compuzest/zlifecycle-event-service/app/util"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"sort"
)

type API interface {
	CompanyStatus(ctx context.Context, company string, eventEntries int, log *logrus.Entry) (TeamStatus, error)
}

type Service struct {
	es event.API
	db *sqlx.DB
}

func NewService(es event.API, db *sqlx.DB) *Service {
	return &Service{es: es, db: db}
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
	isErrorEvent := latest.EventType == event.ValidationError || latest.EventType == event.SchemaValidationError
	if isErrorEvent {
		status = StatusError
		if err := util.CycleJSON(latest.Payload, &errorMessages); err != nil {
			return nil, errors.Wrapf(err, "error unmarshalling event [%s] payload", events[0].ID)
		}
	}
	return &Status{
		Events:      util.Truncate(events, entries),
		Object:      latest.Object,
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
		if !event.IsValidationEvent(e.EventType) || e.Team != team {
			continue
		}
		if e.Team == team {
			if environmentEvents[e.Object] == nil {
				environmentEvents[e.Object] = make([]*event.Event, 0)
			}
			environmentEvents[e.Object] = append(environmentEvents[e.Object], e)
		}
	}

	for _, evts := range environmentEvents {
		sortEvents(evts)
	}

	return environmentEvents
}

func sortEvents(events []*event.Event) {
	sort.Slice(events, func(i, j int) bool {
		return events[i].CreatedAt.After(events[j].CreatedAt)
	})
}
