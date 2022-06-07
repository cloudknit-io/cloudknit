package health

import (
	"time"

	"github.com/compuzest/zlifecycle-event-service/app/event"
)

const (
	StatusOK            StatusType        = "ok"
	StatusUnknown       StatusType        = "unknown"
	StatusError         StatusType        = "error"
	HealthcheckOK       HealthcheckStatus = "ok"
	HealthcheckDegraded HealthcheckStatus = "degraded"
	HealthcheckError    HealthcheckStatus = "failure"
)

type (
	HealthcheckStatus string
	StatusType        string
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
	Component struct {
		Name     string            `json:"name"`
		Status   HealthcheckStatus `json:"status"`
		Critical bool              `json:"critical"`
	}
	Healthcheck struct {
		Status     HealthcheckStatus `json:"status"`
		Code       int               `json:"code"`
		Timestamp  time.Time         `json:"timestamp"`
		Components []*Component      `json:"components,omitempty"`
	}
)
