package event

import (
	"encoding/json"
	"time"

	"github.com/pkg/errors"

	"github.com/compuzest/zlifecycle-event-service/internal/util"

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
	TeamSpecValidationSuccess          Type   = "team_validation_success"
	TeamSpecValidationError            Type   = "team_validation_error"
	TeamSchemaValidationError          Type   = "team_schema_validation_error"
	TeamSchemaValidationSuccess        Type   = "team_schema_validation_success"
	FamilyValidationError              Family = "validation_error"
	FamilyValidationSuccess            Family = "validation_success"
)

type Event struct {
	ID          string          `json:"id" db:"id"`
	Scope       Scope           `json:"scope" db:"scope"`
	Object      string          `json:"object" db:"object"`
	Meta        json.RawMessage `json:"meta" db:"meta"`
	CreatedAt   time.Time       `json:"createdAt" db:"created_at"`
	EventType   Type            `json:"eventType" db:"event_type"`
	EventFamily Family          `json:"eventFamily" db:"event_family"`
	Payload     json.RawMessage `json:"payload" db:"payload"`
	Debug       any             `json:"debug" db:"debug"`
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
		ID:          uuid.New().String(),
		Object:      object,
		Scope:       scope,
		Meta:        util.ToJSONBytes(meta, false),
		CreatedAt:   time.Now(),
		EventType:   eventType,
		EventFamily: family,
		Payload:     util.ToJSONBytes(payload, false),
		Debug:       debug,
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
