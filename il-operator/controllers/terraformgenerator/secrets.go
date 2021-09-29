package terraformgenerator

import (
	stablev1 "github.com/compuzest/zlifecycle-il-operator/api/v1"
	"github.com/compuzest/zlifecycle-il-operator/controllers/secrets"
)

func createSecretsConfig(secretArray []*stablev1.Secret, meta secrets.SecretMeta) (*TerraformSecretsConfig, error) {
	scopedSecrets := make([]Secret, 0, len(secretArray))
	for _, s := range secretArray {
		scope := s.Scope
		if scope == "" {
			scope = "component"
		}
		key, err := secrets.CreateKey(s.Key, scope, meta)
		if err != nil {
			return nil, err
		}
		scopedSecrets = append(scopedSecrets, Secret{Key: key, Name: s.Name})
	}
	conf := TerraformSecretsConfig{Secrets: scopedSecrets}

	return &conf, nil
}
