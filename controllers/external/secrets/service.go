package secrets

//go:generate mockgen --build_flags=--mod=mod -destination=./mock_secrets_api.go -package=secrets "github.com/compuzest/zlifecycle-il-operator/controllers/external/secrets" API
type API interface {
	GetSecret(key string) (*Secret, error)
	GetSecrets(keys ...string) ([]*Secret, error)
}
