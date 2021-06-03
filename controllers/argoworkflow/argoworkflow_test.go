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

	task1 := find(tasks, "networking")
	assert.NotNil(t, task1)
	assert.ElementsMatch(t, task1.Dependencies, []string{})

	task2 := find(tasks, "rebrand")
	assert.NotNil(t, task2)
	assert.ElementsMatch(t, task2.Dependencies, []string{"networking"})

	task3 := find(tasks, "overlay")
	assert.NotNil(t, task3)
	assert.ElementsMatch(t, task3.Dependencies, []string{"networking", "rebrand"})
}

func TestGenerateLegacyWorkflowOfWorkflowsDeletedEnvironment(t *testing.T) {
	env := mocks.GetMockEnv(true)

	wow := GenerateLegacyWorkflowOfWorkflows(env)

	assert.Equal(t, wow.DeletionTimestamp.IsZero(), true)

	tasks := wow.Spec.Templates[0].DAG.Tasks

	task1 := find(tasks, "networking")
	assert.NotNil(t, task1)
	assert.ElementsMatch(t, task1.Dependencies, []string{"rebrand", "overlay"})

	task2 := find(tasks, "rebrand")
	assert.NotNil(t, task2)
	assert.ElementsMatch(t, task2.Dependencies, []string{"overlay"})

	task3 := find(tasks, "overlay")
	assert.NotNil(t, task3)
	assert.ElementsMatch(t, task3.Dependencies, []string{})
}

func find(tasks []v1alpha1.DAGTask, name string) *v1alpha1.DAGTask {
	for _, task := range tasks {
		if task.Name == name {
			return &task
		}
	}
	return nil
}
