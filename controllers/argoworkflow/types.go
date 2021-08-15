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
	"github.com/go-logr/logr"
	"net/http"
)

type Api interface {
	ListWorkflows(opts ListWorkflowOptions) (*ListWorkflowsResponse, *http.Response, error)
	DeleteWorkflow(opts DeleteWorkflowOptions) (*http.Response, error)
}

type HttpApi struct {
	ServerUrl string
	Log       logr.Logger
}

type ListWorkflowOptions struct {
	Namespace string `json:"namespace"`
}

type DeleteWorkflowOptions struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

type ListWorkflowsResponse struct {
	Items []struct {
		Metadata struct {
			Name string `json:"name"`
		} `json:"metadata"`
	} `json:"items"`
}
