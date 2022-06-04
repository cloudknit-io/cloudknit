package health

import (
	"context"
	"sort"

	"github.com/compuzest/zlifecycle-event-service/app/util"

	"github.com/compuzest/zlifecycle-event-service/app/event"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type StatusType string

type API interface {
	Healthcheck(ctx context.Context, company string, log *logrus.Entry) (TeamStatus, error)
}

type Service struct {
	es event.API
}

func NewService(es event.API) *Service {
	return &Service{es: es}
}

func (s *Service) Healthcheck(ctx context.Context, company string, log *logrus.Entry) (TeamStatus, error) {
	log.Infof("Performing healthcheck for company [%s]", company)

	payload := event.ListPayload{Company: company}
	events, err := s.es.List(ctx, &payload, event.ScopeCompany, event.FilterAll, log)
	if err != nil {
		return nil, errors.Wrapf(err, "error listing events for company [%s]", company)
	}

	teamEvents := buildTeamEvents(events)
	teamStatus, err := buildTeamStatus(teamEvents)
	if err != nil {
		return nil, errors.Wrapf(err, "error running healthcheck for company [%s]", company)
	}

	return teamStatus, nil
}

func buildTeamStatus(teamEvents TeamEvents) (TeamStatus, error) {
	teamStatus := make(TeamStatus)

	for team, ee := range teamEvents {
		environmentStatus, err := buildEnvironmentStatus(ee)
		if err != nil {
			return nil, errors.Wrapf(err, "error creating team status for team [%s]", team)
		}
		teamStatus[team] = environmentStatus
	}

	return teamStatus, nil
}

func buildEnvironmentStatus(ee EnvironmentEvents) (EnvironmentStatus, error) {
	environmentStatus := make(EnvironmentStatus)

	for env, events := range ee {
		status, err := NewEnvironmentStatus(events)
		if err != nil {
			return nil, errors.Wrapf(err, "error creating environment status for environment [%s]", env)
		}
		environmentStatus[env] = append(environmentStatus[env], status)
	}

	return environmentStatus, nil
}

func NewEnvironmentStatus(events []*event.Event) (*Status, error) {
	var errorMessages []string
	status := StatusOK
	latest := events[0]
	if latest.EventType == event.ValidationError {
		status = StatusError
		if err := util.CycleJSON(latest.Payload, &errorMessages); err != nil {
			return nil, errors.Wrapf(err, "error unmarshalling event [%s] payload", events[0].ID)
		}
	}
	limit := 5
	if len(events) < limit {
		limit = len(events)
	}
	return &Status{
		Events:      events[0:limit],
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
