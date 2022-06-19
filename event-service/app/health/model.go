package health

import (
	"time"
)

const (
	HealthcheckOK       HealthcheckStatus = "ok"
	HealthcheckDegraded HealthcheckStatus = "degraded"
	HealthcheckError    HealthcheckStatus = "failure"
)

type (
	HealthcheckStatus string
	Component         struct {
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
