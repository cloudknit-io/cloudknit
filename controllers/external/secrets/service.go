package secrets

import (
	"context"
)

//go:generate mockgen --build_flags=--mod=mod -destination=./mock_secrets_api.go -package=secrets "github.com/compuzest/zlifecycle-il-operator/controllers/external/secrets" API
type API interface {
	GetSecret(ctx context.Context, key string) (*Secret, error)
	GetSecrets(ctx context.Context, keys ...string) ([]*Secret, error)
}
