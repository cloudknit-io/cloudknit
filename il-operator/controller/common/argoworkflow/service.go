package argoworkflow

import "net/http"

type API interface {
	ListWorkflows(opts ListWorkflowOptions) (*ListWorkflowsResponse, *http.Response, error)
	DeleteWorkflow(opts DeleteWorkflowOptions) (*http.Response, error)
}
