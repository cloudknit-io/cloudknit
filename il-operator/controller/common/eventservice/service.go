package eventservice

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/compuzest/zlifecycle-il-operator/controller/util"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

//go:generate mockgen --build_flags=--mod=mod -destination=./mock_event_service.go -package=eventservice "github.com/compuzest/zlifecycle-il-operator/controller/common/eventservice" API
type API interface {
	Record(ctx context.Context, e *Event, log *logrus.Entry) error
}

type Service struct {
	apiURL     string
	httpClient *http.Client
}

func NewService(apiURL string) *Service {
	return &Service{apiURL: apiURL, httpClient: util.GetHTTPClient()}
}

func (s *Service) Record(ctx context.Context, e *Event, log *logrus.Entry) error {
	log = log.WithFields(logrus.Fields{
		"eventType": e.EventType,
	})
	log.Infof("Sending [%s] event to event service", e.EventType)

	eventsEndpoint := fmt.Sprintf("%s/events", s.apiURL)

	jsonBody, err := util.ToJSON(e)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, eventsEndpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send POST request to %s: %w", eventsEndpoint, err)
	}
	defer util.CloseBody(resp.Body)

	if resp.StatusCode != http.StatusCreated {
		return errors.Errorf("received non-OK status code: %d", resp.StatusCode)
	}

	return nil
}

var _ API = &Service{}
