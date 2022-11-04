package zlog

import (
	"context"

	"github.com/compuzest/zlifecycle-state-manager/app/env"

	"github.com/newrelic/go-agent/v3/integrations/logcontext/nrlogrusplugin"
	log "github.com/sirupsen/logrus"
)

var logger = initLogger()

func PlainLogger() *log.Logger {
	return logger
}

func CtxLogger(ctx context.Context) *log.Entry {
	return logger.WithContext(ctx)
}

func initLogger() *log.Logger {
	l := log.New()
	if env.Config().DevMode != "true" {
		l.SetFormatter(nrlogrusplugin.ContextFormatter{})
	} else {
		l.SetFormatter(&log.TextFormatter{})
	}

	return l
}
