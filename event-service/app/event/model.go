package event

import (
	"encoding/json"
	"time"

	"github.com/compuzest/zlifecycle-event-service/app/util"

	"github.com/google/uuid"
)

type Type string

const (
	ValidationSuccess Type = "validation_success"
	ValidationError   Type = "validation_error"
)

type Event struct {
	ID          string          `json:"id" db:"id"`
	Object      string          `json:"object" db:"object"`
	Company     string          `json:"company" db:"company"`
	Team        string          `json:"team" db:"team"`
	Environment string          `json:"environment" db:"environment"`
	CreatedAt   time.Time       `json:"createdAt" db:"created_at"`
	EventType   Type            `json:"eventType" db:"event_type"`
	Payload     json.RawMessage `json:"payload" db:"payload"`
	Debug       any             `json:"debug" db:"debug"`
}

func NewEvent(object, company, team, environment string, payload any, eventType Type, debug any) *Event {
	return &Event{
		ID:          uuid.New().String(),
		Object:      object,
		Company:     company,
		Team:        team,
		Environment: environment,
		CreatedAt:   time.Now(),
		EventType:   eventType,
		Payload:     util.ToJSONBytes(payload, false),
		Debug:       debug,
	}
}

type RecordPayload struct {
	Object      string `json:"object"`
	Company     string `json:"company"`
	Team        string `json:"team"`
	Environment string `json:"environment"`
	EventType   string `json:"eventType"`
	Payload     any    `json:"payload"`
	Debug       any    `json:"debug"`
}

type ListPayload struct {
	Company     string `json:"company" db:"company"`
	Team        string `json:"team" db:"team"`
	Environment string `json:"environment" db:"environment"`
}
