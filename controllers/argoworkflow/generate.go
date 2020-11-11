package controllers

import (
	workflow "github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	stablev1alpha1 "github.com/compuzest/environment-operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GenerateWorkflowOfWorkflows(environment stablev1alpha1.Environment) *workflow.Workflow {

	workflowTemplate := "terraform-sync-template"
	customerId := string(environment.Spec.CustomerId)

	var parallelSteps []workflow.ParallelSteps

	for _, terraformConfig := range environment.Spec.TerraformConfigs {
		step := workflow.ParallelSteps{
			Steps: []workflow.WorkflowStep{
				{
					Name: terraformConfig.Name,
					TemplateRef: &workflow.TemplateRef{
						Name:     "workflow-trigger-template",
						Template: "run",
					},
					Arguments: workflow.Arguments{
						Parameters: []workflow.Parameter{
							{
								Name:  "workflowtemplate",
								Value: &workflowTemplate,
							},
							{
								Name:  "module_source",
								Value: &terraformConfig.Module.Source,
							},
							{
								Name:  "module_source_path",
								Value: &terraformConfig.Module.Path,
							},
							{
								Name:  "env_name",
								Value: &environment.Spec.Name,
							},
							{
								Name:  "customer_id",
								Value: &customerId,
							},
							{
								Name:  "name",
								Value: &terraformConfig.Name,
							},
						},
					},
				},
			},
		}

		parallelSteps = append(parallelSteps, step)
	}

	return &workflow.Workflow{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "argoproj.io/v1alpha1",
			Kind:       "Workflow",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      environment.Spec.CustomerId + "-" + environment.Spec.Name,
			Namespace: "argo",
		},
		Spec: workflow.WorkflowSpec{
			Entrypoint: "main",
			Templates: []workflow.Template{
				{
					Name:  "main",
					Steps: parallelSteps,
				},
			},
		},
	}
}
