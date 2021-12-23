package controllers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/compuzest/zlifecycle-state-manager/app/apm"
	"github.com/compuzest/zlifecycle-state-manager/app/il"
	http2 "github.com/compuzest/zlifecycle-state-manager/app/web/http"
	"github.com/compuzest/zlifecycle-state-manager/app/zlog"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/pkg/errors"
)

func StateHandler(w http.ResponseWriter, r *http.Request) {
	txn := newrelic.FromContext(r.Context())

	var err error
	var resp interface{}
	var statusCode int
	switch r.Method {
	case "POST":
		resp, err = postStateHandler(r.Context(), r.Body)
		statusCode = http.StatusOK
	case "DELETE":
		resp, err = deleteStateResourcesHandler(r.Context(), r.Body)
		statusCode = http.StatusOK
	default:
		err := apm.NoticeError(
			txn,
			http2.NewVerboseError("NotFoundError", r.Method, "/state", errors.New("endpoint not implemented")),
		)
		zlog.CtxLogger(r.Context()).Error(err)
		http2.ErrorResponse(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		verr := http2.NewVerboseError("StateError", r.Method, "/state", err)
		_ = apm.NoticeError(txn, verr)
		zlog.CtxLogger(r.Context()).WithError(verr).Errorf("state handler error")
		zlog.CtxLogger(r.Context()).Errorf("%+v", verr.OriginalError)
		http2.ErrorResponse(w, verr.Error(), http.StatusBadRequest)
		return
	}

	http2.Response(w, resp, statusCode)
}

func postStateHandler(ctx context.Context, b io.ReadCloser) (*GetStateResponse, error) {
	var body GetStateRequest
	decoder := json.NewDecoder(b)
	if err := decoder.Decode(&body); err != nil {
		return nil, errors.Wrap(err, "invalid get state body")
	}
	if err := validateZState(body.ZState); err != nil {
		return nil, errors.Wrap(err, "error validating get state resource body")
	}

	s, err := il.FetchState(ctx, body.ZState)
	if err != nil {
		return nil, errors.Wrap(err, "error handling get state request")
	}

	resp := &GetStateResponse{State: s.GetRawState(), Resources: s.ParseResources()}
	return resp, nil
}

type GetStateRequest struct {
	ZState *il.ZState `json:"zstate"`
}

type GetStateResponse struct {
	State     *tfjson.State `json:"state"`
	Resources []string      `json:"resources"`
}

func deleteStateResourcesHandler(ctx context.Context, b io.ReadCloser) (*DeleteStateResourcesResponse, error) {
	var body DeleteStateResourcesRequest
	decoder := json.NewDecoder(b)
	if err := decoder.Decode(&body); err != nil {
		return nil, errors.Wrap(err, "invalid delete state resource request body")
	}
	if err := validateZState(body.ZState); err != nil {
		return nil, errors.Wrap(err, "error validating delete state resource body")
	}

	s, err := il.RemoveStateResources(ctx, body.ZState, body.Resources)
	if err != nil {
		return nil, errors.Wrap(err, "error handling delete state request")
	}

	resp := &DeleteStateResourcesResponse{State: s.GetRawState(), Resources: s.ParseResources()}
	return resp, nil
}

type DeleteStateResourcesRequest struct {
	ZState    *il.ZState `json:"zstate"`
	Resources []string   `json:"resources"`
}

type DeleteStateResourcesResponse struct {
	State     *tfjson.State `json:"state"`
	Resources []string      `json:"resources"`
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}
