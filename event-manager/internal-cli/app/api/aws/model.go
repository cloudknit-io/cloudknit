package aws

const (
	AuthModeProfile = "profile"
	AuthModeDefault = "default"
)

type Credentials struct {
	AccessKeyID     string `json:"accessKeyId"`
	SecretAccessKey string `json:"secretAccessKey"`
	SessionToken    string `json:"sessionToken"`
}

type Auth struct {
	Mode    string `json:"mode"`
	Profile string `json:"profile"`
	Region  string `json:"region"`
}
