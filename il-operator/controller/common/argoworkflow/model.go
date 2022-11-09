package argoworkflow

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
