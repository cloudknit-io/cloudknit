package zlstate

import (
	v1 "github.com/compuzest/zlifecycle-il-operator/api/v1"
	"github.com/compuzest/zlifecycle-il-operator/controller/common/state_manager"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func ReconcileState(api state_manager.API, company, team string, environment *v1.Environment, log *logrus.Entry) error {
	resp, err := api.Get(company, team, environment.Spec.EnvName)
	if err != nil {
		return errors.Wrap(err, "error fetching zlstate")
	}

	for _, ec := range environment.Spec.Components {
		if !componentExists(resp.ZLState.Components, ec.Name) {
			log.Infof(
				"Adding new component [%s] to company [%s], team [%s] and environment [%s] zL state",
				ec.Name, company, team, environment.Spec.EnvName,
			)
			if err := api.PutComponent(company, team, environment.Spec.EnvName, state_manager.ToZLStateComponent(ec)); err != nil {
				return errors.Wrap(err, "error adding component to zlstate")
			}
		}
	}

	return nil
}

func componentExists(components []*state_manager.Component, component string) bool {
	for _, c := range components {
		if c.Name == component {
			return true
		}
	}
	return false
}
