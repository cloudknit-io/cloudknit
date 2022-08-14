package controllers

import (
	"context"
	"github.com/compuzest/zlifecycle-event-service/internal/apm"
	"github.com/compuzest/zlifecycle-event-service/internal/services"
	http2 "github.com/compuzest/zlifecycle-event-service/internal/web/http"
	"github.com/compuzest/zlifecycle-event-service/internal/zlog"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sirupsen/logrus"
	"net/http"
)

func AdminDatabaseHandler(svcs *services.Services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		txn := newrelic.FromContext(r.Context())
		log := zlog.NewCtxEntry(r.Context())
		var err error
		var resp any
		var statusCode int
		switch r.Method {
		case http.MethodDelete:
			resp, err = deleteAdminDatabaseHandler(r.Context(), r, svcs, log)
			statusCode = http.StatusNoContent
		default:
			err := apm.NoticeError(txn, http2.NewNotFoundError(r))
			http2.WriteNotFoundError(err, w, log)
			return
		}
		if err != nil {
			verr := apm.NoticeError(txn, http2.NewVerboseError("AdminDatabaseError", r, err))
			http2.WriteInternalError(w, verr, r, log)
			return
		}

		http2.WriteResponse(w, resp, statusCode)
	}
}

func deleteAdminDatabaseHandler(ctx context.Context, r *http.Request, svcs *services.Services, log *logrus.Entry) (any, error) {
	if err := svcs.AS.WipeDatabase(); err != nil {
		return nil, err
	}

	return struct{}{}, nil
}
