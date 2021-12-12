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
	"time"

	workflow "github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	stablev1 "github.com/compuzest/zlifecycle-il-operator/api/v1"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/common"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/env"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/il"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GenerateWorkflowOfWorkflows create WoW.
func GenerateWorkflowOfWorkflows(environment *stablev1.Environment) *workflow.Workflow {
	workflowTemplate := "terraform-provision-template"

	var tasks []workflow.DAGTask

	for _, environmentComponent := range environment.Spec.Components {
		tfPath := il.EnvironmentComponentTerraformDirectoryPath(environment.Spec.TeamName, environment.Spec.EnvName, environmentComponent.Name)
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
						Value: AnyStringPointer(workflowTemplate),
					},
					{
						Name:  "terraform_version",
						Value: AnyStringPointer(env.Config.DefaultTerraformVersion),
					},
					{
						Name:  "terraform_il_path",
						Value: AnyStringPointer(tfPath),
					},
					{
						// TODO: to be replaced with reference to il.RepoURL(owner, companyName) once company can be extrapolated here
						Name:  "il_repo",
						Value: AnyStringPointer(env.Config.TFILRepoURL),
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
			OnExit:     "exit-handler",
			Templates: []workflow.Template{
				{
					Name: "main",
					DAG: &workflow.DAGTemplate{
						Tasks: tasks,
					},
				},
				exitHandler(environment),
			},
		},
	}
}

func GenerateLegacyWorkflowOfWorkflows(environment *stablev1.Environment) *workflow.Workflow {
	workflowTemplate := "terraform-sync-template"

	var tasks []workflow.DAGTask

	autoApproveAll := environment.Spec.AutoApprove
	destroyAll := !environment.DeletionTimestamp.IsZero() || environment.Spec.Teardown

	tasks = append(tasks, generateAuditTask(environment, destroyAll, "0", nil))

	ecs := environment.Spec.Components

	allComponents := make([]string, 0, len(ecs))

	for _, ec := range ecs {
		allComponents = append(allComponents, ec.Name)

		tfPath := il.EnvironmentComponentTerraformDirectoryPath(environment.Spec.TeamName, environment.Spec.EnvName, ec.Name)

		autoApproveFlag := ec.AutoApprove
		if autoApproveAll {
			autoApproveFlag = true
		}

		dependencies := ec.DependsOn
		destroyFlag := ec.Destroy
		if destroyAll {
			destroyFlag = true
			dependencies = buildInverseDependencies(ecs, ec.Name)
		}

		dependencies = append(dependencies, "trigger-audit")

		parameters := []workflow.Parameter{
			{
				Name:  "workflowtemplate",
				Value: AnyStringPointer(workflowTemplate),
			},
			{
				Name:  "customer_id",
				Value: AnyStringPointer(env.Config.CompanyName),
			},
			{
				Name:  "team_name",
				Value: AnyStringPointer(environment.Spec.TeamName),
			},
			{
				Name:  "env_name",
				Value: AnyStringPointer(environment.Spec.EnvName),
			},
			{
				Name:  "config_name",
				Value: AnyStringPointer(ec.Name),
			},
			{
				Name:  "il_repo",
				Value: AnyStringPointer(env.Config.TFILRepoURL),
			},
			{
				Name:  "terraform_il_path",
				Value: AnyStringPointer(tfPath),
			},
			{
				Name:  "is_destroy",
				Value: AnyStringPointer(destroyFlag),
			},
			{
				Name:  "reconcile_id",
				Value: AnyStringPointer("{{tasks.trigger-audit.outputs.parameters.reconcile_id}}"),
			},
			{
				Name:  "status",
				Value: AnyStringPointer("initializing"),
			},
			{
				Name:  "auto_approve",
				Value: AnyStringPointer(autoApproveFlag),
			},
			{
				Name:  "skip_component",
				Value: AnyStringPointer(skipComponent(ec.DestroyProtection, destroyFlag, environment.Spec.SelectiveReconcile, ec.Tags)),
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

	onExit := ""
	if env.Config.SlackWebhookURL != "" {
		onExit = "exit-handler"
	}

	w := &workflow.Workflow{
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
			OnExit:     onExit,
			PodGC:      &workflow.PodGC{Strategy: workflow.PodGCOnPodSuccess},
			Templates: []workflow.Template{
				{
					Name: "main",
					DAG: &workflow.DAGTemplate{
						Tasks: tasks,
					},
				},
				exitHandler(environment),
			},
		},
		Status: workflow.WorkflowStatus{
			StartedAt:  metav1.Time{Time: getStaticDate()},
			FinishedAt: metav1.Time{Time: getStaticDate()},
		},
	}

	return w
}

func exitHandler(e *stablev1.Environment) workflow.Template {
	return workflow.Template{
		Name: "exit-handler",
		Steps: []workflow.ParallelSteps{
			{
				Steps: []workflow.WorkflowStep{
					{
						Name: "exit-handler",
						TemplateRef: &workflow.TemplateRef{
							Name:     "slack-notification",
							Template: "send-completion",
						},
						When: "{{workflow.status}} != Succeeded",
						Arguments: workflow.Arguments{
							Parameters: []workflow.Parameter{
								{
									Name:  "WORKFLOW_STATUS",
									Value: AnyStringPointer("{{workflow.status}}"),
								},
								{
									Name:  "WORKFLOW_NAME",
									Value: AnyStringPointer("{{workflow.name}}"),
								},
								{
									Name:  "SLACK_WEBHOOK_URL",
									Value: AnyStringPointer(env.Config.SlackWebhookURL),
								},
								{
									Name:  "WORKFLOW_TEAM",
									Value: AnyStringPointer(e.Spec.TeamName),
								},
								{
									Name:  "WORKFLOW_ENVIRONMENT",
									Value: AnyStringPointer(e.Spec.EnvName),
								},
								{
									Name:  "WORKFLOW_FAILURES",
									Value: AnyStringPointer("{{workflow.failures}}"),
								},
								{
									Name: "WORKFLOW_URL",
									Value: AnyStringPointer(fmt.Sprintf(
										"https://%s.zlifecycle.com/%s/%s/infra",
										env.Config.CompanyName,
										e.Spec.TeamName,
										e.Spec.TeamName+"-"+e.Spec.EnvName,
									)),
								},
							},
						},
					},
				},
			},
		},
	}
}

func skipComponent(destroyProtection bool, destroyFlag bool, selectiveReconcile *stablev1.SelectiveReconcile, tags []*stablev1.Tags) string {
	noSkipStatus := "noSkip"
	selectiveReconcileStatus := "selectiveReconcile"
	destroyProtectionStatus := "destroyProtection"

	if destroyProtection && destroyFlag {
		return destroyProtectionStatus
	}

	if selectiveReconcile == nil {
		return noSkipStatus
	}

	if tags == nil {
		if selectiveReconcile.SkipMode {
			return noSkipStatus
		} else {
			return selectiveReconcileStatus
		}
	}

	for _, tag := range tags {
		if tag.Name == selectiveReconcile.TagName {
			for _, sTag := range selectiveReconcile.TagValues {
				if selectiveReconcile.SkipMode && sTag == tag.Value {
					return selectiveReconcileStatus
				} else if !selectiveReconcile.SkipMode && sTag == tag.Value {
					return noSkipStatus
				}
			}
		}
	}

	if selectiveReconcile.SkipMode {
		return noSkipStatus
	}
	return selectiveReconcileStatus
}

func generateAuditTask(environment *stablev1.Environment, destroyAll bool, phase string, dependencies []string) workflow.DAGTask {
	var name string
	var reconcileID string
	var status string

	if phase == "0" {
		name = "trigger-audit"
		status = "initializing"
		reconcileID = "0"
	} else {
		name = "end-audit"
		status = "ended"
		reconcileID = "{{tasks.trigger-audit.outputs.parameters.reconcile_id}}"
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
					Value: AnyStringPointer(environment.Spec.TeamName),
				},
				{
					Name:  "env_name",
					Value: AnyStringPointer(environment.Spec.EnvName),
				},
				{
					Name:  "status",
					Value: AnyStringPointer(status),
				},
				{
					Name:  "is_destroy",
					Value: AnyStringPointer(destroyAll),
				},
				{
					Name:  "phase",
					Value: AnyStringPointer(phase),
				},
				{
					Name:  "reconcile_id",
					Value: AnyStringPointer(reconcileID),
				},
			},
		},
	}

	return task
}

func getStaticDate() time.Time {
	layout := "2006-01-02T15:04:05.000Z"
	someDate := "2019-06-25T15:04:05.000Z"

	t, _ := time.Parse(layout, someDate)

	return t
}

func AnyStringPointer(val interface{}) *workflow.AnyString {
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
