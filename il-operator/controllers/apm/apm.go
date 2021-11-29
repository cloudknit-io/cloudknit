package apm

import (
	"context"
	"github.com/compuzest/zlifecycle-il-operator/controllers/zerrors"
	nr "github.com/newrelic/go-agent/v3/newrelic"
)

type APM interface {
	NoticeError(tx *nr.Transaction, err zerrors.ZError) error
	RecordCustomEvent(event string, params map[string]interface{})
	StartTransaction(name string) *nr.Transaction
	NewContext(ctx context.Context, tx *nr.Transaction) context.Context
}
