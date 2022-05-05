package aws

const (
	AuthModeProfile = "profile"
	AuthModeStatic  = "static"
	AuthModeDefault = "default"
)

type Credentials struct {
	AccessKeyID     string `json:"accessKeyId"`
	SecretAccessKey string `json:"secretAccessKey"`
	SessionToken    string `json:"sessionToken"`
}

type Auth struct {
	Mode            string `json:"mode"`
	Profile         string `json:"profile"`
	Region          string `json:"region"`
	AccessKeyID     string `json:"accessKeyID"`
	SecretAccessKey string `json:"secretAccessKey"`
}
