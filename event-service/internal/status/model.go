package status

import "github.com/compuzest/zlifecycle-event-service/internal/event"

const (
	StatusOK      Type = "ok"
	StatusUnknown Type = "unknown"
	StatusError   Type = "error"
)

type (
	Type     string
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
	EnvironmentStatus        struct {
		Object      string        `json:"object"`
		Company     string        `json:"company"`
		Team        string        `json:"team"`
		Environment string        `json:"environment"`
		Status      *ObjectStatus `json:"status"`
	}
	ObjectStatus struct {
		Events []*event.Event `json:"events"`
		Object string         `json:"object"`
		Meta   any            `json:"meta"`
		Status Type           `json:"status"`
		Errors []string       `json:"errors,omitempty"`
	}
)
