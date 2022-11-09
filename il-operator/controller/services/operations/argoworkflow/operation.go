package argoworkflow

import (
	"strings"

	argoworkflowapi "github.com/compuzest/zlifecycle-il-operator/controller/common/argoworkflow"

	"github.com/compuzest/zlifecycle-il-operator/controller/util"

	"github.com/go-logr/logr"
)

func DeleteWorkflowsWithPrefix(log logr.Logger, prefix string, namespace string, api argoworkflowapi.API) error {
	log.Info("Listing Argo Workflows", "prefix", prefix)
	listOpts := argoworkflowapi.ListWorkflowOptions{Namespace: namespace}
	wfs, listResp, err := api.ListWorkflows(listOpts)
	if err != nil {
		return err
	}
	defer util.CloseBody(listResp.Body)

	for _, wf := range wfs.Items {
		name := wf.Metadata.Name
		if strings.HasPrefix(name, prefix) {
			log.Info("Deleting Argo Workflow", "namespace", namespace, "name", name)
			if err := DeleteWorkflow(name, namespace, api); err != nil {
				return err
			}
		}
	}

	return nil
}

func DeleteWorkflow(name string, namespace string, api argoworkflowapi.API) error {
	deleteOpts := argoworkflowapi.DeleteWorkflowOptions{Name: name, Namespace: namespace}
	resp, err := api.DeleteWorkflow(deleteOpts)
	if err != nil {
		return err
	}
	util.CloseBody(resp.Body)

	return nil
}
