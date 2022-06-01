package validator

import (
	"testing"

	v1 "github.com/compuzest/zlifecycle-il-operator/api/v1"
	"github.com/stretchr/testify/assert"
)

func TestGetOutputFromComponent(t *testing.T) {
	t.Parallel()

	ec := &v1.EnvironmentComponent{
		Outputs: []*v1.Output{
			{Name: "testOutput"},
			{Name: "anotherOutput"},
		},
	}

	output := GetOutputFromComponent("testOutput", ec)
	assert.IsType(t, &v1.Output{}, output)

	output = GetOutputFromComponent("doesNotExist", ec)
	assert.Nil(t, output)
}

func TestSplitValueFrom(t *testing.T) {
	t.Parallel()

	comp, variable, err := SplitValueFrom("blah.varname")

	assert.Equal(t, "blah", comp)
	assert.Equal(t, "varname", variable)
	assert.Nil(t, err)

	comp, variable, err = SplitValueFrom("blah.varname[0]")

	assert.Equal(t, "blah", comp)
	assert.Equal(t, "varname", variable)
	assert.Nil(t, err)

	comp, variable, err = SplitValueFrom("blah.varname[0")

	assert.Equal(t, "", comp)
	assert.Equal(t, "", variable)
	assert.NotNil(t, err)

	comp, variable, err = SplitValueFrom("blah.varname[jfkd]")

	assert.Equal(t, "", comp)
	assert.Equal(t, "", variable)
	assert.NotNil(t, err)

	comp, variable, err = SplitValueFrom("blah-varname")

	assert.Equal(t, "", comp)
	assert.Equal(t, "", variable)
	assert.NotNil(t, err)
}
