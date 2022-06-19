package eventservice

const (
	ValidationSuccess       Type = "validation_success"
	ValidationError         Type = "validation_error"
	SchemaValidationError   Type = "schema_validation_error"
	SchemaValidationSuccess Type = "schema_validation_success"
)

type Type string

type Event struct {
	Object      string `json:"object"`
	Company     string `json:"company"`
	Team        string `json:"team"`
	Environment string `json:"environment"`
	EventType   string `json:"eventType"`
	Payload     any    `json:"payload"`
	Debug       any    `json:"debug"`
}
