package zlog

import (
	"context"

	"github.com/compuzest/zlifecycle-event-service/internal/env"

	"github.com/newrelic/go-agent/v3/integrations/logcontext/nrlogrusplugin"
	log "github.com/sirupsen/logrus"
)

var logger = initLogger()

func NewLogger() *log.Logger {
	return logger
}

func NewPlainEntry() *log.Entry {
	return log.NewEntry(logger)
}

func NewCtxEntry(ctx context.Context) *log.Entry {
	return logger.WithContext(ctx)
}

func initLogger() *log.Logger {
	l := log.New()
	if env.Config().DevMode != "true" {
		l.SetFormatter(nrlogrusplugin.ContextFormatter{})
	} else {
		l.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	}

	return l
}
