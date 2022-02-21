package log

import (
	"github.com/compuzest/zlifecycle-il-operator/controllers/env"
	"github.com/newrelic/go-agent/v3/integrations/logcontext/nrlogrusplugin"
	"github.com/sirupsen/logrus"
)

func NewLogger() *logrus.Logger {
	l := logrus.New()
	if env.Config.Mode != "local" {
		l.SetFormatter(nrlogrusplugin.ContextFormatter{})
	} else {
		l.SetFormatter(&logrus.TextFormatter{})
	}
	return l
}
