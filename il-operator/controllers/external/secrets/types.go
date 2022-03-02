package secrets

type AWSCreds struct {
	AccessKeyID     string `json:"accessKeyID"`
	SecretAccessKey string `json:"secretAccessKey"`
	SessionToken    string `json:"sessionToken"`
}

type Secret struct {
	Exists bool    `json:"exists"`
	Key    string  `json:"key"`
	Value  *string `json:"value"`
}
