package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

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
			http2.NewVerboseError("InternalError", r.Method, "/zl/state", errors.New("internal server error")),
		)
		zlog.CtxLogger(r.Context()).Error(err)
		http2.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	switch r.Method {
	case "POST":
		resp, err = postZLStateComponentHandler(r.Context(), zlog.CtxLogger(r.Context()), r.Body, s3Client)
		statusCode = http.StatusOK
	case "PATCH":
		resp, err = patchZLStateComponentHandler(r.Context(), zlog.CtxLogger(r.Context()), r.Body, s3Client)
		statusCode = http.StatusOK
	default:
		err := apm.NoticeError(
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

func postZLStateComponentHandler(ctx context.Context, log *logrus.Entry, b io.ReadCloser, s3Client zlstate.S3API) (*FetchZLStateComponentResponse, error) {
	var body FetchZLStateComponentRequest
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

	return &FetchZLStateComponentResponse{Component: component}, nil
}

func findComponent(components []*zlstate.Component, targetComponent string) *zlstate.Component {
	for _, c := range components {
		if c.Name == targetComponent {
			return c
		}
	}
	return nil
}

type FetchZLStateComponentRequest struct {
	Company     string `json:"company"`
	Team        string `json:"team"`
	Environment string `json:"environment"`
	Component   string `json:"component"`
}

type FetchZLStateComponentResponse struct {
	Component *zlstate.Component `json:"component"`
}

func patchZLStateComponentHandler(ctx context.Context, log *logrus.Entry, b io.ReadCloser, s3Client zlstate.S3API) (*UpdateZLStateComponentResponse, error) {
	var body UpdateZLStateComponentRequest
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

	zlState, err := backend.Get(key)
	if err != nil {
		return nil, errors.Wrap(err, "error getting zLstate from remote backend")
	}

	var oldStatus string
	updated := false
	for _, c := range zlState.Components {
		if c.Name != body.Component {
			continue
		}
		oldStatus = c.Status
		c.Status = body.Status
		c.UpdatedAt = time.Now().UTC()
		zlState.UpdatedAt = time.Now().UTC()
		updated = true
		break
	}
	if !updated {
		return nil, errors.Errorf("component not found: %s", body.Component)
	}

	if err := backend.Put(key, zlState, true); err != nil {
		return nil, errors.Wrap(err, "error persisting zLstate to remote backend")
	}

	msg := fmt.Sprintf("updated environment component [%s] status from [%s] to [%s]", body.Component, oldStatus, body.Status)
	log.WithFields(logrus.Fields{
		"company":     body.Company,
		"team":        body.Team,
		"environment": body.Environment,
		"component":   body.Component,
	}).Info(msg)
	return &UpdateZLStateComponentResponse{
		Message: fmt.Sprintf("updated environment component [%s] status from [%s] to [%s]", body.Component, oldStatus, body.Status),
	}, nil
}

type UpdateZLStateComponentRequest struct {
	Company     string `json:"company"`
	Team        string `json:"team"`
	Environment string `json:"environment"`
	Component   string `json:"component"`
	Status      string `json:"status"`
}

type UpdateZLStateComponentResponse struct {
	Message string `json:"message"`
}
