package terraform

import (
	stablev1 "github.com/compuzest/zlifecycle-il-operator/api/v1"
	"github.com/compuzest/zlifecycle-il-operator/controller/codegen/secret"
)

func createSecretsConfig(arr []*stablev1.Secret, identifier secret.Identifier) (*TerraformSecretsConfig, error) {
	scopedSecrets := make([]*Secret, 0, len(arr))
	for _, s := range arr {
		scope := s.Scope
		if scope == "" {
			scope = "component"
		}
		key, err := identifier.GenerateKey(s.Key, scope)
		if err != nil {
			return nil, err
		}
		scopedSecrets = append(scopedSecrets, &Secret{Key: key, Name: s.Name})
	}

	return &TerraformSecretsConfig{Secrets: scopedSecrets}, nil
}
