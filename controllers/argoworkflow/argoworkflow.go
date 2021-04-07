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
	workflow "github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	stablev1alpha1 "github.com/compuzest/zlifecycle-il-operator/api/v1alpha1"
	terraformgenerator "github.com/compuzest/zlifecycle-il-operator/controllers/terraformgenerator"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/env"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/il"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GenerateWorkflowOfWorkflows create WoW
func GenerateWorkflowOfWorkflows(environment stablev1alpha1.Environment) *workflow.Workflow {
	envComponentDirectory := il.EnvironmentComponentDirectory(environment.Spec.TeamName, environment.Spec.EnvName)
	workflowTemplate := "terraform-provision-template"

	var tasks []workflow.DAGTask

	for _, environmentComponent := range environment.Spec.EnvironmentComponent {
		tf := terraformgenerator.TerraformGenerator{}
		tfPath := tf.GenerateTerraformIlPath(envComponentDirectory, environmentComponent.Name)
		task := workflow.DAGTask{
			Name: environmentComponent.Name,
			TemplateRef: &workflow.TemplateRef{
				Name:     "provisioner-trigger-template",
				Template: "run",
			},
			Dependencies: environmentComponent.DependsOn,
			Arguments: workflow.Arguments{
				Parameters: []workflow.Parameter{
					{
						Name:  "workflowtemplate",
						Value: &workflowTemplate,
					},
					{
						Name:  "terraform_version",
						Value: &terraformgenerator.DefaultTerraformVersion,
					},
					{
						Name:  "terraform_il_path",
						Value: &tfPath,
					},
					{
						Name:  "il_repo",
						Value: &env.Config.ILRepoURL,
						// to be replaced with reference to il.RepoURL(owner, companyName) once company can be extrapolated here
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
			Name:      "experimental-" + environment.Spec.TeamName + "-" + environment.Spec.EnvName,
			Namespace: "argocd",
			Annotations: map[string]string{
				"argocd.argoproj.io/hook":               "PreSync",
				"argocd.argoproj.io/hook-delete-policy": "BeforeHookCreation",
			},
			Labels: map[string]string{
				"workflows.argoproj.io/completed": "false",
				"terraform/sync":                  "true",
				"zlifecycle.com/model":            "environment-sync-flow",
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

func GenerateLegacyWorkflowOfWorkflows(environment stablev1alpha1.Environment) *workflow.Workflow {

	workflowTemplate := "terraform-sync-template"

	var tasks []workflow.DAGTask

	for _, environmentComponent := range environment.Spec.EnvironmentComponent {
		moduleSource := il.EnvComponentModuleSource(environmentComponent.Module.Source, environmentComponent.Module.Name)
		modulePath := il.EnvComponentModulePath(environmentComponent.Module.Path)

		task := workflow.DAGTask{
			Name: environmentComponent.Name,
			TemplateRef: &workflow.TemplateRef{
				Name:     "workflow-trigger-template",
				Template: "run",
			},
			Dependencies: environmentComponent.DependsOn,
			Arguments: workflow.Arguments{
				Parameters: []workflow.Parameter{
					{
						Name:  "workflowtemplate",
						Value: &workflowTemplate,
					},
					{
						Name:  "module_source",
						Value: &moduleSource,
					},
					{
						Name:  "module_source_path",
						Value: &modulePath,
					},
					{
						Name:  "variables_file_source",
						Value: &environmentComponent.VariablesFile.Source,
					},
					{
						Name:  "variables_file_path",
						Value: &environmentComponent.VariablesFile.Path,
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
						Value: &environmentComponent.Name,
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
				"zlifecycle.com/model":            "environment-sync-flow",
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
