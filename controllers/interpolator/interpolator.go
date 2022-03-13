package interpolator

import (
	"fmt"
	v1 "github.com/compuzest/zlifecycle-il-operator/api/v1"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util"
	"github.com/pkg/errors"
)

const (
	identifierLocal = "zlocals"
)

type Variables map[string]string

func Interpolate(e *v1.Environment) (*v1.Environment, error) {
	vars := BuildZLocalsVariableMap(e.Spec.ZLocals)

	teamName, err := util.Interpolate(e.Spec.TeamName, vars)
	if err != nil {
		return nil, errors.Wrapf(err, "error interpolating environment.spec.teamName: [%s]", teamName)
	}
	e.Spec.TeamName = teamName

	envName, err := util.Interpolate(e.Spec.EnvName, vars)
	if err != nil {
		return nil, errors.Wrapf(err, "error interpolating environment.spec.envName: [%s]", envName)
	}
	e.Spec.EnvName = envName

	for i, ec := range e.Spec.Components {
		interpolated, err := interpolateComponent(ec, i, vars)
		if err != nil {
			return nil, err
		}
		e.Spec.Components[i] = interpolated
	}

	return e, nil
}

func interpolateComponent(ec *v1.EnvironmentComponent, index int, vars Variables) (*v1.EnvironmentComponent, error) {
	name, err := util.Interpolate(ec.Name, vars)
	if err != nil {
		return nil, errors.Wrapf(err, "error interpolating environment.spec.components[%d].name: [%s]", index, name)
	}
	ec.Name = name

	for vIndex, v := range ec.Variables {
		interpolated, err := interpolateVariable(v, index, vIndex, vars)
		if err != nil {
			return nil, err
		}
		ec.Variables[vIndex] = interpolated
	}

	return ec, nil
}

func interpolateVariable(v *v1.Variable, ecIndex, varIndex int, vars Variables) (*v1.Variable, error) {
	name, err := util.Interpolate(v.Name, vars)
	if err != nil {
		return nil, errors.Wrapf(err, "error interpolating environment.spec.components[%d].variables[%d].name: [%s]", ecIndex, varIndex, name)
	}
	v.Name = name

	value, err := util.Interpolate(v.Name, vars)
	if err != nil {
		return nil, errors.Wrapf(err, "error interpolating environment.spec.components[%d].variables[%d].value: [%s]", ecIndex, varIndex, value)
	}
	v.Value = value

	valueFrom, err := util.Interpolate(v.Name, vars)
	if err != nil {
		return nil, errors.Wrapf(err, "error interpolating environment.spec.components[%d].variables[%d].valueFrom: [%s]", ecIndex, varIndex, valueFrom)
	}
	v.ValueFrom = valueFrom

	return v, nil
}

func BuildZLocalsVariableMap(zlocals []*v1.LocalVariable) Variables {
	vars := make(map[string]string, len(zlocals))

	for _, l := range zlocals {
		scopedName := fmt.Sprintf("%s.%s", identifierLocal, l.Name)
		vars[scopedName] = l.Value
	}

	return vars
}
