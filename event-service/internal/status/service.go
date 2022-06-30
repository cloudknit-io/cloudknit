package status

import (
	"context"
	"github.com/tidwall/gjson"
	"sort"

	fpgo "github.com/TeaEntityLab/fpGo/v2"
	"github.com/compuzest/zlifecycle-event-service/internal/event"
	"github.com/compuzest/zlifecycle-event-service/internal/util"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type API interface {
	Calculate(ctx context.Context, company string, eventEntries int, log *logrus.Entry) (*Response, error)
}

type Service struct {
	es event.API
	db *sqlx.DB
}

func NewService(es event.API, db *sqlx.DB) *Service {
	return &Service{es: es, db: db}
}

func (s *Service) Calculate(ctx context.Context, company string, eventEntries int, log *logrus.Entry) (*Response, error) {
	log.Infof("Performing status check for company [%s]", company)

	lp := event.ListPayload{Company: company, Scope: event.ScopeCompany, Filter: event.FilterAll}
	events, err := s.es.List(ctx, &lp, log)
	if err != nil {
		return nil, errors.Wrapf(err, "error listing events for company [%s]", company)
	}

	groupedTeamEvents, err := groupTeamEvents(events)
	if err != nil {
		return nil, errors.Wrap(err, "error grouping team events")
	}
	environmentEvents, err := groupEnvironmentEvents(events)
	if err != nil {
		return nil, errors.Wrap(err, "error grouping environment events")
	}
	teamStatusMap, err := buildTeamStatusMap(groupedTeamEvents, eventEntries)
	if err != nil {
		return nil, err
	}
	environmentStatusMap, err := buildEnvironmentStatusMap(environmentEvents, eventEntries)
	if err != nil {
		return nil, err
	}

	return &Response{
		TeamsStatus:       teamStatusMap,
		EnvironmentStatus: environmentStatusMap,
	}, nil
}

func groupTeamEvents(allEvents []*event.Event) (GroupedTeamEvents, error) {
	teamEvents := extractEventsForScope(allEvents, event.ScopeTeam)
	grouped := make(GroupedTeamEvents)
	for _, evt := range teamEvents {
		grouped[evt.Object] = append(grouped[evt.Object], evt)
	}
	for _, teamEvents := range grouped {
		sortEventsByDescTimestamp(teamEvents)
	}
	return grouped, nil
}

func buildTeamStatusMap(grouped GroupedTeamEvents, eventEntries int) (TeamStatusMap, error) {
	m := make(TeamStatusMap, len(grouped))

	for team, events := range grouped {
		teamStatus, err := NewTeamStatus(events, eventEntries)
		if err != nil {
			return nil, errors.Wrap(err, "error creating team status")
		}
		m[team] = teamStatus
	}

	return m, nil
}

func NewTeamStatus(events []*event.Event, eventEntries int) (*TeamStatus, error) {
	status, err := NewObjectStatus(events, eventEntries)
	if err != nil {
		return nil, errors.Wrap(err, "error creating team object status")
	}
	latest := events[0]

	company := gjson.GetBytes(latest.Meta, "company").String()
	if company == "" {
		return nil, errors.New("event meta is missing team field")
	}
	team := gjson.GetBytes(latest.Meta, "team").String()
	if team == "" {
		return nil, errors.New("event meta is missing team field")
	}
	return &TeamStatus{
		Object:  latest.Object,
		Company: company,
		Team:    team,
		Status:  status,
	}, nil
}

func groupEnvironmentEvents(allEvents []*event.Event) (GroupedEnvironmentEvents, error) {
	environmentEvents := extractEventsForScope(allEvents, event.ScopeEnvironment)
	grouped := make(GroupedEnvironmentEvents)

	for _, evt := range environmentEvents {
		team := gjson.GetBytes(evt.Meta, "team").String()
		if team == "" {
			return nil, errors.New("event meta is missing team field")
		}
		if grouped[team] == nil {
			grouped[team] = make(map[string][]*event.Event)
		}
		grouped[team][evt.Object] = append(grouped[team][evt.Object], evt)
	}
	for _, groupedEnvironmentMap := range grouped {
		for _, events := range groupedEnvironmentMap {
			sortEventsByDescTimestamp(events)
		}
	}
	return grouped, nil
}

func buildEnvironmentStatusMap(grouped GroupedEnvironmentEvents, eventEntries int) (EnvironmentStatusMap, error) {
	m := make(EnvironmentStatusMap, len(grouped))

	for team, environmentEventsMap := range grouped {
		if m[team] == nil {
			m[team] = make(map[string]*EnvironmentStatus)
		}
		for environment, events := range environmentEventsMap {
			environmentStatus, err := NewEnvironmentStatus(events, eventEntries)
			if err != nil {
				return nil, errors.Wrap(err, "error creating environment status")
			}
			m[team][environment] = environmentStatus
		}
	}

	return m, nil
}

func NewEnvironmentStatus(events []*event.Event, eventEntries int) (*EnvironmentStatus, error) {
	status, err := NewObjectStatus(events, eventEntries)
	if err != nil {
		return nil, errors.Wrap(err, "error creating environment object status")
	}
	latest := events[0]

	company := gjson.GetBytes(latest.Meta, "company").String()
	if company == "" {
		return nil, errors.New("event meta is missing company field")
	}
	team := gjson.GetBytes(latest.Meta, "team").String()
	if team == "" {
		return nil, errors.New("event meta is missing team field")
	}
	environment := gjson.GetBytes(latest.Meta, "environment").String()
	if environment == "" {
		return nil, errors.New("event meta is missing environment field")
	}

	return &EnvironmentStatus{
		Object:      latest.Object,
		Company:     company,
		Team:        team,
		Environment: environment,
		Status:      status,
	}, nil
}

func extractEventsForScope(events []*event.Event, scope event.Scope) []*event.Event {
	return fpgo.StreamFromArray(events).Filter(func(e *event.Event, i int) bool {
		return e.Scope == scope
	}).ToArray()
}

func NewObjectStatus(events []*event.Event, entries int) (*ObjectStatus, error) {
	status := StatusOK
	latest := events[0]
	isErrorEvent := latest.EventFamily == event.FamilyValidationError
	var errorMessages []string
	if isErrorEvent {
		status = StatusError
		if err := util.CycleJSON(latest.Payload, &errorMessages); err != nil {
			return nil, errors.Wrapf(err, "error unmarshalling event [%s] payload", events[0].ID)
		}
	}
	return &ObjectStatus{
		Events: util.Truncate(events, entries),
		Object: latest.Object,
		Meta:   latest.Meta,
		Status: status,
		Errors: errorMessages,
	}, nil
}

func sortEventsByDescTimestamp(events []*event.Event) {
	sort.Slice(events, func(i, j int) bool {
		return events[i].CreatedAt.After(events[j].CreatedAt)
	})
}
