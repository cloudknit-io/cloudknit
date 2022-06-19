package controllers

import (
	"context"
	"github.com/compuzest/zlifecycle-event-service/app/status"
	"net/http"
	"strconv"

	"github.com/compuzest/zlifecycle-event-service/app/apm"
	"github.com/compuzest/zlifecycle-event-service/app/services"
	http2 "github.com/compuzest/zlifecycle-event-service/app/web/http"
	"github.com/compuzest/zlifecycle-event-service/app/zlog"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func StatusHandler(svcs *services.Services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		txn := newrelic.FromContext(r.Context())
		log := zlog.CtxLogger(r.Context())
		var err error
		var resp any
		var statusCode int
		switch r.Method {
		case http.MethodGet:
			resp, err = getStatusHandler(r.Context(), r, svcs, log)
			statusCode = http.StatusOK
		default:
			err := apm.NoticeError(txn, http2.NewNotFoundError(r))
			http2.WriteNotFoundError(err, w, log)
			return
		}
		if err != nil {
			verr := apm.NoticeError(txn, http2.NewVerboseError("StatusError", r, err))
			http2.WriteInternalError(w, verr, r, log)
			return
		}

		http2.WriteResponse(w, resp, statusCode)
	}
}

func getStatusHandler(ctx context.Context, r *http.Request, svcs *services.Services, log *logrus.Entry) (*GetStatusResponse, error) {
	company := r.URL.Query().Get("company")
	if company == "" {
		return nil, errors.New("missing query param: company")
	}
	history := 1
	if val := r.URL.Query().Get("history"); val != "" {
		h, err := strconv.Atoi(val)
		if err != nil {
			return nil, errors.Wrapf(err, "invalid value for history parameter, must be integer: %s", val)
		}
		history = h
	}

	status, err := svcs.SS.CompanyStatus(ctx, company, history, log)
	if err != nil {
		return nil, errors.Wrap(err, "error inspecting company status")
	}

	return &GetStatusResponse{Status: status}, nil
}

type GetStatusResponse struct {
	Status status.TeamStatus `json:"status"`
}
