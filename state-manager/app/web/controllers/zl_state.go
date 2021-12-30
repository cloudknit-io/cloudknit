package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/compuzest/zlifecycle-state-manager/app/zlstate"
	"io"
	"net/http"

	"github.com/compuzest/zlifecycle-state-manager/app/apm"
	http2 "github.com/compuzest/zlifecycle-state-manager/app/web/http"
	"github.com/compuzest/zlifecycle-state-manager/app/zlog"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/pkg/errors"
)

func ZLStateHandler(w http.ResponseWriter, r *http.Request) {
	txn := newrelic.FromContext(r.Context())

	var err error
	var resp interface{}
	var statusCode int
	switch r.Method {
	case "POST":
		resp, err = postZLStateHandler(r.Body)
		statusCode = http.StatusOK
	case "PUT":
		resp, err = putZLStateHandler(r.Body)
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

func postZLStateHandler(b io.ReadCloser) (*GetZLStateResponse, error) {
	var body GetZLStateRequest
	decoder := json.NewDecoder(b)
	if err := decoder.Decode(&body); err != nil {
		return nil, errors.Wrap(err, "invalid get zLstate body")
	}
	if err := validateGetZLStateRequest(&body); err != nil {
		return nil, errors.Wrap(err, "error validating get zLstate resource body")
	}

	client := zlstate.NewS3Backend(BuildZLStateBucket(body.Company))

	zlstate, err := client.Get(BuildZLStateKey(body.Team, body.Environment))
	if err != nil {
		return nil, errors.Wrap(err, "error getting zLstate from remote backend")
	}

	resp := &GetZLStateResponse{ZLState: zlstate}
	return resp, nil
}

type GetZLStateRequest struct {
	Company     string `json:"company"`
	Team        string `json:"team"`
	Environment string `json:"environment"`
}

type GetZLStateResponse struct {
	ZLState *zlstate.ZLState `json:"zlstate"`
}

func putZLStateHandler(b io.ReadCloser) (*PostZLStateResponse, error) {
	var body PostZLStateRequest
	decoder := json.NewDecoder(b)
	if err := decoder.Decode(&body); err != nil {
		return nil, errors.Wrap(err, "invalid get zLstate body")
	}
	if err := validatePostZLStateRequest(&body); err != nil {
		return nil, errors.Wrap(err, "error validating post zLstate resource body")
	}

	client := zlstate.NewS3Backend(BuildZLStateBucket(body.Company))

	if err := client.Put(BuildZLStateKey(body.Team, body.Environment), body.ZLState); err != nil {
		return nil, errors.Wrap(err, "error persisting zLstate to remote backend")
	}

	resp := &PostZLStateResponse{}
	return resp, nil
}

type PostZLStateRequest struct {
	Company     string           `json:"company"`
	Team        string           `json:"team"`
	Environment string           `json:"environment"`
	ZLState     *zlstate.ZLState `json:"zlstate"`
}

type PostZLStateResponse struct{}

func validateGetZLStateRequest(req *GetZLStateRequest) error {
	if req.Company == "" {
		return errors.New(`request body is missing field "company"`)
	}
	if req.Team == "" {
		return errors.New(`request body is missing field "team"`)
	}
	if req.Environment == "" {
		return errors.New(`request body is missing field "environment"`)
	}

	return nil
}

func validatePostZLStateRequest(req *PostZLStateRequest) error {
	if req.Company == "" {
		return errors.New(`request body is missing field "company"`)
	}
	if req.Team == "" {
		return errors.New(`request body is missing field "team"`)
	}
	if req.Environment == "" {
		return errors.New(`request body is missing field "environment"`)
	}
	if req.ZLState == nil {
		return errors.New(`request body is missing field "zlstate"`)
	}

	return nil
}

func BuildZLStateKey(team, environment string) string {
	return fmt.Sprintf("%s/%s.zlstate", team, environment)
}

func BuildZLStateBucket(company string) string {
	return fmt.Sprintf("zlifecycle-zlstate-%s", company)
}
