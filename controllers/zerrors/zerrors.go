package zerrors

import "fmt"

type TeamError struct {
	Team string `json:"team"`
	Err  error  `json:"error"`
}

func (e *TeamError) Error() string {
	return fmt.Sprintf("team %s: error %v", e.Team, e.Err)
}

func NewTeamError(team string, err error) error {
	return &TeamError{
		Team: team,
		Err:  err,
	}
}

type EnvironmentError struct {
	Team        string `json:"team"`
	Environment string `json:"environment"`
	Err         error  `json:"error"`
}

func (e *EnvironmentError) Error() string {
	return fmt.Sprintf("environment %s/%s: error %v", e.Team, e.Environment, e.Err)
}

func NewEnvironmentError(team string, environment string, err error) *EnvironmentError {
	return &EnvironmentError{
		Team:        team,
		Environment: environment,
		Err:         err,
	}
}

type EnvironmentComponentError struct {
	Team        string `json:"team"`
	Environment string `json:"environment"`
	Component   string `json:"component"`
	Err         error  `json:"error"`
}

func (e *EnvironmentComponentError) Error() string {
	return fmt.Sprintf("environment %s/%s: component %s: error %v", e.Team, e.Environment, e.Component, e.Err)
}

func NewEnvironmentComponentError(team string, environment string, component string, err error) error {
	return &EnvironmentComponentError{
		Team:        team,
		Environment: environment,
		Component:   component,
		Err:         err,
	}
}
