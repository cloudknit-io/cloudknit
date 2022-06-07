package controllers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/compuzest/zlifecycle-event-service/app/util"

	"github.com/compuzest/zlifecycle-event-service/app/apm"
	"github.com/compuzest/zlifecycle-event-service/app/health"

	"github.com/sirupsen/logrus"

	"github.com/compuzest/zlifecycle-event-service/app/event"
	"github.com/compuzest/zlifecycle-event-service/app/services"

	http2 "github.com/compuzest/zlifecycle-event-service/app/web/http"
	"github.com/compuzest/zlifecycle-event-service/app/zlog"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/pkg/errors"
)

func EventsHandler(svcs *services.Services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		txn := newrelic.FromContext(r.Context())
		log := zlog.CtxLogger(r.Context())
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
	}
	events, err := svcs.ES.List(ctx, &p, event.Scope(scope), event.Filter(filter), log)
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

	p := event.RecordPayload{
		Company:     body.Company,
		Team:        body.Team,
		Environment: body.Environment,
		EventType:   body.EventType,
		Payload:     body.Payload,
		Debug:       body.Debug,
	}
	evt, err := svcs.ES.Record(ctx, &p, log)
	if err != nil {
		return nil, errors.Wrapf(
			err,
			"error recording event for company [%s], team [%s] and environment [%s]",
			p.Company, p.Team, p.Environment,
		)
	}

	status, err := health.NewEnvironmentStatus([]*event.Event{evt})
	if err != nil {
		return nil, errors.Wrapf(
			err,
			"error generating event healthcheck status for company [%s], team [%s] and environment [%s]",
			p.Company, p.Team, p.Environment,
		)
	}

	msg := util.ToJSONBytes(status, false)
	svcs.SSEBroker.Notify(msg)
	log.Infof("Successfully notified %d client(s) for new event", svcs.SSEBroker.Clients())

	return &PostEventsResponse{*evt}, nil
}

type PostEventsRequest struct {
	Company     string `json:"company"`
	Team        string `json:"team"`
	Environment string `json:"environment"`
	EventType   string `json:"eventType"`
	Payload     any    `json:"payload"`
	Debug       any    `json:"debug"`
}

type PostEventsResponse struct {
	event.Event
}
