package github

import "os"

func GetWebhookSecret() string {
	secret, exists := os.LookupEnv("github_webhook_secret")
	if exists {
		return secret
	} else {
		return "C0mpuZ3s7"
	}
}
