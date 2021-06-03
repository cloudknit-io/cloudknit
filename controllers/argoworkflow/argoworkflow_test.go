package argoworkflow

import (
	"github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	"github.com/compuzest/zlifecycle-il-operator/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateLegacyWorkflowOfWorkflows(t *testing.T) {
	env := mocks.GetMockEnv(false)

	wow := GenerateLegacyWorkflowOfWorkflows(env)

	assert.Equal(t, wow.DeletionTimestamp.IsZero(), true)

	tasks := wow.Spec.Templates[0].DAG.Tasks

	task1 := findTasks(tasks, "networking")
	assert.NotNil(t, task1)
	assert.ElementsMatch(t, task1.Dependencies, []string{})
	destroyFlag1 := findParam(task1.Arguments.Parameters, "destroy_flag")
	assert.NotNil(t, destroyFlag1)
	assert.Equal(t, *destroyFlag1.Value, "false")

	task2 := findTasks(tasks, "rebrand")
	assert.NotNil(t, task2)
	assert.ElementsMatch(t, task2.Dependencies, []string{"networking"})
	destroyFlag2 := findParam(task2.Arguments.Parameters, "destroy_flag")
	assert.NotNil(t, destroyFlag2)
	assert.Equal(t, *destroyFlag2.Value, "false")

	task3 := findTasks(tasks, "overlay")
	assert.NotNil(t, task3)
	assert.ElementsMatch(t, task3.Dependencies, []string{"networking", "rebrand"})
	destroyFlag3 := findParam(task3.Arguments.Parameters, "destroy_flag")
	assert.NotNil(t, destroyFlag3)
	assert.Equal(t, *destroyFlag3.Value, "false")
}

func TestGenerateLegacyWorkflowOfWorkflowsDeletedEnvironment(t *testing.T) {
	env := mocks.GetMockEnv(true)

	wow := GenerateLegacyWorkflowOfWorkflows(env)

	assert.Equal(t, wow.DeletionTimestamp.IsZero(), true)

	tasks := wow.Spec.Templates[0].DAG.Tasks

	task1 := findTasks(tasks, "networking")
	assert.NotNil(t, task1)
	assert.ElementsMatch(t, task1.Dependencies, []string{"rebrand", "overlay"})
	destroyFlag1 := findParam(task1.Arguments.Parameters, "destroy_flag")
	assert.NotNil(t, destroyFlag1)
	assert.Equal(t, *destroyFlag1.Value, "true")

	task2 := findTasks(tasks, "rebrand")
	assert.NotNil(t, task2)
	assert.ElementsMatch(t, task2.Dependencies, []string{"overlay"})
	destroyFlag2 := findParam(task2.Arguments.Parameters, "destroy_flag")
	assert.NotNil(t, destroyFlag2)
	assert.Equal(t, *destroyFlag2.Value, "true")

	task3 := findTasks(tasks, "overlay")
	assert.NotNil(t, task3)
	assert.ElementsMatch(t, task3.Dependencies, []string{})
	destroyFlag3 := findParam(task3.Arguments.Parameters, "destroy_flag")
	assert.NotNil(t, destroyFlag3)
	assert.Equal(t, *destroyFlag3.Value, "true")
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
