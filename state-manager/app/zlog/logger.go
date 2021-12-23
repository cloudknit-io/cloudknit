package zlog

import (
	"context"

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
	l.SetFormatter(nrlogrusplugin.ContextFormatter{})

	return l
}
