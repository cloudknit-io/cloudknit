package argocd

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestGetArgocdServerAddrExistingCreds(t *testing.T) {
	expected := ArgocdCredentials{Username: "test", Password: "test"}
	os.Setenv("ARGOCD_USERNAME", expected.Username)
	os.Setenv("ARGOCD_PASSWORD", expected.Password)
	r, _ := getArgocdCredentialsFromEnv()

	assert.Equal(t, expected.Username, r.Username)
	assert.Equal(t, expected.Password, r.Password)
	os.Clearenv()
}

func TestGetArgocdServerAddrMissingCreds(t *testing.T) {
	expected := ArgocdCredentials{Username: "test"}
	os.Setenv("ARGOCD_USERNAME", expected.Username)
	_, err := getArgocdCredentialsFromEnv()
	assert.Error(t, err)
}

func TestGetArgocdServerAddrExistingUrl(t *testing.T) {
	expected := "https://argocd.test.com"
	os.Setenv("ARGOCD_URL", expected)
	r := GetArgocdServerAddr()
	assert.Equal(t, expected, r)
	os.Clearenv()
}

func TestGetArgocdServerAddrDefaultUrl(t *testing.T) {
	expected := "http://argocd-server.argocd.svc.cluster.local"
	r := GetArgocdServerAddr()
	assert.Equal(t, expected, r)
	if expected != r {
		t.Errorf("getArgocdServerAddrExistingUrl with default url failed, expected %s, got %s", expected, r)
	}
}
