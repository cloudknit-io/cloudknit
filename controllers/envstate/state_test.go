package envstate

import (
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/common"
	"github.com/compuzest/zlifecycle-il-operator/mocks"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	controllerruntime "sigs.k8s.io/controller-runtime"
	"testing"
)

func TestGetEnvironmentStateDiff(t *testing.T) {
	mockEnv1 := mocks.GetMockEnv1(false)
	cm := v1.ConfigMap{
		ObjectMeta: controllerruntime.ObjectMeta{
			Name: "environment-state-cm",
			Namespace: "test",
		},
		Data: make(map[string]string),
	}
	err := UpsertStateEntry(&cm, &mockEnv1)
	assert.NoError(t, err)

	mockEnv2 := mocks.GetMockEnv2(false)

	d, err := GetEnvironmentStateDiff(&cm, &mockEnv2)
	assert.NoError(t, err)
	assert.Len(t, d, 1)
	assert.Equal(t, d[0].Name, "overlay")
}

func TestDiff(t *testing.T) {
	mockComponents1 := mocks.GetMockEnv1(false).Spec.EnvironmentComponent
	mockComponents2 := mocks.GetMockEnv2(false).Spec.EnvironmentComponent

	d := diff(mockComponents1, mockComponents2)
	assert.Len(t, d, 1)
	assert.Equal(t, d[0].Name, "overlay")
}

func TestUpdateOrCreateStateEntryEmptyData(t *testing.T) {
	mockEnv := mocks.GetMockEnv1(false)
	mockTeamName := mockEnv.Spec.TeamName
	mockEnvName := mockEnv.Spec.EnvName
	cm := v1.ConfigMap{
		ObjectMeta: controllerruntime.ObjectMeta{
			Name: "environment-state-cm",
			Namespace: "test",
		},
		Data: make(map[string]string),
	}
	err := UpsertStateEntry(&cm, &mockEnv)
	assert.NoError(t, err)

	ymlstr := cm.Data[mockTeamName]
	ts := TeamState{}
	err = common.FromYaml(ymlstr, &ts)
	assert.NoError(t, err, "Error parsing TeamState from yaml string")
	assert.Equal(t, ts.Name, mockTeamName)
	assert.Equal(t, ts.Environments[0].Name, mockEnvName)
	assert.Len(t, ts.Environments[0].EnvironmentComponents, 3)
}

func TestUpdateOrCreateStateEntryExistingState(t *testing.T) {
	mockEnv1 := mocks.GetMockEnv1(false)
	mockTeamName := mockEnv1.Spec.TeamName
	mockEnvName  := mockEnv1.Spec.EnvName
	mockTs := buildNewTeamState(&mockEnv1)
	ymlstring, _ := common.ToYaml(mockTs)
	mockData := make(map[string]string)
	mockData[mockTeamName] = ymlstring
	cm := v1.ConfigMap{
		ObjectMeta: controllerruntime.ObjectMeta{
			Name: "environment-state-cm",
			Namespace: "test",
		},
		Data: mockData,
	}

	mockEnv2 := mocks.GetMockEnv2(false)
	err := UpsertStateEntry(&cm, &mockEnv2)
	assert.NoError(t, err)

	ymlstr := cm.Data[mockTeamName]
	ts := TeamState{}
	err = common.FromYaml(ymlstr, &ts)
	assert.NoError(t, err, "Error parsing TeamState from yaml string")
	assert.Equal(t, ts.Name, mockTeamName)
	assert.Equal(t, ts.Environments[0].Name, mockEnvName)
	assert.Len(t, ts.Environments[0].EnvironmentComponents, 2)
}

func TestUpdateOrCreateStateEntryExistingStateNewEnv(t *testing.T) {
	mockEnv1 := mocks.GetMockEnv1(false)
	mockTeamName := mockEnv1.Spec.TeamName
	mockTs := buildNewTeamState(&mockEnv1)
	ymlstring, _ := common.ToYaml(mockTs)
	mockData := make(map[string]string)
	mockData[mockTeamName] = ymlstring
	cm := v1.ConfigMap{
		ObjectMeta: controllerruntime.ObjectMeta{
			Name: "environment-state-cm",
			Namespace: "test",
		},
		Data: mockData,
	}

	mockEnv3 := mocks.GetMockEnv3(false)
	err := UpsertStateEntry(&cm, &mockEnv3)
	assert.NoError(t, err)

	ymlstr := cm.Data[mockTeamName]
	ts := TeamState{}
	err = common.FromYaml(ymlstr, &ts)
	assert.NoError(t, err, "Error parsing TeamState from yaml string")
	assert.Equal(t, ts.Name, mockTeamName)
	assert.Len(t, ts.Environments, 2)
}

