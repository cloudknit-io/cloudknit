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
