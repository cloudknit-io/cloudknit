package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/compuzest/zlifecycle-state-manager/app/apm"
	http2 "github.com/compuzest/zlifecycle-state-manager/app/web/http"
	"github.com/compuzest/zlifecycle-state-manager/app/zlog"
	"github.com/compuzest/zlifecycle-state-manager/app/zlstate"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"time"
)

func ZLStateHandler(w http.ResponseWriter, r *http.Request) {
	txn := newrelic.FromContext(r.Context())

	var err error
	var resp interface{}
	var statusCode int
	switch r.Method {
	case "POST":
		resp, err = postZLStateHandler(r.Context(), r.Body)
		statusCode = http.StatusOK
	case "PUT":
		resp, err = putZLStateHandler(r.Context(), r.Body)
		statusCode = http.StatusOK
	case "PATCH":
		resp, err = patchZLStateHandler(r.Context(), r.Body)
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
		verr := http2.NewVerboseError("StateError", r.Method, "/zl/state", err)
		_ = apm.NoticeError(txn, verr)
		zlog.CtxLogger(r.Context()).WithError(verr).Errorf("zl state handler error")
		zlog.CtxLogger(r.Context()).Errorf("%+v", verr.OriginalError)
		http2.ErrorResponse(w, verr.Error(), http.StatusBadRequest)
		return
	}

	http2.Response(w, resp, statusCode)
}

func postZLStateHandler(ctx context.Context, b io.ReadCloser) (*GetZLStateResponse, error) {
	var body GetZLStateRequest
	decoder := json.NewDecoder(b)
	if err := decoder.Decode(&body); err != nil {
		return nil, errors.Wrap(err, "invalid get zLstate body")
	}
	if err := validateGetZLStateRequest(&body); err != nil {
		return nil, errors.Wrap(err, "error validating get zLstate resource body")
	}

	client, err := zlstate.NewS3Backend(ctx, BuildZLStateBucket(body.Company))
	if err != nil {
		return nil, errors.Wrap(err, "error instantiating s3 backend for zLstate manager")
	}

	zlState, err := client.Get(BuildZLStateKey(body.Team, body.Environment))
	if err != nil {
		return nil, errors.Wrap(err, "error getting zLstate from remote backend")
	}

	return &GetZLStateResponse{ZLState: zlState}, nil
}

type GetZLStateRequest struct {
	Company     string `json:"company"`
	Team        string `json:"team"`
	Environment string `json:"environment"`
}

type GetZLStateResponse struct {
	ZLState *zlstate.ZLState `json:"zlstate"`
}

func putZLStateHandler(ctx context.Context, b io.ReadCloser) (*PutZLStateResponse, error) {
	var body PutZLStateRequest
	decoder := json.NewDecoder(b)
	if err := decoder.Decode(&body); err != nil {
		return nil, errors.Wrap(err, "invalid put zLstate body")
	}
	if err := validatePutZLStateRequest(&body); err != nil {
		return nil, errors.Wrap(err, "error validating put zLstate resource body")
	}

	client, err := zlstate.NewS3Backend(ctx, BuildZLStateBucket(body.Company))
	if err != nil {
		return nil, errors.Wrap(err, "error instantiating s3 backend for zLstate manager")
	}

	if err := client.Put(BuildZLStateKey(body.Team, body.Environment), body.ZLState); err != nil {
		if errors.Is(err, zlstate.ErrKeyExists) {
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

func patchZLStateHandler(ctx context.Context, b io.ReadCloser) (*PatchZLStateResponse, error) {
	var body PatchZLStateRequest
	decoder := json.NewDecoder(b)
	if err := decoder.Decode(&body); err != nil {
		return nil, errors.Wrap(err, "invalid patch zLstate body")
	}
	if err := validatePatchZLStateRequest(&body); err != nil {
		return nil, errors.Wrap(err, "error validating patch zLstate resource body")
	}

	client, err := zlstate.NewS3Backend(ctx, BuildZLStateBucket(body.Company))
	if err != nil {
		return nil, errors.Wrap(err, "error instantiating s3 backend for zLstate manager")
	}

	key := BuildZLStateKey(body.Team, body.Environment)

	zlState, err := client.Get(key)
	if err != nil {
		return nil, errors.Wrap(err, "error getting zLstate from remote backend")
	}

	updated := false
	for _, c := range zlState.Components {
		if c.Name != body.Component {
			continue
		}
		c.Status = body.Status
		c.UpdatedAt = time.Now().UTC()
		zlState.UpdatedAt = time.Now().UTC()
		updated = true
		break
	}
	if !updated {
		return nil, errors.Errorf("component not found: %s", body.Component)
	}

	if err := client.Put(key, zlState); err != nil {
		return nil, errors.Wrap(err, "error persisting zLstate to remote backend")
	}

	return &PatchZLStateResponse{
		ZLState: zlState,
	}, nil
}

type PatchZLStateRequest struct {
	Company     string `json:"company"`
	Team        string `json:"team"`
	Environment string `json:"environment"`
	Component   string `json:"component"`
	Status      string `json:"statu"`
}

type PatchZLStateResponse struct {
	ZLState *zlstate.ZLState `json:"zlstate"`
}

func BuildZLStateKey(team, environment string) string {
	return fmt.Sprintf("%s/%s.zlstate", team, environment)
}

func BuildZLStateBucket(company string) string {
	return fmt.Sprintf("zlifecycle-zlstate-%s", company)
}
