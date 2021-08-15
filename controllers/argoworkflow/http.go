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
	"fmt"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/common"
	"github.com/go-logr/logr"
	"net/http"
)

func NewHttpClient(l logr.Logger, serverUrl string) Api {
	return HttpApi{Log: l, ServerUrl: serverUrl}
}

func (api HttpApi) ListWorkflows(opts ListWorkflowOptions) (*ListWorkflowsResponse, *http.Response, error) {
	listWorkflowsUrl := fmt.Sprintf("%s/api/v1/workflows/%s?fields=items.metadata.name", api.ServerUrl, opts.Namespace)
	resp, err := http.Get(listWorkflowsUrl)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to send GET request to /api/v1/workflows/%s: %v", opts.Namespace, err)
	}

	respBody, err := common.ReadBody(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	r := new(ListWorkflowsResponse)
	if err := common.FromJson(r, respBody); err != nil {
		common.CloseBody(resp.Body)
		return nil, nil, err
	}

	return r, resp, nil
}

func (api HttpApi) DeleteWorkflow(opts DeleteWorkflowOptions) (*http.Response, error) {
	deleteWorkflowUrl := fmt.Sprintf("%s/api/v1/workflows/%s/%s", api.ServerUrl, opts.Namespace, opts.Name)
	req, err := http.NewRequest("DELETE", deleteWorkflowUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to send DELETE request to /api/v1/workflows/%s/%s: %v", opts.Namespace, opts.Name, err)
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer common.CloseBody(resp.Body)

	if resp.StatusCode == 404 {
		api.Log.Info("Argo Workflow does not exist", "namespace", opts.Namespace, "workflow", opts.Name)
		return resp, nil
	}

	if resp.StatusCode != 200 {
		common.LogBody(api.Log, resp.Body)
		return nil, fmt.Errorf("delete workflow returned non-OK response: %d", resp.StatusCode)
	}

	return resp, nil
}
