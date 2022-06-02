package event

import "time"

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

type ListPayload struct {
	Company     string `json:"company" db:"company"`
	Team        string `json:"team" db:"team"`
	Environment string `json:"environment" db:"environment"`
}
