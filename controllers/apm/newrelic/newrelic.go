package newrelic

import (
	"context"
	"fmt"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/env"
	"github.com/compuzest/zlifecycle-il-operator/controllers/zerrors"
	nr "github.com/newrelic/go-agent/v3/newrelic"
	"github.com/pkg/errors"
)

type APM interface {
	NoticeError(tx *nr.Transaction, err zerrors.ZError) error
	RecordCustomEvent(event string, params map[string]interface{})
	StartTransaction(name string) *nr.Transaction
	NewContext(ctx context.Context, tx *nr.Transaction) context.Context
}

type App struct {
	app *nr.Application
}

func NewApp() (*App, error) {
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
				"env": env.Config.CompanyName,
			}
		},
	)
	if err != nil {
		return nil, err
	}

	return &App{app: app}, nil
}

func (a *App) NoticeError(tx *nr.Transaction, err zerrors.ZError) error {
	a.app.RecordCustomMetric(err.Metric(), 1)
	if tx != nil {
		tx.NoticeError(nr.Error{
			Message:    err.Error(),
			Class:      err.Class(),
			Attributes: err.Attributes(),
		})
	}

	return err
}

func (a *App) RecordCustomEvent(event string, params map[string]interface{}) {
	fullEvent := fmt.Sprintf("com.zlifecycle.%s", event)
	a.app.RecordCustomEvent(fullEvent, params)
}

func (a *App) StartTransaction(name string) *nr.Transaction {
	fullName := fmt.Sprintf("com.zlifecycle.%s", name)
	return a.app.StartTransaction(fullName)
}

func (a *App) NewContext(ctx context.Context, tx *nr.Transaction) context.Context {
	if tx != nil {
		return nr.NewContext(ctx, tx)
	}
	return ctx
}

var _ APM = (*App)(nil)

type Noop struct{}

func NewNoop() *Noop {
	return &Noop{}
}

func (n *Noop) NoticeError(tx *nr.Transaction, err zerrors.ZError) error {
	return err
}

func (n *Noop) RecordCustomEvent(event string, params map[string]interface{}) {}

func (n *Noop) StartTransaction(name string) *nr.Transaction {
	return nil
}

func (n *Noop) NewContext(ctx context.Context, tx *nr.Transaction) context.Context {
	return ctx
}

var _ APM = (*Noop)(nil)
