package health

import (
	"context"
	"sort"

	"github.com/compuzest/zlifecycle-event-service/app/event"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type StatusType string

const (
	StatusOK      StatusType = "ok"
	StatusUnknown StatusType = "unknown"
	StatusError   StatusType = "error"
)

type API interface {
	Healthcheck(ctx context.Context, company string, log *logrus.Entry) (TeamStatus, error)
}

type Service struct {
	es event.API
}

func NewService(es event.API) *Service {
	return &Service{es: es}
}

type (
	TeamStatus        map[string]EnvironmentStatus
	TeamEvents        map[string]EnvironmentEvents
	EnvironmentStatus map[string][]*Status
	EnvironmentEvents map[string][]*event.Event
	Status            struct {
		Events []*event.Event `json:"events"`
		Status StatusType     `json:"status"`
	}
)

func (s *Service) Healthcheck(ctx context.Context, company string, log *logrus.Entry) (TeamStatus, error) {
	log.Infof("Performing healthcheck for company [%s]", company)

	payload := event.ListPayload{Company: company}
	events, err := s.es.List(ctx, &payload, event.ScopeCompany, event.FilterAll, log)
	if err != nil {
		return nil, errors.Wrapf(err, "error listing events for company [%s]", company)
	}

	teamEvents := buildTeamEvents(events)

	return buildTeamStatus(teamEvents), nil
}

func buildTeamStatus(teamEvents TeamEvents) TeamStatus {
	teamStatus := make(TeamStatus)

	for team, ee := range teamEvents {
		teamStatus[team] = buildEnvironmentStatus(ee)
	}

	return teamStatus
}

func buildEnvironmentStatus(ee EnvironmentEvents) EnvironmentStatus {
	environmentStatus := make(EnvironmentStatus)

	for env, events := range ee {
		environmentStatus[env] = append(environmentStatus[env], newEnvironmentStatus(events))
	}

	return environmentStatus
}

func newEnvironmentStatus(events []*event.Event) *Status {
	status := StatusOK
	if events[0].EventType == event.ValidationError {
		status = StatusError
	}
	return &Status{Status: status, Events: events}
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
		if e.Team != team {
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
