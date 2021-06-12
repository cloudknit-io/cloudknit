package envstate

import (
	stablev1alpha1 "github.com/compuzest/zlifecycle-il-operator/api/v1alpha1"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/common"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/env"
	v1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sync"
)

var (
	mutex = sync.Mutex{}
)

func GetEnvironmentStateObjectKey() client.ObjectKey {
	return client.ObjectKey{
		Name: env.Config.EnvironmentStateConfigMap,
		Namespace: env.Config.ZlifecycleOperatorNamespace,
	}
}

func GetEnvironmentStateConfigMap() *v1.ConfigMap {
	return &v1.ConfigMap{
		ObjectMeta: ctrl.ObjectMeta{
			Name: env.Config.EnvironmentStateConfigMap,
			Namespace: env.Config.ZlifecycleOperatorNamespace,
		},
	}
}

func GetEnvironmentStateDiff(
	cm *v1.ConfigMap,
	e *stablev1alpha1.Environment,
	) (d []*stablev1alpha1.EnvironmentComponent, err error) {
	teamName := e.Spec.TeamName
	envName := e.Spec.EnvName
	ymlstr := cm.Data[teamName]
	ts := TeamState{}
	if err := common.FromYaml(ymlstr, &ts); err != nil {
		return nil, err
	}
	var oldState []*stablev1alpha1.EnvironmentComponent
	for _, es := range ts.Environments {
		if es.Name == envName {
			oldState = es.EnvironmentComponents
		}
	}
	newState := e.Spec.EnvironmentComponent

	return diff(oldState, newState), nil
}

func diff(
	old []*stablev1alpha1.EnvironmentComponent,
	new []*stablev1alpha1.EnvironmentComponent,
	) []*stablev1alpha1.EnvironmentComponent {
	var d []*stablev1alpha1.EnvironmentComponent
	for _, oldEc := range old {
		found := false
		for _, newEc := range new {
			if oldEc.Name == newEc.Name {
				found = true
				break
			}
		}
		if !found {
			d = append(d, oldEc)
		}
	}
	return d
}

func DeleteStateEntry(cm *v1.ConfigMap, e *stablev1alpha1.Environment) {
	mutex.Lock()
	defer mutex.Unlock()

	if cm.Data == nil {
		return
	}
	team := e.Spec.TeamName
	cm.Data[team] = ""

	return
}

func UpsertStateEntry(cm *v1.ConfigMap, e *stablev1alpha1.Environment) error {
	mutex.Lock()
	defer mutex.Unlock()
	if cm.Data == nil {
		cm.Data = make(map[string]string)
	}
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
	ts.Environments[envName] = es
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

func buildNewEnvironmentState(e *stablev1alpha1.Environment) EnvironmentState {
	return EnvironmentState{
		Name: e.Spec.EnvName,
		EnvironmentComponents: e.Spec.EnvironmentComponent,
	}
}

func buildNewTeamState(e *stablev1alpha1.Environment) TeamState {
	es := EnvironmentState{
		Name: e.Spec.EnvName,
		EnvironmentComponents: e.Spec.EnvironmentComponent,
	}

	m := make(map[string]EnvironmentState)
	envName := e.Spec.EnvName
	m[envName] = es
	return TeamState{
		Name: e.Spec.TeamName,
		Environments: m,
	}
}
