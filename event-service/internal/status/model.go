package status

import (
	"time"

	"github.com/cloudknit-io/cloudknit/event-service/internal/event"
)

const (
	StateOK      State = "ok"
	StateUnknown State = "unknown"
	StateError   State = "error"
)

type (
	State    string
	Response struct {
		TeamsStatus       TeamStatusMap        `json:"teamsStatus"`
		EnvironmentStatus EnvironmentStatusMap `json:"environmentStatus"`
	}
	TeamStatusMap     map[string]*TeamStatus
	GroupedTeamEvents map[string][]*event.Event
	TeamStatus        struct {
		Object  string        `json:"object"`
		Company string        `json:"company"`
		Team    string        `json:"team"`
		Status  *ObjectStatus `json:"status"`
	}
	EnvironmentStatusMap     map[string]map[string]*EnvironmentStatus
	GroupedEnvironmentEvents map[string]map[string][]*event.Event
	GroupedEventsByFamily    map[event.Family][]*event.Event
	EnvironmentStatus        struct {
		Object      string        `json:"object"`
		Company     string        `json:"company"`
		Team        string        `json:"team"`
		Environment string        `json:"environment"`
		Status      *ObjectStatus `json:"status"`
	}
	ObjectStatus struct {
		Events []*event.Event           `json:"events"`
		Object string                   `json:"object"`
		Meta   any                      `json:"meta"`
		Status map[event.Family]*Status `json:"status"`
	}
)

type Status struct {
	State     State        `json:"state"`
	Family    event.Family `json:"family"`
	Errors    []string     `json:"errors,omitempty"`
	Timestamp time.Time    `json:"timestamp"`
}
