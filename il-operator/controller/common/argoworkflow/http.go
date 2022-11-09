package argoworkflow

import (
	"context"
	"fmt"
	"net/http"

	"github.com/compuzest/zlifecycle-il-operator/controller/util"

	"github.com/go-errors/errors"
)

type HTTPAPI struct {
	ctx       context.Context
	serverURL string
}

func NewHTTPClient(ctx context.Context, serverURL string) API {
	return &HTTPAPI{ctx: ctx, serverURL: serverURL}
}

func (api *HTTPAPI) ListWorkflows(opts ListWorkflowOptions) (*ListWorkflowsResponse, *http.Response, error) {
	listWorkflowsURL := fmt.Sprintf("%s/api/v1/workflows/%s?fields=items.metadata.name", api.serverURL, opts.Namespace)
	req, err := http.NewRequestWithContext(api.ctx, http.MethodGet, listWorkflowsURL, http.NoBody)
	if err != nil {
		return nil, nil, err
	}

	client := util.GetHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to send GET request to /api/v1/workflows/%s: %w", opts.Namespace, err)
	}

	respBody, err := util.ReadBody(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	r := new(ListWorkflowsResponse)
	if err := util.FromJSON(r, respBody); err != nil {
		util.CloseBody(resp.Body)
		return nil, nil, err
	}

	return r, resp, nil
}

func (api *HTTPAPI) DeleteWorkflow(opts DeleteWorkflowOptions) (*http.Response, error) {
	deleteWorkflowURL := fmt.Sprintf("%s/api/v1/workflows/%s/%s", api.serverURL, opts.Namespace, opts.Name)
	req, err := http.NewRequestWithContext(api.ctx, http.MethodDelete, deleteWorkflowURL, http.NoBody)
	if err != nil {
		return nil, err
	}

	client := util.GetHTTPClient()

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send DELETE request to /api/v1/workflows/%s/%s: %w", opts.Namespace, opts.Name, err)
	}
	defer util.CloseBody(resp.Body)

	if resp.StatusCode == 404 {
		return resp, nil
	}

	if resp.StatusCode != 200 {
		return nil, errors.Errorf("delete workflow returned non-OK response: %d", resp.StatusCode)
	}

	return resp, nil
}
