package controllers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/compuzest/zlifecycle-event-service/internal/status"

	"github.com/compuzest/zlifecycle-event-service/internal/util"

	"github.com/compuzest/zlifecycle-event-service/internal/apm"
	"github.com/sirupsen/logrus"

	"github.com/compuzest/zlifecycle-event-service/internal/event"
	"github.com/compuzest/zlifecycle-event-service/internal/services"

	http2 "github.com/compuzest/zlifecycle-event-service/internal/web/http"
	"github.com/compuzest/zlifecycle-event-service/internal/zlog"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/pkg/errors"
)

func EventsHandler(svcs *services.Services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		txn := newrelic.FromContext(r.Context())
		log := zlog.NewCtxEntry(r.Context())
		var err error
		var resp any
		var statusCode int
		switch r.Method {
		case http.MethodGet:
			resp, err = getEventsHandler(r.Context(), r, svcs, log)
			statusCode = http.StatusOK
		case http.MethodPost:
			resp, err = postEventsHandler(r.Context(), r, svcs, log)
			statusCode = http.StatusCreated
		default:
			err := apm.NoticeError(txn, http2.NewNotFoundError(r))
			http2.WriteNotFoundError(err, w, log)
			return
		}
		if err != nil {
			verr := apm.NoticeError(txn, http2.NewVerboseError("EventsError", r, err))
			http2.WriteInternalError(w, verr, r, log)
			return
		}

		http2.WriteResponse(w, resp, statusCode)
	}
}

func getEventsHandler(ctx context.Context, r *http.Request, svcs *services.Services, log *logrus.Entry) (*GetEventsResponse, error) {
	filter := r.URL.Query().Get("filter")
	scope := r.URL.Query().Get("scope")
	company := r.URL.Query().Get("company")
	team := r.URL.Query().Get("team")
	environment := r.URL.Query().Get("environment")

	p := event.ListPayload{
		Company:     company,
		Team:        team,
		Environment: environment,
		Scope:       event.Scope(scope),
		Filter:      event.Filter(filter),
	}
	events, err := svcs.ES.List(ctx, &p, log)
	if err != nil {
		return nil, errors.Wrap(err, "error listing events")
	}

	return &GetEventsResponse{Events: events}, nil
}

type GetEventsRequest struct{}

type GetEventsResponse struct {
	Events []*event.Event `json:"events"`
}

func postEventsHandler(ctx context.Context, r *http.Request, svcs *services.Services, log *logrus.Entry) (*PostEventsResponse, error) {
	var body PostEventsRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&body); err != nil {
		return nil, errors.Wrap(err, "invalid post events request body")
	}

	var recordMeta event.Meta
	if err := util.CycleJSON(body.Meta, &recordMeta); err != nil {
		return nil, errors.Wrap(err, "error unmarshalling meta field")
	}
	p := event.RecordPayload{
		Scope:     event.Scope(body.Scope),
		Object:    body.Object,
		Meta:      &recordMeta,
		EventType: body.EventType,
		Payload:   body.Payload,
		Debug:     body.Debug,
	}
	evt, err := svcs.ES.Record(ctx, &p, log)
	if err != nil {
		return nil, errors.Wrap(
			err,
			"error recording event",
		)
	}

	status, err := status.NewEnvironmentStatus([]*event.Event{evt}, 1)
	if err != nil {
		return nil, errors.Wrapf(
			err,
			"error generating event status for company [%s], team [%s] and environment [%s]",
			status.Company, status.Team, status.Environment,
		)
	}

	msg := util.ToJSONBytes(status, false)
	svcs.SSEBroker.Notify(msg)
	log.Infof("Successfully notified %d client(s) for new event", svcs.SSEBroker.Clients())

	return &PostEventsResponse{*evt}, nil
}

type PostEventsRequest struct {
	Scope     string `json:"scope"`
	Object    string `json:"object"`
	Meta      any    `json:"meta"`
	EventType string `json:"eventType"`
	Payload   any    `json:"payload"`
	Debug     any    `json:"debug"`
}

type PostEventsResponse struct {
	event.Event
}
