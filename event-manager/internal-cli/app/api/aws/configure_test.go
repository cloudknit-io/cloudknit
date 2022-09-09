package aws_test

import (
	"testing"

	"github.com/compuzest/zlifecycle-internal-cli/app/api/aws"
	"github.com/compuzest/zlifecycle-internal-cli/app/env"
	"github.com/stretchr/testify/assert"
)

func TestGenerateAWSCredentialsEntry(t *testing.T) {
	t.Parallel()
	if env.TestMode != env.TestModeUnit && env.TestMode != env.TestModeAll {
		t.Skip()
	}

	creds := aws.GenerateAWSCredentialsEntry("client-test", "xxx", "yyy", "us-west-1")

	expected := `[client-test]
aws_access_key_id = xxx
aws_secret_access_key = yyy
region = us-west-1`

	assert.Equal(t, creds, expected)
}
