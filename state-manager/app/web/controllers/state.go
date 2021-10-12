package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/compuzest/zlifecycle-state-manager/app/il"
	http2 "github.com/compuzest/zlifecycle-state-manager/app/web/http"
	"github.com/compuzest/zlifecycle-state-manager/app/zlog"
	tfjson "github.com/hashicorp/terraform-json"
	"io"
	"net/http"
)

var ctx = context.Background()

func StateHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var resp interface{}
	var statusCode int
	switch r.Method {
	case "POST":
		resp, err = PostStateHandler(r.Body)
		statusCode = http.StatusOK
	case "DELETE":
		resp, err = DeleteStateResourcesHandler(r.Body)
		statusCode = http.StatusOK
	default:
		err := fmt.Errorf("endpoint not implemented")
		zlog.Logger.Error(err)
		http2.ErrorResponse(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		zlog.Logger.Error(err)
		http2.ErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	http2.Response(w, resp, statusCode)
}

func PostStateHandler(b io.ReadCloser) (*GetStateResponse, error) {
	var body GetStateRequest
	decoder := json.NewDecoder(b)
	if err := decoder.Decode(&body); err != nil {
		return nil, err
	}

	s, err := il.FetchState(ctx, body.ZState)
	if err != nil {
		return nil, err
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

func DeleteStateResourcesHandler(b io.ReadCloser) (*DeleteStateResourcesResponse, error) {
	var body DeleteStateResourcesRequest
	decoder := json.NewDecoder(b)
	if err := decoder.Decode(&body); err != nil {
		return nil, err
	}

	s, err := il.RemoveStateResources(ctx, body.ZState, body.Resources)

	if err != nil {
		return nil, err
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
