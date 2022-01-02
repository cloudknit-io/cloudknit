package log

import (
	"github.com/compuzest/zlifecycle-internal-cli/app/env"
	"github.com/sirupsen/logrus"
	"io"
)

func NewLogger() *logrus.Logger {
	l := logrus.New()
	if !env.Verbose {
		l.Out = io.Discard
	}
	return l
}
