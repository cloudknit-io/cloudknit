package newrelic

import (
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/env"
	nr "github.com/newrelic/go-agent/v3/newrelic"
	"github.com/pkg/errors"
)

var nrapp *nr.Application

func InitNewRelic() (app *nr.Application, err error) {
	if env.Config.NewRelicAPIKey == "" {
		return nil, errors.New("missing NEW_RELIC_API_KEY env var")
	}
	app, err = nr.NewApplication(
		nr.ConfigAppName("zlifecycle-il-operator"),
		nr.ConfigEnabled(env.Config.Mode != "local"),
		nr.ConfigLicense(env.Config.NewRelicAPIKey),
	)
	if err != nil {
		return nil, err
	}

	nrapp = app

	return nrapp, nil
}

func GetApp() *nr.Application {
	return nrapp
}
