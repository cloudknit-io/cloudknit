package controllers

import (
	"context"
	"net/http"

	"github.com/cloudknit-io/cloudknit/event-service/internal/apm"
	"github.com/cloudknit-io/cloudknit/event-service/internal/health"
	"github.com/cloudknit-io/cloudknit/event-service/internal/services"
	http2 "github.com/cloudknit-io/cloudknit/event-service/internal/web/http"
	"github.com/cloudknit-io/cloudknit/event-service/internal/zlog"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sirupsen/logrus"
)

func HealthHandler(svcs *services.Services, fullCheck bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		txn := newrelic.FromContext(r.Context())
		log := zlog.NewCtxEntry(r.Context())
		var err error
		var resp any
		var statusCode int
		switch r.Method {
		case http.MethodGet:
			hc := getHealthHandler(r.Context(), svcs, fullCheck, log)
			statusCode = hc.Code
			resp = hc
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

func getHealthHandler(ctx context.Context, svcs *services.Services, fullCheck bool, log *logrus.Entry) *GetHealthResponse {
	hc := svcs.HS.Healthcheck(ctx, fullCheck, log)

	return &GetHealthResponse{*hc}
}

type GetHealthResponse struct {
	health.Healthcheck
}
