package controllers

import (
	"context"
	"net/http"

	"github.com/compuzest/zlifecycle-event-service/app/event"
	"github.com/compuzest/zlifecycle-event-service/app/services"

	"github.com/compuzest/zlifecycle-event-service/app/apm"
	http2 "github.com/compuzest/zlifecycle-event-service/app/web/http"
	"github.com/compuzest/zlifecycle-event-service/app/zlog"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/pkg/errors"
)

func EventsHandler(svcs *services.Services) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		txn := newrelic.FromContext(r.Context())

		var err error
		var resp interface{}
		var statusCode int
		switch r.Method {
		case http.MethodGet:
			resp, err = getEventsHandler(r.Context(), r, svcs)
			statusCode = http.StatusOK
		default:
			err := apm.NoticeError(
				txn,
				http2.NewVerboseError("NotFoundError", r.Method, "/events", errors.New("endpoint not implemented")),
			)
			zlog.CtxLogger(r.Context()).Error(err)
			http2.ErrorResponse(w, err.Error(), http.StatusNotFound)
			return
		}
		if err != nil {
			verr := http2.NewVerboseError("EventsError", r.Method, "/events", err)
			_ = apm.NoticeError(txn, verr)
			zlog.CtxLogger(r.Context()).WithError(verr).Errorf("events handler error")
			zlog.CtxLogger(r.Context()).Errorf("%+v", verr.OriginalError)
			http2.ErrorResponse(w, verr.Error(), http.StatusBadRequest)
			return
		}

		http2.Response(w, resp, statusCode)
	}
}

func getEventsHandler(ctx context.Context, r *http.Request, svcs *services.Services) (*GetEventsResponse, error) {
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
	events, err := svcs.ES.List(ctx, &p, event.Scope(scope), event.Filter(filter))
	if err != nil {
		return nil, errors.Wrap(err, "error listing events")
	}

	return &GetEventsResponse{events: events}, nil
}

type GetEventsRequest struct{}

type GetEventsResponse struct {
	events []*event.Event
}
