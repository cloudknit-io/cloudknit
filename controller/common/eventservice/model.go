package eventservice

type (
	Scope string
	Type  string
)

const (
	ScopeCompany                       Scope = "company"
	ScopeTeam                          Scope = "team"
	ScopeEnvironment                   Scope = "environment"
	EnvironmentValidationSuccess       Type  = "environment_validation_success"
	EnvironmentValidationError         Type  = "environment_validation_error"
	EnvironmentSchemaValidationError   Type  = "environment_schema_validation_error"
	EnvironmentSchemaValidationSuccess Type  = "environment_schema_validation_success"
	EnvironmentReconcileError          Type  = "environment_reconcile_error"
	EnvironmentReconcileSuccess        Type  = "environment_reconcile_success"
	TeamValidationSuccess              Type  = "team_validation_success"
	TeamValidationError                Type  = "team_validation_error"
	TeamSchemaValidationError          Type  = "team_schema_validation_error"
	TeamSchemaValidationSuccess        Type  = "team_schema_validation_success"
	TeamReconcileError                 Type  = "team_reconcile_error"
	TeamReconcileSuccess               Type  = "team_reconcile_success"
)

type Event struct {
	Scope     string `json:"scope"`
	Object    string `json:"object"`
	Meta      *Meta  `json:"meta"`
	EventType string `json:"eventType"`
	Payload   any    `json:"payload"`
	Debug     any    `json:"debug"`
}

type Meta struct {
	Company     string `json:"company"`
	Team        string `json:"team,omitempty"`
	Environment string `json:"environment,omitempty"`
}
