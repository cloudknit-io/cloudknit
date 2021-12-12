package apm

import (
	"context"
	"fmt"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/env"
	"github.com/compuzest/zlifecycle-il-operator/controllers/zerrors"
	nr "github.com/newrelic/go-agent/v3/newrelic"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"runtime/debug"
)

type NewRelic struct {
	app *nr.Application
}

func NewNewRelic() (*NewRelic, error) {
	if env.Config.NewRelicAPIKey == "" {
		return nil, errors.New("missing NEW_RELIC_API_KEY env var")
	}
	app, err := nr.NewApplication(
		nr.ConfigAppName(env.Config.App),
		nr.ConfigEnabled(env.Config.EnableNewRelic == "true"),
		nr.ConfigLicense(env.Config.NewRelicAPIKey),
		nr.ConfigDistributedTracerEnabled(true),
		func(c *nr.Config) {
			c.Labels = map[string]string{
				"instance": env.Config.CompanyName,
			}
		},
	)
	if err != nil {
		return nil, err
	}

	return &NewRelic{app: app}, nil
}

func (a *NewRelic) NoticeError(tx *nr.Transaction, log *logrus.Entry, err zerrors.ZError) error {
	a.app.RecordCustomMetric(err.Metric(), 1)
	if tx != nil {
		tx.NoticeError(nr.Error{
			Message:    err.Error(),
			Class:      err.Class(),
			Attributes: err.Attributes(),
		})
	}
	log.WithError(err).Errorf("error during reconcile\nstack trace:\n%s", string(debug.Stack()))

	return err
}

func (a *NewRelic) RecordCustomEvent(event string, params map[string]interface{}) {
	fullEvent := fmt.Sprintf("com.zlifecycle.%s", event)
	a.app.RecordCustomEvent(fullEvent, params)
}

func (a *NewRelic) StartTransaction(name string) *nr.Transaction {
	fullName := fmt.Sprintf("com.zlifecycle.%s", name)
	return a.app.StartTransaction(fullName)
}

func (a *NewRelic) NewContext(ctx context.Context, tx *nr.Transaction) context.Context {
	if tx != nil {
		return nr.NewContext(ctx, tx)
	}
	return ctx
}

var _ APM = (*NewRelic)(nil)
