package interpolator

import (
	"fmt"

	v1 "github.com/compuzest/zlifecycle-il-operator/api/v1"
	"github.com/compuzest/zlifecycle-il-operator/controller/util"
	"github.com/pkg/errors"
)

const (
	identifierLocal = "zlocals"
)

func InterpolateTFVars(tfvars string, zlocals []*v1.LocalVariable) (string, error) {
	vars := BuildZLocalsVariableMap(zlocals)

	interpolated, err := util.Interpolate(tfvars, vars)
	if err != nil {
		return "", errors.Wrapf(err, "error interpolating tfvars")
	}

	return interpolated, nil
}

func Interpolate(e v1.Environment) (*v1.Environment, error) { // nolint
	vars := BuildZLocalsVariableMap(e.Spec.ZLocals)

	teamName, err := util.Interpolate(e.Spec.TeamName, vars)
	if err != nil {
		return nil, errors.Wrapf(err, "error interpolating environment.spec.teamName: [%s]", e.Spec.TeamName)
	}
	e.Spec.TeamName = teamName

	envName, err := util.Interpolate(e.Spec.EnvName, vars)
	if err != nil {
		return nil, errors.Wrapf(err, "error interpolating environment.spec.envName: [%s]", e.Spec.EnvName)
	}
	e.Spec.EnvName = envName

	for i, ec := range e.Spec.Components {
		interpolated, err := interpolateComponent(*ec, i, vars)
		if err != nil {
			return nil, err
		}
		e.Spec.Components[i] = interpolated
	}

	return &e, nil
}

func interpolateComponent(ec v1.EnvironmentComponent, index int, vars util.Variables) (*v1.EnvironmentComponent, error) { // nolint
	name, err := util.Interpolate(ec.Name, vars)
	if err != nil {
		return nil, errors.Wrapf(err, "error interpolating environment.spec.components[%d].name: [%s]", index, ec.Name)
	}
	ec.Name = name

	if ec.VariablesFile != nil {
		interpolated, err := interpolateVariablesFile(*ec.VariablesFile, index, vars)
		if err != nil {
			return nil, err
		}
		ec.VariablesFile = interpolated
	}

	for vIndex, v := range ec.Variables {
		interpolated, err := interpolateVariable(*v, index, vIndex, vars)
		if err != nil {
			return nil, err
		}
		ec.Variables[vIndex] = interpolated
	}

	for tIndex, t := range ec.Tags {
		interpolated, err := interpolateTags(*t, index, tIndex, vars)
		if err != nil {
			return nil, err
		}
		ec.Tags[tIndex] = interpolated
	}

	return &ec, nil
}

func interpolateVariablesFile(vf v1.VariablesFile, ecIndex int, vars util.Variables) (*v1.VariablesFile, error) {
	source, err := util.Interpolate(vf.Source, vars)
	if err != nil {
		return nil, errors.Wrapf(err, "error interpolating environment.spec.components[%d].variableaFile.source: [%s]", ecIndex, vf.Source)
	}
	vf.Source = source

	path, err := util.Interpolate(vf.Path, vars)
	if err != nil {
		return nil, errors.Wrapf(err, "error interpolating environment.spec.components[%d].variableaFile.path: [%s]", ecIndex, vf.Path)
	}
	vf.Path = path

	ref, err := util.Interpolate(vf.Ref, vars)
	if err != nil {
		return nil, errors.Wrapf(err, "error interpolating environment.spec.components[%d].variableaFile.path: [%s]", ecIndex, vf.Ref)
	}
	vf.Ref = ref

	return &vf, nil
}

func interpolateVariable(v v1.Variable, ecIndex, varIndex int, vars util.Variables) (*v1.Variable, error) {
	name, err := util.Interpolate(v.Name, vars)
	if err != nil {
		return nil, errors.Wrapf(err, "error interpolating environment.spec.components[%d].variables[%d].name: [%s]", ecIndex, varIndex, name)
	}
	v.Name = name

	value, err := util.Interpolate(v.Value, vars)
	if err != nil {
		return nil, errors.Wrapf(err, "error interpolating environment.spec.components[%d].variables[%d].value: [%s]", ecIndex, varIndex, value)
	}
	v.Value = value

	valueFrom, err := util.Interpolate(v.ValueFrom, vars)
	if err != nil {
		return nil, errors.Wrapf(err, "error interpolating environment.spec.components[%d].variables[%d].valueFrom: [%s]", ecIndex, varIndex, valueFrom)
	}
	v.ValueFrom = valueFrom

	return &v, nil
}

func interpolateTags(t v1.Tag, ecIndex, tagIndex int, vars util.Variables) (*v1.Tag, error) {
	name, err := util.Interpolate(t.Name, vars)
	if err != nil {
		return nil, errors.Wrapf(err, "error interpolating environment.spec.components[%d].tags[%d].name: [%s]", ecIndex, tagIndex, name)
	}
	t.Name = name

	value, err := util.Interpolate(t.Value, vars)
	if err != nil {
		return nil, errors.Wrapf(err, "error interpolating environment.spec.components[%d].tags[%d].value: [%s]", ecIndex, tagIndex, value)
	}
	t.Value = value

	return &t, nil
}

func BuildZLocalsVariableMap(zlocals []*v1.LocalVariable) util.Variables {
	vars := make(map[string]string, len(zlocals))

	for _, l := range zlocals {
		scopedName := fmt.Sprintf("%s.%s", identifierLocal, l.Name)
		vars[scopedName] = l.Value
	}

	return vars
}
