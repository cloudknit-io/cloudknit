package event

import (
	"encoding/json"
	"time"

	"github.com/pkg/errors"

	"github.com/cloudknit-io/cloudknit/event-service/internal/util"

	"github.com/google/uuid"
)

type (
	Type   string
	Family string
)

const (
	ScopeCompany                       Scope  = "company"
	ScopeTeam                          Scope  = "team"
	ScopeEnvironment                   Scope  = "environment"
	FilterAll                          Filter = "all"
	FilterLatest                       Filter = "latest"
	EnvironmentSpecValidationSuccess   Type   = "environment_validation_success"
	EnvironmentSpecValidationError     Type   = "environment_validation_error"
	EnvironmentSchemaValidationError   Type   = "environment_schema_validation_error"
	EnvironmentSchemaValidationSuccess Type   = "environment_schema_validation_success"
	EnvironmentReconcileSuccess        Type   = "environment_reconcile_success"
	EnvironmentReconcileError          Type   = "environment_reconcile_error"
	TeamSpecValidationSuccess          Type   = "team_validation_success"
	TeamSpecValidationError            Type   = "team_validation_error"
	TeamSchemaValidationError          Type   = "team_schema_validation_error"
	TeamSchemaValidationSuccess        Type   = "team_schema_validation_success"
	TeamReconcileSuccess               Type   = "team_reconcile_success"
	TeamReconcileError                 Type   = "team_reconcile_error"
	FamilyValidation                   Family = "validation"
	FamilyReconcile                    Family = "reconcile"
)

type Event struct {
	ID        string          `json:"id" db:"id"`
	Scope     Scope           `json:"scope" db:"scope"`
	Object    string          `json:"object" db:"object"`
	Meta      json.RawMessage `json:"meta" db:"meta"`
	CreatedAt time.Time       `json:"createdAt" db:"created_at"`
	EventType Type            `json:"eventType" db:"event_type"`
	Family    Family          `json:"family" db:"family"`
	Payload   json.RawMessage `json:"payload" db:"payload"`
	Debug     any             `json:"debug" db:"debug"`
}

type Meta struct {
	Company     string `json:"company"`
	Team        string `json:"team"`
	Environment string `json:"environment"`
}

func NewEvent(scope Scope, object string, meta *Meta, payload any, eventType Type, debug any) (*Event, error) {
	family, err := GetFamilyForType(eventType)
	if err != nil {
		return nil, errors.Wrap(err, "error decoding event family from event type")
	}
	return &Event{
		ID:        uuid.New().String(),
		Object:    object,
		Scope:     scope,
		Meta:      util.ToJSONBytes(meta, false),
		CreatedAt: time.Now(),
		EventType: eventType,
		Family:    family,
		Payload:   util.ToJSONBytes(payload, false),
		Debug:     debug,
	}, nil
}

type RecordPayload struct {
	Scope     Scope  `json:"scope"`
	Object    string `json:"object"`
	Meta      *Meta  `json:"meta"`
	EventType string `json:"eventType"`
	Payload   any    `json:"payload"`
	Debug     any    `json:"debug"`
}

type ListPayload struct {
	Scope       Scope  `json:"scope"`
	Filter      Filter `json:"filter"`
	Company     string `json:"company"`
	Team        string `json:"team"`
	Environment string `json:"environment"`
}
