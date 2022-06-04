package health

import "github.com/compuzest/zlifecycle-event-service/app/event"

const (
	StatusOK      StatusType = "ok"
	StatusUnknown StatusType = "unknown"
	StatusError   StatusType = "error"
)

type (
	TeamStatus        map[string]EnvironmentStatus
	TeamEvents        map[string]EnvironmentEvents
	EnvironmentStatus map[string][]*Status
	EnvironmentEvents map[string][]*event.Event
	Status            struct {
		Events      []*event.Event `json:"events"`
		Company     string         `json:"company"`
		Team        string         `json:"team"`
		Environment string         `json:"environment"`
		Status      StatusType     `json:"status"`
		Errors      []string       `json:"errors,omitempty"`
	}
)
