package argoworkflow_test

import (
	"testing"

	"github.com/compuzest/zlifecycle-il-operator/controller/codegen/workflow"

	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/compuzest/zlifecycle-il-operator/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGenerateWorkflowOfWorkflows(t *testing.T) {
	t.Parallel()

	env := mocks.GetMockEnv1(false)

	wow := workflow.GenerateWorkflowOfWorkflows(&env, nil)

	assert.Equal(t, wow.DeletionTimestamp.IsZero(), true)

	tasks := wow.Spec.Templates[0].DAG.Tasks

	task1 := findTasks(tasks, "networking")
	assert.NotNil(t, task1)
	assert.ElementsMatch(t, task1.Dependencies, []string{"trigger-audit"})
	destroyFlag1 := findParam(task1.Arguments.Parameters, "is_destroy")
	assert.NotNil(t, destroyFlag1)
	assert.Equal(t, destroyFlag1.Value, workflow.AnyStringPointer("false"))
	autoApproveFlag1 := findParam(task1.Arguments.Parameters, "auto_approve")
	assert.NotNil(t, autoApproveFlag1)
	assert.Equal(t, autoApproveFlag1.Value, workflow.AnyStringPointer("false"))

	task2 := findTasks(tasks, "rebrand")
	assert.NotNil(t, task2)
	assert.ElementsMatch(t, task2.Dependencies, []string{"networking", "trigger-audit"})
	destroyFlag2 := findParam(task2.Arguments.Parameters, "is_destroy")
	assert.NotNil(t, destroyFlag2)
	assert.Equal(t, destroyFlag2.Value, workflow.AnyStringPointer("false"))
	autoApproveFlag2 := findParam(task2.Arguments.Parameters, "auto_approve")
	assert.NotNil(t, autoApproveFlag2)
	assert.Equal(t, autoApproveFlag2.Value, workflow.AnyStringPointer("true"))

	task3 := findTasks(tasks, "overlay")
	assert.NotNil(t, task3)
	assert.ElementsMatch(t, task3.Dependencies, []string{"networking", "rebrand", "trigger-audit"})
	destroyFlag3 := findParam(task3.Arguments.Parameters, "is_destroy")
	assert.NotNil(t, destroyFlag3)
	assert.Equal(t, destroyFlag3.Value, workflow.AnyStringPointer("false"))
	autoApproveFlag3 := findParam(task3.Arguments.Parameters, "auto_approve")
	assert.NotNil(t, autoApproveFlag3)
	assert.Equal(t, autoApproveFlag3.Value, workflow.AnyStringPointer("true"))
}

func TestGenerateWorkflowOfWorkflowsDeletedEnvironment(t *testing.T) {
	t.Parallel()

	env := mocks.GetMockEnv1(true)

	wow := workflow.GenerateWorkflowOfWorkflows(&env, nil)

	assert.Equal(t, wow.DeletionTimestamp.IsZero(), true)

	tasks := wow.Spec.Templates[0].DAG.Tasks

	task1 := findTasks(tasks, "networking")
	assert.NotNil(t, task1)
	assert.ElementsMatch(t, task1.Dependencies, []string{"rebrand", "overlay", "trigger-audit"})
	destroyFlag1 := findParam(task1.Arguments.Parameters, "is_destroy")
	assert.NotNil(t, destroyFlag1)
	assert.Equal(t, destroyFlag1.Value, workflow.AnyStringPointer("true"))

	task2 := findTasks(tasks, "rebrand")
	assert.NotNil(t, task2)
	assert.ElementsMatch(t, task2.Dependencies, []string{"overlay", "trigger-audit"})
	destroyFlag2 := findParam(task2.Arguments.Parameters, "is_destroy")
	assert.NotNil(t, destroyFlag2)
	assert.Equal(t, destroyFlag2.Value, workflow.AnyStringPointer("true"))

	task3 := findTasks(tasks, "overlay")
	assert.NotNil(t, task3)
	assert.ElementsMatch(t, task3.Dependencies, []string{"trigger-audit"})
	destroyFlag3 := findParam(task3.Arguments.Parameters, "is_destroy")
	assert.NotNil(t, destroyFlag3)
	assert.Equal(t, destroyFlag3.Value, workflow.AnyStringPointer("true"))
}

func findTasks(tasks []v1alpha1.DAGTask, name string) *v1alpha1.DAGTask {
	for _, task := range tasks {
		if task.Name == name {
			return &task
		}
	}
	return nil
}

func findParam(params []v1alpha1.Parameter, name string) *v1alpha1.Parameter {
	for _, param := range params {
		if param.Name == name {
			return &param
		}
	}
	return nil
}
