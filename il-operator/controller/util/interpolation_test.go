package util_test

import (
	"testing"

	"github.com/compuzest/zlifecycle-il-operator/controller/util"
	"github.com/stretchr/testify/assert"
)

func TestInterpolate(t *testing.T) {
	t.Parallel()

	vars1 := map[string]string{"global.foo": "Hello", "global.bar": "world"}
	template1 := "${global.foo} ${global.bar}!"
	expected1 := "Hello world!"
	interpolated1, err := util.Interpolate(template1, vars1)
	assert.NoError(t, err)
	assert.Equal(t, interpolated1, expected1)

	vars2 := map[string]string{"global.bar": "Hello"}
	template2 := "${global.bar} Dejan, ${global.bar} Adarsh"
	expected2 := "Hello Dejan, Hello Adarsh"
	interpolated2, err := util.Interpolate(template2, vars2)
	assert.NoError(t, err)
	assert.Equal(t, interpolated2, expected2)

	vars3 := map[string]string{"local.test": "Programming"}
	template3 := "${local.test} is fun but ${global.else} is funnier"
	expected3 := "Programming is fun but ${global.else} is funnier"
	interpolated3, err := util.Interpolate(template3, vars3)
	assert.NoError(t, err)
	assert.Equal(t, interpolated3, expected3)

	vars4 := map[string]string{"local.test": "Programming", "global.test": "drinking"}
	template4 := "${local.test} is fun but ${global.test} is funnier"
	expected4 := "Programming is fun but drinking is funnier"
	interpolated4, err := util.Interpolate(template4, vars4)
	assert.NoError(t, err)
	assert.Equal(t, interpolated4, expected4)

	vars5 := map[string]string{"local.test": "Programming", "global.test": "drinking"}
	template5 := "${test} is fun but ${global.test} is funnier"
	_, err = util.Interpolate(template5, vars5)
	assert.Error(t, err)
}
