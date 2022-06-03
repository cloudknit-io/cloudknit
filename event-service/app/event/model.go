package event

import (
	"time"

	"github.com/google/uuid"
)

type Type string

const (
	ValidationSuccess Type = "validation_success"
	ValidationError   Type = "validation_error"
)

type Event struct {
	ID          string    `json:"id" db:"id"`
	Company     string    `json:"company" db:"company"`
	Team        string    `json:"team" db:"team"`
	Environment string    `json:"environment" db:"environment"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	EventType   Type      `json:"eventType" db:"event_type"`
	Message     string    `json:"message" db:"message"`
	Debug       any       `json:"debug" db:"debug"`
}

func NewEvent(company, team, environment, message string, eventType Type, debug any) *Event {
	return &Event{
		ID:          uuid.New().String(),
		Company:     company,
		Team:        team,
		Environment: environment,
		CreatedAt:   time.Now(),
		EventType:   eventType,
		Message:     message,
		Debug:       debug,
	}
}

type RecordPayload struct {
	Company     string `json:"company"`
	Team        string `json:"team"`
	Environment string `json:"environment"`
	EventType   string `json:"eventType"`
	Message     string `json:"message"`
	Debug       any    `json:"debug"`
}

type ListPayload struct {
	Company     string `json:"company" db:"company"`
	Team        string `json:"team" db:"team"`
	Environment string `json:"environment" db:"environment"`
}
