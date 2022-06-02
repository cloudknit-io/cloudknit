package apm

import (
	"github.com/compuzest/zlifecycle-event-service/app/env"
	"github.com/compuzest/zlifecycle-event-service/app/web/http"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/pkg/errors"
)

func Init() (*newrelic.Application, error) {
	license := env.Config().NewRelicAPIKey
	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName(env.Config().App),
		newrelic.ConfigLicense(license),
		newrelic.ConfigDistributedTracerEnabled(true),
		func(c *newrelic.Config) {
			c.Labels = map[string]string{
				"instance": env.Config().Instance,
			}
		},
	)
	if err != nil {
		return nil, err
	}

	return app, nil
}

func NoticeError(txn *newrelic.Transaction, err *http.VerboseError) error {
	if txn != nil {
		txn.NoticeError(newrelic.Error{
			Message: err.Error(),
			Class:   err.Class,
			Attributes: map[string]interface{}{
				"instance": env.Config().Instance,
			},
			Stack: stackTrace(err.OriginalError),
		})
	}

	return err
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func stackTrace(err error) []uintptr {
	st := deepestStackTrace(err)
	if st == nil {
		return nil
	}
	return transformStackTrace(st)
}

func deepestStackTrace(err error) errors.StackTrace {
	var last stackTracer
	for err != nil {
		//nolint
		if err, ok := err.(stackTracer); ok {
			last = err
		}
		//nolint
		cause, ok := err.(interface {
			Cause() error
		})
		if !ok {
			break
		}
		err = cause.Cause()
	}

	if last == nil {
		return nil
	}
	return last.StackTrace()
}

func transformStackTrace(orig errors.StackTrace) []uintptr {
	st := make([]uintptr, len(orig))
	for i, frame := range orig {
		st[i] = uintptr(frame)
	}
	return st
}
