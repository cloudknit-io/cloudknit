package secrets

type TerraformStateConfig struct {
	Bucket    string `json:"bucket"`
	LockTable string `json:"lockTable"`
}

type Secret struct {
	Exists bool    `json:"exists"`
	Key    string  `json:"key"`
	Value  *string `json:"value"`
}

type AWSCredentials struct {
	AccessKeyID     string `json:"accessKeyId"`
	SecretAccessKey string `json:"secretAccessKey"`
	SessionToken    string `json:"sessionToken"`
}
