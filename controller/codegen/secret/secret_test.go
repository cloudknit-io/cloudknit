package secret_test

import (
	"testing"

	"github.com/compuzest/zlifecycle-il-operator/controller/codegen/secret"

	"github.com/stretchr/testify/assert"
)

func TestCreateKey(t *testing.T) {
	t.Parallel()

	meta1 := secret.Identifier{
		Company:              "zbank",
		Team:                 "payments",
		Environment:          "dev",
		EnvironmentComponent: "overlay",
	}
	key1, err1 := meta1.GenerateKey("secret1", "component")
	expected1 := "/zbank/payments/dev/overlay/secret1"
	assert.NoError(t, err1)
	assert.Equal(t, key1, expected1)

	meta2 := secret.Identifier{
		Company:     "zbank",
		Team:        "checkout",
		Environment: "prod",
	}
	key2, err2 := meta2.GenerateKey("secret1", "environment")
	expected2 := "/zbank/checkout/prod/secret1"
	assert.NoError(t, err2)
	assert.Equal(t, key2, expected2)

	meta3 := secret.Identifier{
		Company: "zbank",
		Team:    "platform",
	}
	key3, err3 := meta3.GenerateKey("secret1", "team")
	expected3 := "/zbank/platform/secret1"
	assert.NoError(t, err3)
	assert.Equal(t, key3, expected3)

	meta4 := secret.Identifier{
		Company: "zbank",
	}
	key4, err4 := meta4.GenerateKey("secret1", "org")
	expected4 := "/zbank/secret1"
	assert.NoError(t, err4)
	assert.Equal(t, key4, expected4)

	meta5 := secret.Identifier{
		Company: "zbank",
		Team:    "platform",
	}
	_, err5 := meta5.GenerateKey("secret1", "component")
	assert.Error(t, err5)
}
