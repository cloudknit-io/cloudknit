package interpolator_test

import (
	v1 "github.com/compuzest/zlifecycle-il-operator/api/v1"
	"github.com/compuzest/zlifecycle-il-operator/controllers/interpolator"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInterpolate(t *testing.T) {
	t.Parallel()

	e := &v1.Environment{
		Spec: v1.EnvironmentSpec{
			ZLocals: []*v1.LocalVariable{
				{Name: "env", Value: "dev"},
				{Name: "team", Value: "checkout"},
				{Name: "component", Value: "k8s"},
			},
			EnvName:    "test-${zlocals.env}",
			TeamName:   "test-${zlocals.team}",
			Components: []*v1.EnvironmentComponent{{Name: "test-ec-${zlocals.component}"}},
		},
	}

	expected := &v1.Environment{
		Spec: v1.EnvironmentSpec{
			ZLocals: []*v1.LocalVariable{
				{Name: "test", Value: "dev"},
				{Name: "team", Value: "checkout"},
				{Name: "component", Value: "k8s"},
			},
			EnvName:    "test-dev",
			TeamName:   "test-checkout",
			Components: []*v1.EnvironmentComponent{{Name: "test-ec-k8s"}},
		},
	}

	interpolated, err := interpolator.Interpolate(e)
	assert.NoError(t, err)
	assert.Equal(t, interpolated.Spec.TeamName, expected.Spec.TeamName)
	assert.Equal(t, interpolated.Spec.EnvName, expected.Spec.EnvName)
	assert.Equal(t, interpolated.Spec.Components[0].Name, expected.Spec.Components[0].Name)
}

func TestBuildVariableMap(t *testing.T) {
	t.Parallel()

	zlocals := []*v1.LocalVariable{{Name: "foo", Value: "bar"}, {Name: "test", Value: "baz"}}
	actualVars := interpolator.BuildZLocalsVariableMap(zlocals)
	expectedVars := map[string]string{"zlocals.foo": "bar", "zlocals.test": "baz"}
	assert.EqualValues(t, actualVars, expectedVars)
}
