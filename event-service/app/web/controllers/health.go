package controllers

import (
	"context"
	"github.com/compuzest/zlifecycle-event-service/app/apm"
	"github.com/compuzest/zlifecycle-event-service/app/health"
	"github.com/compuzest/zlifecycle-event-service/app/services"
	http2 "github.com/compuzest/zlifecycle-event-service/app/web/http"
	"github.com/compuzest/zlifecycle-event-service/app/zlog"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net/http"
)

func HealthHandler(svcs *services.Services) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		txn := newrelic.FromContext(r.Context())
		log := zlog.CtxLogger(r.Context())
		var err error
		var resp any
		var statusCode int
		switch r.Method {
		case http.MethodGet:
			resp, err = getHealthHandler(r.Context(), r, svcs, log)
			statusCode = http.StatusOK
		default:
			err := apm.NoticeError(txn, http2.NewNotFoundError(r))
			http2.WriteNotFoundError(err, w, log)
			return
		}
		if err != nil {
			verr := apm.NoticeError(txn, http2.NewVerboseError("HealthError", r, err))
			http2.WriteInternalError(w, verr, r, log)
			return
		}

		http2.WriteResponse(w, resp, statusCode)
	}
}

func getHealthHandler(ctx context.Context, r *http.Request, svcs *services.Services, log *logrus.Entry) (*GetHealthResponse, error) {
	company := r.URL.Query().Get("company")
	if company == "" {
		return nil, errors.New("missing query param: company")
	}

	healthcheck, err := svcs.HS.Healthcheck(ctx, company, log)
	if err != nil {
		return nil, errors.Wrap(err, "error inspecting company health")
	}

	return &GetHealthResponse{Healthcheck: healthcheck}, nil
}

type GetHealthResponse struct {
	Healthcheck health.TeamStatus `json:"health"`
}
