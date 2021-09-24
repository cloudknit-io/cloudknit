package secrets_test

import (
	"testing"

	"github.com/compuzest/zlifecycle-il-operator/controllers/secrets"
	"github.com/stretchr/testify/assert"
)

func TestCreateKey(t *testing.T) {
	t.Parallel()

	meta1 := secrets.SecretMeta{
		Company:              "zbank",
		Team:                 "payments",
		Environment:          "dev",
		EnvironmentComponent: "overlay",
	}
	key1, err1 := secrets.CreateKey("secret1", "component", meta1)
	expected1 := "/zbank/payments/dev/overlay/secret1"
	assert.NoError(t, err1)
	assert.Equal(t, key1, expected1)

	meta2 := secrets.SecretMeta{
		Company:     "zbank",
		Team:        "checkout",
		Environment: "prod",
	}
	key2, err2 := secrets.CreateKey("secret1", "environment", meta2)
	expected2 := "/zbank/checkout/prod/secret1"
	assert.NoError(t, err2)
	assert.Equal(t, key2, expected2)

	meta3 := secrets.SecretMeta{
		Company: "zbank",
		Team:    "platform",
	}
	key3, err3 := secrets.CreateKey("secret1", "team", meta3)
	expected3 := "/zbank/platform/secret1"
	assert.NoError(t, err3)
	assert.Equal(t, key3, expected3)

	meta4 := secrets.SecretMeta{
		Company: "zbank",
	}
	key4, err4 := secrets.CreateKey("secret1", "org", meta4)
	expected4 := "/zbank/secret1"
	assert.NoError(t, err4)
	assert.Equal(t, key4, expected4)

	meta5 := secrets.SecretMeta{
		Company: "zbank",
		Team:    "platform",
	}
	_, err5 := secrets.CreateKey("secret1", "component", meta5)
	assert.Error(t, err5)
}
