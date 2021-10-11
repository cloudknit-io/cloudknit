package controllers

import (
	"context"
	"encoding/json"
	"github.com/compuzest/zlifecycle-state-manager/app/il"
	"github.com/compuzest/zlifecycle-state-manager/app/zlog"
	tfjson "github.com/hashicorp/terraform-json"
	"io"
	"net/http"
)

var ctx = context.Background()

func StateHandler(w http.ResponseWriter, r *http.Request) {
	var err        error
	var resp       interface{}
	var statusCode int
	switch r.Method {
	case "GET":
 		resp, err  = GetStateHandler(r.Body)
 		statusCode = http.StatusOK
	case "DELETE":
		resp, err  = DeleteStateResourcesHandler(r.Body)
		statusCode = http.StatusNoContent
	}
	if err != nil {
		zlog.Logger.Error(err)
		ErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	Response(w, resp, statusCode)
}

func GetStateHandler(b io.ReadCloser) (*GetStateResponse, error) {
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
