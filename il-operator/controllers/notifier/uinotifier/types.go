package uinotifier

import (
	"github.com/sirupsen/logrus"
)

type UINotifier struct {
	apiURL string
	log    *logrus.Entry
}

func NewUINotifier(log *logrus.Entry, apiURL string) *UINotifier {
	return &UINotifier{apiURL: apiURL, log: log}
}
