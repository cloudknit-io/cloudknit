package event

import (
	"github.com/pkg/errors"
)

func validateRecordPayload(p *RecordPayload) error {
	if p == nil {
		return errors.New("record payload is nil")
	}

	if err := validateEventType(p.EventType); err != nil {
		return err
	}

	if p.Object == "" {
		return errors.New("object field cannot be empty")
	}

	if p.Meta == nil {
		return errors.New("record meta is nil")
	}

	if err := validateScope(p.Scope, p.Meta); err != nil {
		return err
	}

	return nil
}

func validateScope(scope Scope, meta *Meta) error {
	switch scope {
	case ScopeCompany:
		if meta.Company == "" {
			return errors.New("meta.company field cannot be empty")
		}
	case ScopeTeam:
		if meta.Company == "" {
			return errors.New("meta.company field cannot be empty")
		}
		if meta.Team == "" {
			return errors.New("meta.team field cannot be empty")
		}
	case ScopeEnvironment:
		if meta.Company == "" {
			return errors.New("meta.company field cannot be empty")
		}
		if meta.Team == "" {
			return errors.New("meta.team field cannot be empty")
		}
		if meta.Environment == "" {
			return errors.New("meta.environment field cannot be empty")
		}
	default:
		return errors.Errorf("invalid scope: %s", scope)
	}

	return nil
}

func validateEventType(eventType string) error {
	if eventType == "" {
		return errors.New("eventType cannot be empty")
	}

	if !isSupportedEvent(Type(eventType)) {
		return errors.Errorf("unsupported event type: %s", eventType)
	}

	return nil
}

func validateListPayload(p *ListPayload) error {
	switch p.Scope {
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
		return errors.Errorf("invalid scope: %s", p.Scope)
	}
	return nil
}
