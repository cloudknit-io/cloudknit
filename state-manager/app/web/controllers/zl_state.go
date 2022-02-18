package controllers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/compuzest/zlifecycle-state-manager/app/apm"
	http2 "github.com/compuzest/zlifecycle-state-manager/app/web/http"
	"github.com/compuzest/zlifecycle-state-manager/app/zlog"
	"github.com/compuzest/zlifecycle-state-manager/app/zlstate"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func ZLStateHandler(w http.ResponseWriter, r *http.Request) {
	txn := newrelic.FromContext(r.Context())

	var err error
	var resp interface{}
	var statusCode int

	s3Client, err := zlstate.NewS3Client(r.Context())
	if err != nil {
		err = apm.NoticeError(
			txn,
			http2.NewVerboseError("InternalError", r.Method, "/zl/state", errors.New("internal server error")),
		)
		zlog.CtxLogger(r.Context()).Error(err)
		http2.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	switch r.Method {
	case "POST":
		resp, err = postZLStateHandler(r.Context(), zlog.CtxLogger(r.Context()), r.Body, s3Client)
		statusCode = http.StatusOK
	case "PUT":
		resp, err = putZLStateHandler(r.Context(), zlog.CtxLogger(r.Context()), r.Body, s3Client)
		statusCode = http.StatusOK
	default:
		err := apm.NoticeError(
			txn,
			http2.NewVerboseError("NotFoundError", r.Method, "/zl/state", errors.New("endpoint not implemented")),
		)
		zlog.CtxLogger(r.Context()).Error(err)
		http2.ErrorResponse(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		verr := http2.NewVerboseError("ZLStateError", r.Method, "/zl/state", err)
		_ = apm.NoticeError(txn, verr)
		zlog.CtxLogger(r.Context()).WithError(verr).Errorf("zl state handler error")
		zlog.CtxLogger(r.Context()).Errorf("%+v", verr.OriginalError)
		http2.ErrorResponse(w, verr.Error(), http.StatusBadRequest)
		return
	}

	http2.Response(w, resp, statusCode)
}

func postZLStateHandler(ctx context.Context, log *logrus.Entry, b io.ReadCloser, s3Client zlstate.S3API) (*FetchZLStateResponse, error) {
	var body FetchZLStateRequest
	decoder := json.NewDecoder(b)
	if err := decoder.Decode(&body); err != nil {
		return nil, errors.Wrap(err, "invalid get zLstate body")
	}
	if err := validateGetZLStateRequest(&body); err != nil {
		return nil, errors.Wrap(err, "error validating get zLstate request body")
	}

	log.WithField("body", body).Info("Handling get zLstate request")

	backend := zlstate.NewS3Backend(ctx, log, BuildZLStateBucketName(body.Company), s3Client)

	zlState, err := backend.Get(BuildZLStateKey(body.Team, body.Environment))
	if err != nil {
		return nil, errors.Wrap(err, "error getting zLstate from remote backend")
	}

	return &FetchZLStateResponse{ZLState: zlState}, nil
}

type FetchZLStateRequest struct {
	Company     string `json:"company"`
	Team        string `json:"team"`
	Environment string `json:"environment"`
}

type FetchZLStateResponse struct {
	ZLState *zlstate.ZLState `json:"zlstate"`
}

func putZLStateHandler(ctx context.Context, log *logrus.Entry, b io.ReadCloser, s3Client zlstate.S3API) (*PutZLStateResponse, error) {
	var body PutZLStateRequest
	decoder := json.NewDecoder(b)
	if err := decoder.Decode(&body); err != nil {
		return nil, errors.Wrap(err, "invalid put zLstate body")
	}
	if err := validatePutZLStateRequest(&body); err != nil {
		return nil, errors.Wrap(err, "error validating put zLstate request body")
	}

	log.WithField("body", body).Info("Handling put zLstate request")

	client := zlstate.NewS3Backend(ctx, log, BuildZLStateBucketName(body.Company), s3Client)

	if err := client.Put(BuildZLStateKey(body.Team, body.Environment), body.ZLState, false); err != nil {
		if errors.Is(err, zlstate.ErrKeyAlreadyExists) {
			return &PutZLStateResponse{Message: "zLstate already exists"}, nil
		}
		return nil, errors.Wrap(err, "error persisting zLstate to remote backend")
	}

	return &PutZLStateResponse{Message: "zLstate created successfully"}, nil
}

type PutZLStateRequest struct {
	Company     string           `json:"company"`
	Team        string           `json:"team"`
	Environment string           `json:"environment"`
	ZLState     *zlstate.ZLState `json:"zlstate"`
}

type PutZLStateResponse struct {
	Message string `json:"message"`
}
