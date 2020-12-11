package controllers

import (
	workflow "github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	stablev1alpha1 "github.com/compuzest/environment-operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GenerateWorkflowOfWorkflows(environment stablev1alpha1.Environment) *workflow.Workflow {

	workflowTemplate := "terraform-sync-template"

	var tasks []workflow.DAGTask

	for _, terraformConfig := range environment.Spec.TerraformConfigs {

		task := workflow.DAGTask{
			Name: terraformConfig.ConfigName,
			TemplateRef: &workflow.TemplateRef{
				Name:     "workflow-trigger-template",
				Template: "run",
			},
			Dependencies: terraformConfig.DependsOn,
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
						Name:  "variables_file_source",
						Value: &terraformConfig.VariablesFile.Source,
					},
					{
						Name:  "variables_file_path",
						Value: &terraformConfig.VariablesFile.Path,
					},
					{
						Name:  "team_name",
						Value: &environment.Spec.TeamName,
					},
					{
						Name:  "env_name",
						Value: &environment.Spec.EnvName,
					},
					{
						Name:  "config_name",
						Value: &terraformConfig.ConfigName,
					},
				},
			},
		}

		tasks = append(tasks, task)
	}

	return &workflow.Workflow{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "argoproj.io/v1alpha1",
			Kind:       "Workflow",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      environment.Spec.TeamName + "-" + environment.Spec.EnvName,
			Namespace: "argocd",
			Annotations: map[string]string{
				"argocd.argoproj.io/hook":               "PreSync",
				"argocd.argoproj.io/hook-delete-policy": "BeforeHookCreation",
			},
			Labels: map[string]string{
				"workflows.argoproj.io/completed": "false",
				"terraform/sync":                  "true",
			},
		},
		Spec: workflow.WorkflowSpec{
			Entrypoint: "main",
			Templates: []workflow.Template{
				{
					Name: "main",
					DAG: &workflow.DAGTemplate{
						Tasks: tasks,
					},
				},
			},
		},
	}
}
