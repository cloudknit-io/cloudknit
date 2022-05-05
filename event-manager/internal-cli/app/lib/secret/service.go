package secret

type API interface {
	GetSecret(key string) (*Secret, error)
	GetSecrets(keys ...string) ([]*Secret, error)
}

type Secret struct {
	Value  *string
	Key    string
	Exists bool
}
