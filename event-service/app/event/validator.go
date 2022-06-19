package event

import (
	"github.com/pkg/errors"
)

func validateRecordPayload(e *RecordPayload) error {
	if e == nil {
		return errors.New("record payload is nil")
	}

	if err := validateEventType(e.EventType); err != nil {
		return err
	}
	if e.Company == "" {
		return errors.New("company cannot be empty")
	}
	if e.Team == "" {
		return errors.New("team cannot be empty")
	}
	if e.Environment == "" {
		return errors.New("environment cannot be empty")
	}

	return nil
}

func validateEventType(eventType string) error {
	if eventType == "" {
		return errors.New("eventType cannot be empty")
	}

	if !IsValidationEvent(Type(eventType)) {
		return errors.Errorf("unsupported event type: %s", eventType)
	}

	return nil
}

func validateListPayload(p *ListPayload, scope Scope) error {
	switch scope {
	case ScopeCompany:
		if p.Company == "" {
			return errors.New("company must be defined when scope is set to company")
		}
	case ScopeTeam:
		if p.Company == "" {
			return errors.New("company must be defined when scope is set to team")
		}
		if p.Team == "" {
			return errors.New("company must be defined when scope is set to team")
		}
	case ScopeEnvironment:
		if p.Company == "" {
			return errors.New("company must be defined when scope is set to environment")
		}
		if p.Team == "" {
			return errors.New("company must be defined when scope is set to environment")
		}
		if p.Environment == "" {
			return errors.New("company must be defined when scope is set to environment")
		}
	default:
		return errors.Errorf("invalid scope: %s", scope)
	}
	return nil
}
