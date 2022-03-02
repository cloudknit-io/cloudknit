package apm

import (
	"context"

	"github.com/compuzest/zlifecycle-il-operator/controllers/zerrors"
	nr "github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sirupsen/logrus"
)

type Noop struct{}

func NewNoop() *Noop {
	return &Noop{}
}

func (n *Noop) NoticeError(tx *nr.Transaction, log *logrus.Entry, err zerrors.ZError) error {
	log.WithError(err).Errorf("error during reconcile\nerror: %s\nstack trace:%+v", err.Error(), deepestStackTrace(err.OriginalError()))
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
