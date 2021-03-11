package github

import "os"

func GetWebhookSecret() string {
	secret, exists := os.LookupEnv("GITHUB_WEBHOOK_SECRET")
	if exists {
		return secret
	} else {
		return ""
	}
}
