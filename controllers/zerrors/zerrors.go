package zerrors

import (
	"fmt"

	"github.com/compuzest/zlifecycle-il-operator/controllers/util/env"
)

type ZError interface {
	Error() string
	Attributes() map[string]interface{}
	Class() string
	Metric() string
	OriginalError() error
}

type CompanyError struct {
	Company string `json:"company"`
	Err     error  `json:"err"`
}

func (e *CompanyError) Error() string {
	return fmt.Sprintf("company reconcile failed for company [%s]: %v", e.Company, e.Err)
}

func (e *CompanyError) Attributes() map[string]interface{} {
	return map[string]interface{}{
		"company": e.Company,
	}
}

func (e *CompanyError) Class() string {
	return "CompanyReconcilerError"
}

func (e *CompanyError) Metric() string {
	return "com.zlifecycle.companyreconciler.error"
}

func (e *CompanyError) OriginalError() error {
	return e.Err
}

var _ ZError = (*CompanyError)(nil)

func NewCompanyError(company string, err error) *CompanyError {
	return &CompanyError{
		Company: company,
		Err:     err,
	}
}

type TeamError struct {
	Company string `json:"company"`
	Team    string `json:"team"`
	Err     error  `json:"error"`
}

func (e *TeamError) Error() string {
	return fmt.Sprintf("team reconcile failed for company [%s] -> team [%s]: %v", e.Company, e.Team, e.Err)
}

func (e *TeamError) Attributes() map[string]interface{} {
	return map[string]interface{}{
		"company": e.Company,
		"team":    e.Team,
	}
}

func (e *TeamError) Class() string {
	return "TeamReconcilerError"
}

func (e *TeamError) Metric() string {
	return "com.zlifecycle.teamreconciler.error"
}

func (e *TeamError) OriginalError() error {
	return e.Err
}

var _ ZError = (*TeamError)(nil)

func NewTeamError(team string, err error) *TeamError {
	return &TeamError{
		Team: team,
		Err:  err,
	}
}

type EnvironmentError struct {
	Company     string `json:"company"`
	Team        string `json:"team"`
	Environment string `json:"environment"`
	Err         error  `json:"error"`
}

func (e *EnvironmentError) Error() string {
	return fmt.Sprintf("environment reconcile failed for company [%s] -> team [%s] -> environment [%s]: %v", e.Company, e.Team, e.Environment, e.Err)
}

func (e *EnvironmentError) Attributes() map[string]interface{} {
	return map[string]interface{}{
		"company":     env.Config.CompanyName,
		"team":        e.Team,
		"environment": e.Environment,
	}
}

func (e *EnvironmentError) Class() string {
	return "EnvironmentReconcilerError"
}

func (e *EnvironmentError) Metric() string {
	return "com.zlifecycle.environmentreconciler.error"
}

func (e *EnvironmentError) OriginalError() error {
	return e.Err
}

var _ ZError = (*EnvironmentError)(nil)

func NewEnvironmentError(team string, environment string, err error) *EnvironmentError {
	return &EnvironmentError{
		Company:     env.Config.CompanyName,
		Team:        team,
		Environment: environment,
		Err:         err,
	}
}

type EnvironmentComponentError struct {
	Component string `json:"component"`
	Err       error  `json:"error"`
}

func (e *EnvironmentComponentError) Error() string {
	return fmt.Sprintf("error reconciling environment component [%s]: %v", e.Component, e.Err)
}

func NewEnvironmentComponentError(component string, err error) error {
	return &EnvironmentComponentError{
		Component: component,
		Err:       err,
	}
}
