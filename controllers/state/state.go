package state

import (
	stablev1alpha1 "github.com/compuzest/zlifecycle-il-operator/api/v1alpha1"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/common"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/env"
	v1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

func CreateEnvironmentStateConfigMap() *v1.ConfigMap {
	return &v1.ConfigMap{
		ObjectMeta: ctrl.ObjectMeta{
			Name: env.Config.EnvironmentStateConfigMap,
			Namespace: env.Config.ZlifecycleOperatorNamespace,
		},
		Data: make(map[string]string),
	}
}

func UpsertStateEntry(cm *v1.ConfigMap, e *stablev1alpha1.Environment) error {
	team := e.Spec.TeamName
	ymlstr := cm.Data[team]

	entryExists := ymlstr != ""
	if entryExists {
		newYmlStr, err := updateStateEntry(e, ymlstr)
		if err != nil {
			return err
		}
		cm.Data[team] = newYmlStr
	} else {
		ymlstr, err := createStateEntry(e)
		if err != nil {
			return err
		}
		cm.Data[team] = ymlstr
	}
	return nil
}

func updateStateEntry(e *stablev1alpha1.Environment, ymlstr string) (string, error) {
	envName := e.Spec.EnvName
	ts := TeamState{}
	if err := common.FromYaml(ymlstr, &ts); err != nil {
		return "", err
	}
	es := buildNewEnvironmentState(e)
	if i := indexOf(ts, envName); i != -1 {
		ts.Environments[i] = es
	} else {
		ts.Environments = append(ts.Environments, es)
	}
	newYmlstr, err := common.ToYaml(&ts)
	if err != nil {
		return "", err
	}
	return newYmlstr, nil
}

func createStateEntry(e *stablev1alpha1.Environment) (string, error) {
	ts := buildNewTeamState(e)
	ymlstr, err := common.ToYaml(ts)
	if err != nil {
		return "", err
	}
	return ymlstr, nil
}

func indexOf(ts TeamState, name string) int {
	for i, e := range ts.Environments {
		if e.Name == name {
			return i
		}
	}
	return -1
}

func find(ts TeamState, name string) *EnvironmentState {
	for _, e := range ts.Environments {
		if e.Name == name {
			return &e
		}
	}
	return nil
}

func buildNewEnvironmentState(e *stablev1alpha1.Environment) EnvironmentState {
	var ecs []stablev1alpha1.EnvironmentComponent
	for _, c := range e.Spec.EnvironmentComponent {
		ecs = append(ecs, *c)
	}
	return EnvironmentState{
		Name: e.Spec.EnvName,
		EnvironmentComponents: ecs,
	}
}

func buildNewTeamState(e *stablev1alpha1.Environment) TeamState {
	var ecs []stablev1alpha1.EnvironmentComponent
	for _, c := range e.Spec.EnvironmentComponent {
		ecs = append(ecs, *c)
	}
	es := EnvironmentState{
		Name: e.Spec.EnvName,
		EnvironmentComponents: ecs,
	}

	return TeamState{
		Name: e.Spec.TeamName,
		Environments: []EnvironmentState{es},
	}
}
