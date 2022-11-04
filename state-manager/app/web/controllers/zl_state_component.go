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

func ZLStateComponentHandler(w http.ResponseWriter, r *http.Request) {
	txn := newrelic.FromContext(r.Context())

	var err error
	var resp interface{}
	var statusCode int

	s3Client, err := zlstate.NewS3Client(r.Context())
	if err != nil {
		err = apm.NoticeError(
			txn,
			http2.NewVerboseError("InternalError", r.Method, "/zl/state/component", errors.New("internal server error")),
		)
		zlog.CtxLogger(r.Context()).Error(err)
		http2.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	switch r.Method {
	case http.MethodPost:
		resp, err = postZLStateComponentHandler(r.Context(), zlog.CtxLogger(r.Context()), r.Body, s3Client)
		statusCode = http.StatusOK
	case http.MethodPatch:
		resp, err = patchZLStateComponentHandler(r.Context(), zlog.CtxLogger(r.Context()), r.Body, s3Client)
		statusCode = http.StatusOK
	case http.MethodPut:
		resp, err = putZLStateComponentHandler(r.Context(), zlog.CtxLogger(r.Context()), r.Body, s3Client)
		statusCode = http.StatusOK
	case http.MethodDelete:
		resp, err = deleteZLStateComponentHandler(r.Context(), zlog.CtxLogger(r.Context()), r.Body, s3Client)
		statusCode = http.StatusOK
	default:
		err = apm.NoticeError(
			txn,
			http2.NewVerboseError("NotFoundError", r.Method, "/zl/state/component", errors.New("endpoint not implemented")),
		)
		zlog.CtxLogger(r.Context()).Error(err)
		http2.ErrorResponse(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		verr := http2.NewVerboseError("ZLStateComponentError", r.Method, "/zl/state/component", err)
		_ = apm.NoticeError(txn, verr)
		zlog.CtxLogger(r.Context()).WithError(verr).Errorf("zLstate component handler error")
		zlog.CtxLogger(r.Context()).Errorf("%+v", verr.OriginalError)
		http2.ErrorResponse(w, verr.Error(), http.StatusBadRequest)
		return
	}

	http2.Response(w, resp, statusCode)
}

func postZLStateComponentHandler(ctx context.Context, log *logrus.Entry, b io.ReadCloser, s3Client zlstate.S3API) (*PostZLStateComponentResponse, error) {
	var body PostZLStateComponentRequest
	decoder := json.NewDecoder(b)
	if err := decoder.Decode(&body); err != nil {
		return nil, errors.Wrap(err, "invalid get zLstate body")
	}
	if err := validateGetZLStateComponentRequest(&body); err != nil {
		return nil, errors.Wrap(err, "error validating get zLstate component request body")
	}

	log.WithField("body", body).Info("Handling get zLstate component request")

	backend := zlstate.NewS3Backend(ctx, log, BuildZLStateBucketName(body.Company), s3Client)

	zlState, err := backend.Get(BuildZLStateKey(body.Team, body.Environment))
	if err != nil {
		return nil, errors.Wrap(err, "error getting zLstate from remote backend")
	}

	component := findComponent(zlState.Components, body.Component)
	if component == nil {
		return nil, errors.Errorf("zLstate component does not exist: [%s]", body.Component)
	}

	return &PostZLStateComponentResponse{Component: component}, nil
}

func findComponent(components []*zlstate.Component, targetComponent string) *zlstate.Component {
	for _, c := range components {
		if c.Name == targetComponent {
			return c
		}
	}
	return nil
}

type PostZLStateComponentRequest struct {
	Company     string `json:"company"`
	Team        string `json:"team"`
	Environment string `json:"environment"`
	Component   string `json:"component"`
}

type PostZLStateComponentResponse struct {
	Component *zlstate.Component `json:"component"`
}

func patchZLStateComponentHandler(ctx context.Context, log *logrus.Entry, b io.ReadCloser, s3Client zlstate.S3API) (*PatchZLStateComponentResponse, error) {
	var body PatchZLStateComponentRequest
	decoder := json.NewDecoder(b)
	if err := decoder.Decode(&body); err != nil {
		return nil, errors.Wrap(err, "invalid patch zLstate body")
	}
	if err := validatePatchZLStateComponentRequest(&body); err != nil {
		return nil, errors.Wrap(err, "error validating patch zLstate resource body")
	}

	log.WithField("body", body).Info("Handling patch zLstate component status request")

	backend := zlstate.NewS3Backend(ctx, log, BuildZLStateBucketName(body.Company), s3Client)

	key := BuildZLStateKey(body.Team, body.Environment)

	zlst, err := backend.PatchComponent(key, body.Component, body.Status)
	if err != nil {
		return nil, errors.Wrap(err, "error patching component status")
	}

	return &PatchZLStateComponentResponse{ZLState: zlst}, nil
}

type PatchZLStateComponentRequest struct {
	Company     string `json:"company"`
	Team        string `json:"team"`
	Environment string `json:"environment"`
	Component   string `json:"component"`
	Status      string `json:"status"`
}

type PatchZLStateComponentResponse struct {
	ZLState *zlstate.ZLState `json:"zlstate"`
}

func putZLStateComponentHandler(ctx context.Context, log *logrus.Entry, b io.ReadCloser, s3Client zlstate.S3API) (*PutZLStateComponentResponse, error) {
	var body PutZLStateComponentRequest
	decoder := json.NewDecoder(b)
	if err := decoder.Decode(&body); err != nil {
		return nil, errors.Wrap(err, "invalid put zLstate body")
	}
	if err := validatePutZLStateComponentRequest(&body); err != nil {
		return nil, errors.Wrap(err, "error validating put zLstate resource body")
	}

	log.WithField("body", body).Info("Handling put zLstate component status request")

	backend := zlstate.NewS3Backend(ctx, log, BuildZLStateBucketName(body.Company), s3Client)

	key := BuildZLStateKey(body.Team, body.Environment)

	zlst, err := backend.UpsertComponent(key, body.Component)
	if err != nil {
		return nil, errors.Wrapf(err, "error upserting component %s", body.Component.Name)
	}

	return &PutZLStateComponentResponse{ZLState: zlst}, nil
}

type PutZLStateComponentRequest struct {
	Company     string             `json:"company"`
	Team        string             `json:"team"`
	Environment string             `json:"environment"`
	Component   *zlstate.Component `json:"component"`
}

type PutZLStateComponentResponse struct {
	ZLState *zlstate.ZLState `json:"zlstate"`
}

func deleteZLStateComponentHandler(ctx context.Context, log *logrus.Entry, b io.ReadCloser, s3Client zlstate.S3API) (*DeleteZLStateComponentResponse, error) {
	var body DeleteZLStateComponentRequest
	decoder := json.NewDecoder(b)
	if err := decoder.Decode(&body); err != nil {
		return nil, errors.Wrap(err, "invalid delete zLstate body")
	}
	if err := validateDeleteZLStateComponentRequest(&body); err != nil {
		return nil, errors.Wrap(err, "error validating delete zLstate resource body")
	}

	log.WithField("body", body).Info("Handling delete zLstate component status request")

	backend := zlstate.NewS3Backend(ctx, log, BuildZLStateBucketName(body.Company), s3Client)

	key := BuildZLStateKey(body.Team, body.Environment)

	zlst, err := backend.DeleteComponent(key, body.Component)
	if err != nil {
		return nil, errors.Wrapf(err, "error upserting component %s", body.Component)
	}

	return &DeleteZLStateComponentResponse{ZLState: zlst}, nil
}

type DeleteZLStateComponentRequest struct {
	Company     string `json:"company"`
	Team        string `json:"team"`
	Environment string `json:"environment"`
	Component   string `json:"component"`
}

type DeleteZLStateComponentResponse struct {
	ZLState *zlstate.ZLState `json:"zlstate"`
}
