package zlstate

import (
	"context"

	v1 "github.com/compuzest/zlifecycle-il-operator/api/v1"
	"github.com/compuzest/zlifecycle-il-operator/controller/common/statemanager"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func ReconcileState(ctx context.Context, api statemanager.API, company, team string, environment *v1.Environment, log *logrus.Entry) error {
	resp, err := api.Get(ctx, company, team, environment.Spec.EnvName, log)
	if err != nil {
		return errors.Wrap(err, "error fetching zlstate")
	}

	for _, ec := range environment.Spec.Components {
		if !componentExists(resp.ZLState.Components, ec.Name) {
			log.Infof(
				"Adding new component [%s] to company [%s], team [%s] and environment [%s] zL state",
				ec.Name, company, team, environment.Spec.EnvName,
			)
			if err := api.PutComponent(ctx, company, team, environment.Spec.EnvName, statemanager.ToZLStateComponent(ec), log); err != nil {
				return errors.Wrap(err, "error adding component to zlstate")
			}
		}
	}

	return nil
}

func componentExists(components []*statemanager.Component, component string) bool {
	for _, c := range components {
		if c.Name == component {
			return true
		}
	}
	return false
}
