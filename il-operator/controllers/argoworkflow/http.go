/* Copyright (C) 2020 CompuZest, Inc. - All Rights Reserved
 *
 * Unauthorized copying of this file, via any medium, is strictly prohibited
 * Proprietary and confidential
 *
 * NOTICE: All information contained herein is, and remains the property of
 * CompuZest, Inc. The intellectual and technical concepts contained herein are
 * proprietary to CompuZest, Inc. and are protected by trade secret or copyright
 * law. Dissemination of this information or reproduction of this material is
 * strictly forbidden unless prior written permission is obtained from CompuZest, Inc.
 */

package argoworkflow

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-errors/errors"

	"github.com/compuzest/zlifecycle-il-operator/controllers/util/common"
)

func NewHTTPClient(ctx context.Context, serverURL string) API {
	return &HTTPAPI{ctx: ctx, serverURL: serverURL}
}

func (api *HTTPAPI) ListWorkflows(opts ListWorkflowOptions) (*ListWorkflowsResponse, *http.Response, error) {
	listWorkflowsURL := fmt.Sprintf("%s/api/v1/workflows/%s?fields=items.metadata.name", api.serverURL, opts.Namespace)
	req, err := http.NewRequestWithContext(api.ctx, "GET", listWorkflowsURL, http.NoBody)
	if err != nil {
		return nil, nil, err
	}

	client := common.GetHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to send GET request to /api/v1/workflows/%s: %w", opts.Namespace, err)
	}

	respBody, err := common.ReadBody(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	r := new(ListWorkflowsResponse)
	if err := common.FromJSON(r, respBody); err != nil {
		common.CloseBody(resp.Body)
		return nil, nil, err
	}

	return r, resp, nil
}

func (api *HTTPAPI) DeleteWorkflow(opts DeleteWorkflowOptions) (*http.Response, error) {
	deleteWorkflowURL := fmt.Sprintf("%s/api/v1/workflows/%s/%s", api.serverURL, opts.Namespace, opts.Name)
	req, err := http.NewRequestWithContext(api.ctx, "DELETE", deleteWorkflowURL, http.NoBody)
	if err != nil {
		return nil, err
	}

	client := common.GetHTTPClient()

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send DELETE request to /api/v1/workflows/%s/%s: %w", opts.Namespace, opts.Name, err)
	}
	defer common.CloseBody(resp.Body)

	if resp.StatusCode == 404 {
		return resp, nil
	}

	if resp.StatusCode != 200 {
		return nil, errors.Errorf("delete workflow returned non-OK response: %d", resp.StatusCode)
	}

	return resp, nil
}
