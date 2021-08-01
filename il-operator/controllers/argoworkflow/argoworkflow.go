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
	workflow "github.com/argoproj/argo/v2/pkg/apis/workflow/v1alpha1"
	stablev1 "github.com/compuzest/zlifecycle-il-operator/api/v1"
	terraformgenerator "github.com/compuzest/zlifecycle-il-operator/controllers/terraformgenerator"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/common"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/env"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/il"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GenerateWorkflowOfWorkflows create WoW
func GenerateWorkflowOfWorkflows(environment stablev1.Environment) *workflow.Workflow {
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
			Dependencies: append(environmentComponent.DependsOn, "trigger-audit"),
			Arguments: workflow.Arguments{
				Parameters: []workflow.Parameter{
					{
						Name:  "workflowtemplate",
						Value: anyStringPointer(workflowTemplate),
					},
					{
						Name:  "terraform_version",
						Value: anyStringPointer(terraformgenerator.DefaultTerraformVersion),
					},
					{
						Name:  "terraform_il_path",
						Value: anyStringPointer(tfPath),
					},
					{
						Name:  "il_repo",
						Value: anyStringPointer(env.Config.ILRepoURL),
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
			Labels: map[string]string{
				"terraform/sync":       "true",
				"zlifecycle.com/model": "environment-sync-flow",
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

func GenerateLegacyWorkflowOfWorkflows(environment stablev1.Environment) *workflow.Workflow {
	workflowTemplate := "terraform-sync-template"
	envComponentDirectory := il.EnvironmentComponentDirectory(environment.Spec.TeamName, environment.Spec.EnvName)
	tf := terraformgenerator.TerraformGenerator{}

	var tasks []workflow.DAGTask

	destroyAll := !environment.DeletionTimestamp.IsZero()

	tasks = append(tasks, generateAuditTask(environment, destroyAll, "0", nil))

	ecs := environment.Spec.EnvironmentComponent

	var allComponents []string

	for _, ec := range ecs {
		allComponents = append(allComponents, ec.Name)

		moduleSource := il.EnvComponentModuleSource(ec.Module.Source, ec.Module.Name)
		modulePath := il.EnvComponentModulePath(ec.Module.Path)
		tfPath := tf.GenerateTerraformIlPath(envComponentDirectory, ec.Name)

		destroyFlag := false
		if ec.MarkedForDeletion {
			destroyFlag = true
		}
		var dependencies []string
		if destroyAll {
			dependencies = buildInverseDependencies(ecs, ec.Name)
			destroyFlag = true
		} else {
			dependencies = ec.DependsOn
		}

		dependencies = append(dependencies, "trigger-audit")

		parameters := []workflow.Parameter{
			{
				Name:  "workflowtemplate",
				Value: anyStringPointer(workflowTemplate),
			},
			{
				Name:  "module_source",
				Value: anyStringPointer(moduleSource),
			},
			{
				Name:  "module_source_path",
				Value: anyStringPointer(modulePath),
			},
			{
				Name:  "team_name",
				Value: anyStringPointer(environment.Spec.TeamName),
			},
			{
				Name:  "env_name",
				Value: anyStringPointer(environment.Spec.EnvName),
			},
			{
				Name:  "config_name",
				Value: anyStringPointer(ec.Name),
			},
			{
				Name:  "il_repo",
				Value: anyStringPointer(env.Config.ILRepoURL),
			},
			{
				Name:  "terraform_il_path",
				Value: anyStringPointer(tfPath),
			},
			{
				Name:  "is_destroy",
				Value: anyStringPointer(destroyFlag),
			},
			{
				Name:  "reconcile_id",
				Value: anyStringPointer("{{tasks.trigger-audit.outputs.parameters.reconcile_id}}"),
			},
			{
				Name:  "status",
				Value: anyStringPointer("initializing"),
			},
		}

		task := workflow.DAGTask{
			Name: ec.Name,
			TemplateRef: &workflow.TemplateRef{
				Name:     "workflow-trigger-template",
				Template: "run",
			},
			Dependencies: dependencies,
			Arguments: workflow.Arguments{
				Parameters: parameters,
			},
		}

		tasks = append(tasks, task)
	}

	tasks = append(tasks, generateAuditTask(environment, destroyAll, "1", allComponents))

	return &workflow.Workflow{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "argoproj.io/v1alpha1",
			Kind:       "Workflow",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      environment.Spec.TeamName + "-" + environment.Spec.EnvName,
			Namespace: "argocd",
			Labels: map[string]string{
				"terraform/sync":       "true",
				"zlifecycle.com/model": "environment-sync-flow",
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

func generateAuditTask(environment stablev1.Environment, destroyAll bool, phase string, dependencies []string) workflow.DAGTask {
	var name string
	var reconcileId string
	var status string

	if phase == "0" {
		name = "trigger-audit"
		status = "initializing"
		reconcileId = "0"
	} else {
		name = "end-audit"
		status = "ended"
		reconcileId = "{{tasks.trigger-audit.outputs.parameters.reconcile_id}}"
	}

	task := workflow.DAGTask{
		Name: name,
		TemplateRef: &workflow.TemplateRef{
			Name:     "audit-run-template",
			Template: "run_audit",
		},
		Dependencies: dependencies,
		Arguments: workflow.Arguments{
			Parameters: []workflow.Parameter{
				{
					Name:  "team_name",
					Value: anyStringPointer(environment.Spec.TeamName),
				},
				{
					Name:  "env_name",
					Value: anyStringPointer(environment.Spec.EnvName),
				},
				{
					Name:  "status",
					Value: anyStringPointer(status),
				},
				{
					Name:  "is_destroy",
					Value: anyStringPointer(destroyAll),
				},
				{
					Name:  "phase",
					Value: anyStringPointer(phase),
				},
				{
					Name:  "reconcile_id",
					Value: anyStringPointer(reconcileId),
				},
			},
		},
	}

	return task
}

func anyStringPointer(val interface{}) *workflow.AnyString {
	s := workflow.ParseAnyString(val)
	return &s
}

func buildInverseDependencies(components []*stablev1.EnvironmentComponent, component string) []string {
	var dependencies []string
	for _, c := range components {
		if component == c.Name {
			continue
		} else if common.Contains(c.DependsOn, component) {
			dependencies = append(dependencies, c.Name)
		}
	}

	return dependencies
}