package secret

import (
	"context"
)

//go:generate mockgen --build_flags=--mod=mod -destination=./mock_secret_api.go -package=secret "github.com/compuzest/zlifecycle-il-operator/controller/common/secret" API
type API interface {
	GetSecret(ctx context.Context, key string) (*Secret, error)
	GetSecrets(ctx context.Context, keys ...string) ([]*Secret, error)
}
