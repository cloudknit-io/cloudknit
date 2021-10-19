package uinotifier

import "github.com/go-logr/logr"

type UINotifier struct {
	apiURL string
	log    logr.Logger
}

func NewUINotifier(log logr.Logger, apiURL string) *UINotifier {
	return &UINotifier{apiURL: apiURL, log: log}
}
