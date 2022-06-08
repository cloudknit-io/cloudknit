package eventservice

const (
	ValidationSuccess Type = "validation_success"
	ValidationError   Type = "validation_error"
)

type Type string

type Event struct {
	Company     string `json:"company"`
	Team        string `json:"team"`
	Environment string `json:"environment"`
	EventType   string `json:"eventType"`
	Payload     any    `json:"payload"`
	Debug       any    `json:"debug"`
}
