package controllers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/cloudknit-io/cloudknit/event-service/internal/status"

	"github.com/cloudknit-io/cloudknit/event-service/internal/apm"
	"github.com/cloudknit-io/cloudknit/event-service/internal/services"
	http2 "github.com/cloudknit-io/cloudknit/event-service/internal/web/http"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func StatusHandler(svcs *services.Services, baseLogger *logrus.Entry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		txn := newrelic.FromContext(r.Context())
		log := baseLogger.WithContext(r.Context())
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

	status, err := svcs.SS.Calculate(ctx, company, history, log)
	if err != nil {
		return nil, errors.Wrap(err, "error inspecting company status")
	}

	return &GetStatusResponse{Status: status}, nil
}

type GetStatusResponse struct {
	// The status model
	// in: body
	Status *status.Response `json:"status"`
}
